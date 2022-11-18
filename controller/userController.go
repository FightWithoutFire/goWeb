package controller

import (
	"github.com/gin-gonic/gin"
	"goWeb/middleware"
	"goWeb/model"
	"goWeb/service"
	"log"
	"net/http"
)

const (
	AddUserUri  = "/users"
	RestUserUri = "/users/:id"
)

func init() {
	middleware.GinEngine.POST(AddUserUri, CreateUserFun())
	middleware.GinEngine.PUT(RestUserUri, UpdateUserFun())
	middleware.GinEngine.GET(RestUserUri, GetUserFun())
	middleware.GinEngine.DELETE(RestUserUri, DeleteUserFun())
}

func CreateUserFun() func(context *gin.Context) {
	return func(context *gin.Context) {

		user := &model.User{}
		err := context.Bind(user)
		if err != nil {
			log.Println("error bind user")
			return
		}
		service.UserServiceImpl.CreateUser(user)
		context.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": map[string]string{
				"id": user.ID,
			},
		})
	}
}

func GetUserFun() func(context *gin.Context) {
	return func(context *gin.Context) {
		user := new(model.User)
		err := context.ShouldBindUri(user)
		if err != nil {
			log.Printf("err	%s \n", err)
		}
		service.UserServiceImpl.GetUser(user)
		if err != nil {
			log.Printf("err	%s \n", err)
		}
		context.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": user,
		})
	}
}

func DeleteUserFun() func(context *gin.Context) {
	return func(context *gin.Context) {
		user := &model.User{}
		err := context.ShouldBindUri(user)
		if err != nil {
			log.Fatal("err", err)
		}
		service.UserServiceImpl.DeleteUser(user.ID)
		context.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": map[string]string{
				"id": user.ID,
			},
			"msg": "success",
		})
	}
}

func UpdateUserFun() func(context *gin.Context) {
	return func(context *gin.Context) {
		user := &model.User{}
		err := context.Bind(user)
		if err != nil {
			log.Println("error", err)
		}
		service.UserServiceImpl.UpdateUser(user)
		context.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": map[string]string{
				"id":   user.ID,
				"name": user.Name,
			},
			"msg": "success",
		})

	}
}

func BatchCreateUser() {

}

func BatchUpdateUser() {

}

func BatchGetUser() {

}

func BatchDeleteUser() {

}

func Page() {

}
