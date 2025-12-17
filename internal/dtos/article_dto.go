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
