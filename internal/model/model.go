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
		ID          int
		UserID      int
		ImageURL    string
	}
)

/*

   id SERIAL PRIMARY KEY,
   username VARCHAR(255) UNIQUE NOT NULL,
   email VARCHAR(255) UNIQUE NOT NULL,
   password_hash VARCHAR(255) NOT NULL

*/
