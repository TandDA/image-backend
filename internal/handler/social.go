package handler

import (
	"net/http"

	"github.com/TandDA/image-beckend/internal/model"
	jwt "github.com/TandDA/image-beckend/internal/util"
	"github.com/labstack/echo/v4"
)

func (h *Handler) likePhoto(c echo.Context) error {
	photoId := c.QueryParam("photoId")
	if photoId == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing path param photoId"})
	}
	userId := jwt.GetUserIdFromContext(c)

	query := `
		INSERT INTO likes(photo_id, user_id) VALUES($1,$2);
	`

	_, err := h.db.Exec(query, photoId, userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Faileto to insert like: " + err.Error()})
	}
	return c.NoContent(http.StatusOK)
}

func (h *Handler) postComment(c echo.Context) error {
	var comm model.Comment
	if err := c.Bind(&comm); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Wrong JSON struct:" + err.Error()})
	}
	userId := jwt.GetUserIdFromContext(c)

	query := `
		INSERT INTO comments(photo_id,user_id,text) VALUES($1,$2,$3)
	`
	_, err := h.db.Exec(query, comm.PhotoId, userId, comm.Text)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Faileto to insert comment: " + err.Error()})
	}
	return c.NoContent(http.StatusOK)
}

type PlainComment struct {
	Username string `json:"username"`
	Text string `json:"text"`
}

func (h *Handler) getCommentsByPhotoId(c echo.Context) error {
	photoId := c.Param("id")
	query := `
	SELECT u.username, c.text FROM comments AS c
	JOIN users AS u ON u.id = c.user_id 
	WHERE photo_id = $1
	`

	rows, err := h.db.Query(query, photoId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Faileto to get all comments: " + err.Error()})
	}
	defer rows.Close()

	comments := make([]PlainComment, 0)
	for rows.Next() {
		var comment PlainComment
		err = rows.Scan(&comment.Username, &comment.Text)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Faileto read comment: " + err.Error()})
		}
		comments = append(comments, comment)
	}
	return c.JSON(http.StatusOK, comments)
}