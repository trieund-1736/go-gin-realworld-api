package dtos

type ListArticlesQuery struct {
	Tag       string `form:"tag"`
	Author    string `form:"author"`
	Favorited string `form:"favorited"`
	Limit     int    `form:"limit,default=20"`
	Offset    int    `form:"offset,default=0"`
}

type ArticleAuthorResponse struct {
	Username string `json:"username"`
}

type ArticleResponse struct {
	Slug           string                `json:"slug"`
	Title          string                `json:"title"`
	Description    string                `json:"description"`
	Body           string                `json:"body"`
	TagList        []string              `json:"tagList"`
	CreatedAt      string                `json:"createdAt"`
	UpdatedAt      string                `json:"updatedAt"`
	Favorited      bool                  `json:"favorited"`
	FavoritesCount int                   `json:"favoritesCount"`
	Author         ArticleAuthorResponse `json:"author"`
}

type ArticlesListResponse struct {
	Articles      []ArticleResponse `json:"articles"`
	ArticlesCount int               `json:"articlesCount"`
}

type CreateArticleRequest struct {
	Article struct {
		Title       string   `json:"title" binding:"required"`
		Description string   `json:"description" binding:"required"`
		Body        string   `json:"body" binding:"required"`
		TagList     []string `json:"tagList"`
	} `json:"article" binding:"required"`
}

type UpdateArticleRequest struct {
	Article struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Body        string   `json:"body"`
		TagList     []string `json:"tagList"`
	} `json:"article" binding:"required"`
}

type ArticleDetailResponse struct {
	Article ArticleResponse `json:"article"`
}
