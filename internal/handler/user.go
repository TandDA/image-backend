package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/TandDA/image-beckend/internal/model"
	jwt "github.com/TandDA/image-beckend/internal/util"
	"github.com/labstack/echo/v4"
)

func (h *Handler) createUser(c echo.Context) error {
	query := `
	INSERT INTO users(
		username, email, password_hash)
		VALUES ($1, $2, $3);
	`
	var dto model.UserDTO
	if err := c.Bind(&dto); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON provided: " + err.Error()})
	}
	_, err := h.db.Exec(query, dto.Username, dto.Email, dto.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusOK)
}

type userAuthDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) AuthUser(c echo.Context) error {
	requestUser := userAuthDTO{}
	if err := c.Bind(&requestUser); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON: " + err.Error()})
	}
	fmt.Println(requestUser)
	query := `
	SELECT id, password_hash FROM users WHERE username = $1;
	`
	row := h.db.QueryRow(query, requestUser.Username)
	var user model.User
	if err := row.Scan(&user.Id, &user.PasswordHash); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get user: " + err.Error()})
	}
	if user.PasswordHash != requestUser.Password {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}
	jwtStr, err := jwt.GenerateJWT(user.Id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "InternalServerError: " + err.Error()})
	}
	return c.JSON(200, map[string]string{"jwt": jwtStr})
}
