package model

import (
	"github.com/go-playground/validator/v10"
	_ "github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"log"
)

type User struct {
	gorm.Model
	ID   string `gorm:"id" uri:"id" json:"id"`
	Name string `gorm:"name" json:"name"`
}

type Attach struct {
}

type AddUserParam struct {
	Name string
}

type RequestParam struct {
	Username string `json:"name" validate:"gt=0"`
	Password string `json:"password" validate:"gt=0"`
	Age      int    `json:"age" validate:"gt=0"`
	Email    string `json:"email" validate:"email"`
}

var validate = validator.New()

func ValidateFuc(req RequestParam) error {

	err := validate.Struct(req)
	if err != nil {
		log.Println("error:", err)
	}

	return err
}
