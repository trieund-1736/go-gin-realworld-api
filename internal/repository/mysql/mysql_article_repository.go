package mysql

import (
	"go-gin-realworld-api/internal/models"

	"gorm.io/gorm"
)

type MySqlArticleRepository struct {
}

func NewMySqlArticleRepository() *MySqlArticleRepository {
	return &MySqlArticleRepository{}
}

// ListArticles lists articles with optional filtering and pagination
func (r *MySqlArticleRepository) ListArticles(db *gorm.DB, tag, author string, favorited *bool, currentUserID *int64, limit, offset int) ([]*models.Article, int64, error) {
	var articles []*models.Article
	var total int64

	query := db

	// Filter by tag
	if tag != "" {
		query = query.
			Joins("JOIN article_tags ON article_tags.article_id = articles.id").
			Joins("JOIN tags ON tags.id = article_tags.tag_id").
			Where("tags.name = ?", tag)
	}

	// Filter by author
	if author != "" {
		query = query.
			Joins("JOIN users AS author_user ON author_user.id = articles.author_id").
			Where("author_user.username = ?", author)
	}

	// Filter by favorited (whether current user has favorited the article)
	if favorited != nil && currentUserID != nil {
		if *favorited {
			// Get articles favorited by current user
			query = query.
				Joins("JOIN favorites ON favorites.article_id = articles.id").
				Where("favorites.user_id = ?", *currentUserID)
		} else {
			// Get articles NOT favorited by current user
			query = query.
				Where("articles.id NOT IN (SELECT article_id FROM favorites WHERE user_id = ?)", *currentUserID)
		}
	}

	// Get total count before pagination
	if err := query.Model(&models.Article{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting and pagination
	if err := query.
		Preload("Author").
		Preload("ArticleTags.Tag").
		Preload("Favorites").
		Order("articles.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// FeedArticles gets articles from followed users
func (r *MySqlArticleRepository) FeedArticles(db *gorm.DB, userID int64, limit, offset int) ([]*models.Article, int64, error) {
	var articles []*models.Article
	var total int64

	query := db.
		Joins("JOIN follows ON follows.followee_id = articles.author_id").
		Where("follows.follower_id = ?", userID)

	// Get total count
	if err := query.Model(&models.Article{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting and pagination
	if err := query.
		Preload("Author").
		Preload("ArticleTags.Tag").
		Preload("Favorites").
		Order("articles.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// FindArticleBySlug finds an article by slug
func (r *MySqlArticleRepository) FindArticleBySlug(db *gorm.DB, slug string) (*models.Article, error) {
	var article *models.Article
	if err := db.
		Preload("Author").
		Preload("ArticleTags.Tag").
		Preload("Favorites").
		Where("slug = ?", slug).
		First(&article).Error; err != nil {
		return nil, err
	}
	return article, nil
}

// CreateArticle creates a new article
func (r *MySqlArticleRepository) CreateArticle(db *gorm.DB, article *models.Article) error {
	if err := db.Create(article).Error; err != nil {
		return err
	}
	return nil
}

// UpdateArticle updates an article
func (r *MySqlArticleRepository) UpdateArticle(db *gorm.DB, article *models.Article) error {
	if err := db.Model(article).Omit("Favorites").Save(article).Error; err != nil {
		return err
	}
	return nil
}

// DeleteArticleBySlug deletes an article by slug
func (r *MySqlArticleRepository) DeleteArticleBySlug(db *gorm.DB, slug string) error {
	if err := db.Where("slug = ?", slug).Delete(&models.Article{}).Error; err != nil {
		return err
	}
	return nil
}

// AssignTagsToArticle associates tags with an article
func (r *MySqlArticleRepository) AssignTagsToArticle(db *gorm.DB, articleID int64, tagNames []string) error {
	if len(tagNames) == 0 {
		return nil
	}

	// Delete existing article tags
	if err := db.Where("article_id = ?", articleID).Delete(&models.ArticleTag{}).Error; err != nil {
		return err
	}

	// Find existing tags by name
	var existingTags []*models.Tag
	if err := db.Where("name IN ?", tagNames).Find(&existingTags).Error; err != nil {
		return err
	}

	// Create a map of existing tag names for quick lookup
	existingTagMap := make(map[string]*models.Tag)
	for _, tag := range existingTags {
		existingTagMap[tag.Name] = tag
	}

	// Identify tags that need to be created
	var tagsToCreate []*models.Tag
	for _, tagName := range tagNames {
		if _, exists := existingTagMap[tagName]; !exists {
			tagsToCreate = append(tagsToCreate, &models.Tag{Name: tagName})
		}
	}

	// Create new tags in bulk
	if len(tagsToCreate) > 0 {
		if err := db.CreateInBatches(tagsToCreate, 100).Error; err != nil {
			return err
		}
		// Add newly created tags to the map
		for _, tag := range tagsToCreate {
			existingTagMap[tag.Name] = tag
		}
	}

	// Create article tags in bulk
	var articleTags []*models.ArticleTag
	for _, tagName := range tagNames {
		tag := existingTagMap[tagName]
		articleTags = append(articleTags, &models.ArticleTag{
			ArticleID: articleID,
			TagID:     tag.ID,
		})
	}

	if err := db.CreateInBatches(articleTags, 100).Error; err != nil {
		return err
	}

	return nil
}
