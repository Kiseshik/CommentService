package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Kiseshik/CommentService.git/internal/controllers/api/dto"
	"github.com/Kiseshik/CommentService.git/internal/core/port"
)

// проверка на этапе компиляции что структура api реализует интерфейс port.Handler
// var _ port.Handler = (*api)(nil)

type api struct {
	postService    port.PostService
	commentService port.CommentService
}

func NewApiImplementation(
	postService port.PostService,
	commentService port.CommentService,
) port.Handler {
	return &api{
		postService:    postService,
		commentService: commentService,
	}
}

// RegisterPublicHandlers регистрирует публичные эндпоинты (без авторизации)
// RegisterPrivateHandlers регистрирует приватные эндпоинты (с авторизацией)
// RegisterInternalHandlers регистрирует внутренние эндпоинты (для других сервисов)

func (api *api) RegisterPublicHandlers(group *gin.RouterGroup) {
	group.POST("/health", api.Health)

	group.POST("/posts/list", api.ListPosts)
	group.POST("/posts/create", api.CreatePost)
	group.POST("/posts/get", api.GetPostByID)
	group.POST("/posts/update", api.UpdatePost)
	group.POST("/posts/toggle-comments", api.ToggleComments)

	group.POST("/comments/list", api.ListComments)
	group.POST("/comments/create", api.CreateComment)
}

func (api *api) RegisterPrivateHandlers(group *gin.RouterGroup) {
	// TODO: добавить JWT middleware и перенести сюда приватные эндпоинты
}

func (api *api) RegisterInternalHandlers(group *gin.RouterGroup) {
}

func (api *api) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (api *api) ListPosts(c *gin.Context) {
	posts, err := api.postService.ListPosts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.PostResponse, len(posts))
	for i, p := range posts {
		response[i] = dto.PostResponse{
			ID:              p.ID,
			Title:           p.Title,
			Content:         p.Content,
			AuthorID:        p.AuthorID,
			CommentsEnabled: p.CommentsEnabled,
			CreatedAt:       p.CreatedAt,
			UpdatedAt:       p.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

func (api *api) CreatePost(c *gin.Context) {
	var req dto.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := &port.CreatePostParams{
		Title:           req.Title,
		Content:         req.Content,
		AuthorID:        req.AuthorID,
		CommentsEnabled: req.CommentsEnabled,
	}

	post, err := api.postService.CreatePost(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.PostResponse{
		ID:              post.ID,
		Title:           post.Title,
		Content:         post.Content,
		AuthorID:        post.AuthorID,
		CommentsEnabled: post.CommentsEnabled,
		CreatedAt:       post.CreatedAt,
		UpdatedAt:       post.UpdatedAt,
	})
}

func (api *api) GetPostByID(c *gin.Context) {
	var req dto.GetPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post, err := api.postService.GetPostByID(c.Request.Context(), req.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.PostResponse{
		ID:              post.ID,
		Title:           post.Title,
		Content:         post.Content,
		AuthorID:        post.AuthorID,
		CommentsEnabled: post.CommentsEnabled,
		CreatedAt:       post.CreatedAt,
		UpdatedAt:       post.UpdatedAt,
	})
}

func (api *api) UpdatePost(c *gin.Context) {
	var req dto.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	params := &port.UpdatePostParams{
		ID:              req.ID,
		Title:           req.Title,
		Content:         req.Content,
		CommentsEnabled: req.CommentsEnabled,
	}
	post, err := api.postService.UpdatePost(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.PostResponse{
		ID:              post.ID,
		Title:           post.Title,
		Content:         post.Content,
		AuthorID:        post.AuthorID,
		CommentsEnabled: post.CommentsEnabled,
		CreatedAt:       post.CreatedAt,
		UpdatedAt:       post.UpdatedAt,
	})
}

func (api *api) ToggleComments(c *gin.Context) {
	var req dto.ToggleCommentsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := api.postService.ToggleComments(c.Request.Context(), req.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	post, err := api.postService.GetPostByID(c.Request.Context(), req.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.PostResponse{
		ID:              post.ID,
		Title:           post.Title,
		Content:         post.Content,
		AuthorID:        post.AuthorID,
		CommentsEnabled: post.CommentsEnabled,
		CreatedAt:       post.CreatedAt,
		UpdatedAt:       post.UpdatedAt,
	})
}

func (api *api) ListComments(c *gin.Context) {
	var req dto.ListCommentsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := &port.ListCommentParams{
		PostID:   req.PostID,
		ParentID: req.ParentID,
		Cursor:   req.Cursor,
		Limit:    req.Limit,
	}

	result, err := api.commentService.ListComments(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	comments := make([]dto.CommentResponse, len(result.Comments))
	for i, c := range result.Comments {
		comments[i] = dto.CommentResponse{
			ID:        c.ID,
			PostID:    c.PostID,
			ParentID:  c.ParentID,
			Content:   c.Content,
			AuthorID:  c.AuthorID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, dto.CommentListResponse{
		Comments:    comments,
		HasNextPage: result.HasNextPage,
		NextCursor:  result.NextCursor,
	})
}

func (api *api) CreateComment(c *gin.Context) {
	var req dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := &port.CreateCommentParams{
		PostID:   req.PostID,
		ParentID: req.ParentID,
		AuthorID: req.AuthorID,
		Content:  req.Content,
	}

	comment, err := api.commentService.CreateComment(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.CommentResponse{
		ID:        comment.ID,
		PostID:    comment.PostID,
		ParentID:  comment.ParentID,
		Content:   comment.Content,
		AuthorID:  comment.AuthorID,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	})
}
