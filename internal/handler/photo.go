package handler

import (
	"io"
	"net/http"
	"os"

	"github.com/TandDA/image-beckend/internal/model"
	jwt "github.com/TandDA/image-beckend/internal/util"
	"github.com/labstack/echo/v4"
)

func (h *Handler) getAllPhotos(c echo.Context) error {
	query := `
	SELECT p.id, p.user_id, p.image_url, COUNT(photo_id) AS likes FROM photos AS p
	LEFT JOIN likes AS l ON l.photo_id = p.id
	GROUP BY p.id;
	`
	rows, err := h.db.Query(query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	var photos []model.Photo
	for rows.Next() {
		var photo model.Photo
		err = rows.Scan(&photo.ID, &photo.UserID, &photo.ImageURL, &photo.Likes)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		photos = append(photos, photo)
	}
	return c.JSON(http.StatusOK, photos)
}

func (h *Handler) getPhotoByPath(c echo.Context) error {
	filePath := c.QueryParam("path")
	if filePath == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Не указан путь к файлу")
	}

	// Проверка существования файла
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return echo.NewHTTPError(http.StatusNotFound, "Файл не найден")
	}

	// Чтение файла
	file, err := os.ReadFile(filePath)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Ошибка при чтении файла")
	}

	// Возврат файла
	return c.Blob(http.StatusOK, "image/jpeg", file)
}

func (h *Handler) postPhoto(c echo.Context) error {
	// Парсинг формы с файлом
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer src.Close()

	// Определение пути для сохранения файла
	dst, err := os.Create("img/" + file.Filename)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "File create error: " + err.Error()})
	}
	defer dst.Close()

	// Копирование файла в папку проекта
	if _, err = io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	userId := jwt.GetUserIdFromContext(c)
	query := `
	INSERT INTO photos(
		user_id, image_url)
		VALUES ($1, $2) RETURNING id;
	`

	row := h.db.QueryRow(query, userId, "img/"+file.Filename)
	var photoId int
	if err = row.Scan(&photoId); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]int{
		"id": photoId,
	})
}

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
