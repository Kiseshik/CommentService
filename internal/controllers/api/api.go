package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Kiseshik/CommentService.git/internal/controllers/api/dto"
	"github.com/Kiseshik/CommentService.git/internal/core/port"
	"github.com/Kiseshik/CommentService.git/internal/core/service"
)

// проверка на этапе компиляции что структура api реализует интерфейс port.Handler
// var _ port.Handler = (*api)(nil)

type api struct {
	postService    *service.PostService
	commentService *service.CommentService
}

func NewApiImplementation(
	postService *service.PostService,
	commentService *service.CommentService,
) port.Handler {
	return &api{
		postService:    postService,
		commentService: commentService,
	}
}

// RegisterPublicHandlers регистрирует публичные эндпоинты (без авторизации)
// RegisterPrivateHandlers регистрирует приватные эндпоинты (с авторизацией)
// RegisterInternalHandlers регистрирует внутренние эндпоинты (для других сервисов)

func (a *api) RegisterPublicHandlers(group *gin.RouterGroup) {
	group.POST("/health", a.Health)

	group.POST("/posts/list", a.ListPosts)
	group.POST("/posts/create", a.CreatePost)
	group.POST("/posts/get", a.GetPostByID)
	group.POST("/posts/update", a.UpdatePost)
	group.POST("/posts/toggle-comments", a.ToggleComments)

	group.POST("/comments/list", a.ListComments)
	group.POST("/comments/create", a.CreateComment)
}

func (a *api) RegisterPrivateHandlers(group *gin.RouterGroup) {
	// TODO: добавить JWT middleware и перенести сюда приватные эндпоинты
}

func (a *api) RegisterInternalHandlers(group *gin.RouterGroup) {
}

func (a *api) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (a *api) ListPosts(c *gin.Context) {
	posts, err := a.postService.ListPosts(c.Request.Context())
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

func (a *api) CreatePost(c *gin.Context) {
	var req dto.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post, err := a.postService.CreatePost(
		c.Request.Context(),
		req.Title,
		req.Content,
		req.AuthorID,
		req.CommentsEnabled,
	)
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

func (a *api) GetPostByID(c *gin.Context) {
	var req dto.GetPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post, err := a.postService.GetPostByID(c.Request.Context(), req.ID)
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

func (a *api) UpdatePost(c *gin.Context) {
	var req dto.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post, err := a.postService.GetPostByID(c.Request.Context(), req.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if req.Title != nil {
		post.Title = *req.Title
	}
	if req.Content != nil {
		post.Content = *req.Content
	}
	if req.CommentsEnabled != nil {
		post.CommentsEnabled = *req.CommentsEnabled
	}

	if err := a.postService.UpdatePost(c.Request.Context(), post); err != nil {
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

func (a *api) ToggleComments(c *gin.Context) {
	var req dto.ToggleCommentsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := a.postService.ToggleComments(c.Request.Context(), req.ID, req.Enabled); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	post, err := a.postService.GetPostByID(c.Request.Context(), req.ID)
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

func (a *api) ListComments(c *gin.Context) {
	var req dto.ListCommentsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := a.commentService.ListComments(
		c.Request.Context(),
		req.PostID,
		req.ParentID,
		req.Cursor,
		req.Limit,
	)
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

func (a *api) CreateComment(c *gin.Context) {
	var req dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment, err := a.commentService.CreateComment(
		c.Request.Context(),
		req.PostID,
		req.ParentID,
		req.AuthorID,
		req.Content,
	)
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
