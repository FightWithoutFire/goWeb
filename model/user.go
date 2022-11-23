package model

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"log"
)

type User struct {
	gorm.Model
	ID   string `gorm:"id" uri:"id" json:"id"`
	Name string `gorm:"name" json:"name"`
	//CheckIn  time.Time `gorm:"check_in" json:"checkIn" binding:"required,bookabledate" time_format:"2006-01-02"`
	//CheckOut time.Time `gorm:"check_out" json:"checkOut" binding:"required,gtfield=CheckIn,bookabledate" time_format:"2006-01-02"`
}

func (u *User) BeforeSave(db *gorm.DB) error {

	return nil
}

//var bookableDate validator.Func = func(fl validator.FieldLevel) bool {
//	date, ok := fl.Field().Interface().(time.Time)
//	if ok {
//		today := time.Now()
//		if today.After(date) {
//			return false
//		}
//	}
//	return true
//}

//func init() {
//	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
//		v.RegisterValidation("bookabledate", bookableDate)
//	}
//}

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
