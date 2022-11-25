package model

import (
	"crypto/rand"
	"crypto/rsa"
	"database/sql/driver"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"goWeb/middleware"
	"gorm.io/gorm"
	"log"
	"strings"
)

func init() {
	err := middleware.DbClient.Migrator().AutoMigrate(&User{})
	if err != nil {
		log.Fatalln("AutoMigrate error", err)
	}
}

type PageObj struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

func Paginate(context *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		log.Println(context.Request.RequestURI)
		pageObj := &PageObj{}
		err := context.ShouldBindQuery(pageObj)
		if err != nil {
			log.Println("error binding:", err)
		}
		log.Println("page", pageObj.Page)
		if pageObj.Page == 0 {
			pageObj.Page = 1
		}

		log.Println("page_size", pageObj.PageSize)

		switch {
		case pageObj.PageSize > 100:
			pageObj.PageSize = 100
		case pageObj.PageSize <= 0:
			pageObj.PageSize = 10
		}

		offset := (pageObj.Page - 1) * pageObj.PageSize
		return db.Offset(offset).Limit(pageObj.PageSize)
	}
}

type EncryptString string

func (e *EncryptString) Value() (driver.Value, error) {
	password := fmt.Sprint(*e)
	if strings.HasPrefix("!E!", password) {
		return password, nil
	}

	encryptedBytes, err := rsa.EncryptPKCS1v15(
		rand.Reader,
		middleware.PublicKey,
		[]byte(password))
	if err != nil {
		log.Println("value err", err)
		return nil, err
	}
	value := "!E!" + base64.StdEncoding.EncodeToString(encryptedBytes)

	return value, nil
}

func (e *EncryptString) Scan(value interface{}) error {
	bytes, ok := value.(string)
	if !ok {
		return errors.New(fmt.Sprint("Failed to decrypt value:", value))
	}

	stemp := bytes
	log.Println(stemp)
	stemp = stemp[len("!E!"):]
	log.Println(stemp)
	decodeString, err := base64.StdEncoding.DecodeString(stemp)
	if err != nil {
		return errors.New(fmt.Sprint("Failed to Decode from base64 format:", stemp))
	}
	values, err := rsa.DecryptPKCS1v15(rand.Reader,
		middleware.PrivateKey,
		decodeString)
	if err != nil {
		return errors.New(fmt.Sprint("!!Failed to decrypt value1:", err))
	}
	result := string(values)
	*e = EncryptString(result)
	return err
}

type User struct {
	gorm.Model
	ID       string `gorm:"id" uri:"id" json:"id"`
	Name     string `gorm:"name" json:"name"`
	Password *EncryptString
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
