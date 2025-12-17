package dtos

type CommentAuthorResponse struct {
	Username string `json:"username"`
}

type CommentResponse struct {
	ID        int64                 `json:"id"`
	CreatedAt string                `json:"createdAt"`
	UpdatedAt string                `json:"updatedAt"`
	Body      string                `json:"body"`
	Author    CommentAuthorResponse `json:"author"`
}

type CommentsListResponse struct {
	Comments []CommentResponse `json:"comments"`
}

type CreateCommentRequest struct {
	Comment struct {
		Body string `json:"body" binding:"required"`
	} `json:"comment" binding:"required"`
}

type CommentDetailResponse struct {
	Comment CommentResponse `json:"comment"`
}
