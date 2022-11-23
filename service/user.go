package service

import (
	"encoding/json"
	"fmt"
	"goWeb/middleware"
	"goWeb/model"
	"gorm.io/gorm"
	"log"
	"time"
)

type IUserService interface {
	CreateUser(user *model.User)
	DeleteUser(id string)
	GetUser(user *model.User)
	UpdateUser(user *model.User)
	BatchUpdateUser(users []*model.User)
	BatchCreateUser(user []*model.User)
	PageUser(users *[]*model.User)
}

var UserServiceImpl IUserService

type UserService struct {
}

func (u UserService) PageUser(users *[]*model.User) {
	middleware.DbClient.Find(users)
	for i, val := range *users {
		log.Println("user: ", i, val.ID, val.Name)

	}
}

func (u UserService) CreateUser(user *model.User) {
	db := middleware.DbClient
	db.Create(user)
	marshal, err := json.Marshal(user)
	if err != nil {
		log.Println("error cache", user.ID)
		return
	}
	middleware.RedisClient.Set(fmt.Sprintf("user::%s", user.ID), string(marshal), time.Hour)
}

func (u UserService) DeleteUser(id string) {
	middleware.DbClient.Delete(&model.User{}, id)
	middleware.RedisClient.Del(fmt.Sprintf("user::%s", id))
}

func (u UserService) GetUser(user *model.User) {
	userRedis := middleware.RedisClient.Get(fmt.Sprintf("user::%s", user.ID))
	bytes, err := userRedis.Bytes()
	if len(bytes) > 0 {

		err = json.Unmarshal(bytes, user)
		if err != nil {
			log.Printf("err %s \n", err)
		}
		log.Println("from redis")

	} else {

		db := middleware.DbClient
		db.First(user, user.ID)
		marshal, err := json.Marshal(user)
		if err != nil {
			log.Println("error cache", user.ID)
		}
		middleware.RedisClient.Set(fmt.Sprintf("user::%s", user.ID), string(marshal), time.Hour)

		log.Println("from db")
	}
}

func (u UserService) UpdateUser(user *model.User) {
	middleware.DbClient.Save(&user)
	middleware.RedisClient.Del(fmt.Sprintf("user::%s", user.ID))
}

func (u UserService) BatchUpdateUser(users []*model.User) {
	middleware.DbClient.Transaction(func(tx *gorm.DB) error {
		tx.Transaction(func(tx *gorm.DB) error {
			return nil
		})

		rows, err := tx.Model(&model.User{}).Find(&model.User{}).Rows()
		if err != nil {
			log.Println("error get allUser")
			return err
		}
		var users []*model.User
		rows.Scan(users)
		for rows.Next() {
			row := &model.User{}
			err := rows.Scan(row)
			log.Println("user result:", row)
			if err != nil {
				return err
			}
		}
		return nil

	})
	return
}

func (u UserService) BatchCreateUser(user []*model.User) {
	//TODO implement me
	panic("implement me")
}

func init() {
	UserServiceImpl = &UserService{}
}
