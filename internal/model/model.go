package model

type (
	User struct {
		Id           int
		Username     string
		Email        string
		PasswordHash string
	}
	UserDTO struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	Photo struct {
		ID       int
		UserID   int    `json:"user_id"`
		ImageURL string `json:"image_url"`
		Likes    int    `json:"likes"`
	}
	Like struct {
		Id      int `json:"-"`
		PhotoId int `json:"photo_id"`
		UserId  int `json:"user_id"`
	}
)

/*

   id SERIAL PRIMARY KEY,
   username VARCHAR(255) UNIQUE NOT NULL,
   email VARCHAR(255) UNIQUE NOT NULL,
   password_hash VARCHAR(255) NOT NULL

*/
