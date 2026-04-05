package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/Kiseshik/CommentService.git/internal/adapters/repository/memory"
	"github.com/Kiseshik/CommentService.git/internal/adapters/repository/postgres"
	"github.com/Kiseshik/CommentService.git/internal/config"
	"github.com/Kiseshik/CommentService.git/internal/controllers/api"
	"github.com/Kiseshik/CommentService.git/internal/core/port"
	"github.com/Kiseshik/CommentService.git/internal/core/service"
)

type App struct {
	config       *config.Config
	ctx          context.Context
	server       *http.Server
	database     *sqlx.DB
	repositories struct {
		post    port.PostRepository
		comment port.CommentRepository
	}
	services struct {
		post    *service.PostService
		comment *service.CommentService
	}
	stopHandlers []func()
}

func New(cfg *config.Config) (*App, error) {
	app := &App{
		config: cfg,
		ctx:    context.Background(),
	}

	if err := app.initRepositories(); err != nil {
		return nil, fmt.Errorf("initRepositories: %w", err)
	}

	app.initServices()

	if err := app.initHTTPServer(); err != nil {
		return nil, fmt.Errorf("initHTTPServer: %w", err)
	}

	return app, nil
}

func (app *App) initRepositories() error {
	if app.config.IsPostgresStorage() {
		db, err := sqlx.Connect("postgres", app.config.PostgresDSN)
		if err != nil {
			return fmt.Errorf("failed to connect to postgres: %w", err)
		}
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(5 * time.Minute)
		app.database = db
		app.repositories.post = postgres.NewPostRepository(db)
		app.repositories.comment = postgres.NewCommentRepository(db)
		app.RegisterStopHandler(func() {
			_ = app.database.Close()
		})
	} else {
		app.repositories.post = memory.NewPostRepository()
		app.repositories.comment = memory.NewCommentRepository()
	}
	return nil
}

func (app *App) initServices() {
	app.services.post = service.NewPostService(app.repositories.post)
	app.services.comment = service.NewCommentService(app.repositories.comment, app.repositories.post)
}

func (app *App) initHTTPServer() error {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())

	apiHandler := api.NewApiImplementation(app.services.post, app.services.comment)

	internalGroup := engine.Group("/internal/api/v1")
	apiHandler.RegisterInternalHandlers(internalGroup)

	privateGroup := engine.Group("/api/v1")
	//todo jwt auth
	//privateGroup.Use(middleware.AuthorizeJWT())
	apiHandler.RegisterPrivateHandlers(privateGroup)

	publicGroup := engine.Group("/public/api/v1")
	apiHandler.RegisterPublicHandlers(publicGroup)

	app.server = &http.Server{
		Addr:    app.config.ListenAddr,
		Handler: engine,
	}

	app.RegisterStopHandler(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = app.server.Shutdown(ctx)
	})

	return nil
}

func (app *App) Run() error {
	go func() {
		log.Printf("server starting on %s", app.config.ListenAddr)
		if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("server error: %v", err)
		}
	}()
	<-app.ctx.Done()
	return nil
}

func (app *App) Shutdown() error {
	for _, handler := range app.stopHandlers {
		handler()
	}
	return nil
}

func (app *App) RegisterStopHandler(handler func()) {
	app.stopHandlers = append(app.stopHandlers, handler)
}
