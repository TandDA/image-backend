package handler

import (
	"database/sql"

	jwt "github.com/TandDA/image-beckend/internal/util"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) Start() {
	e := echo.New()

	userGroup := e.Group("/user")
	userGroup.POST("/create", h.createUser)
	userGroup.POST("/auth", h.AuthUser)

	photoGroup := e.Group("/photo")
	photoGroup.Use(userAuthMiddleware)
	photoGroup.POST("/post", h.postPhoto)
	photoGroup.POST("/post/comment", h.postComment)
	photoGroup.GET("/comment/:id", h.getCommentsByPhotoId)
	photoGroup.POST("/like", h.likePhoto)
	photoGroup.GET("/path", h.getPhotoByPath)
	photoGroup.GET("", h.getAllPhotos)
	e.Logger.Fatal(e.Start(":8080"))
}

func userAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := jwt.ValidateJWT(c)
		if err != nil {
			return c.JSON(401, map[string]string{"error": err.Error()})
		}
		return next(c)
	}
}
