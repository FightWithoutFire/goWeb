package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"goWeb/middleware"
	"goWeb/model"
	"goWeb/service"
	"log"
	"net/http"
)

const (
	AddUserEndpoint    = "/users"
	RestUserEndpoint   = "/users/:id"
	BookableEndpoint   = "/bookable"
	RedirectEndpoint   = "/redirect"
	SetCookie          = "/setCookie"
	MultipartEndpoint  = "/multipart"
	SingleFileEndpoint = "/single"
	PageEndpoint       = "/page"
)

func init() {
	middleware.GinEngine.POST(AddUserEndpoint, createUserFun)
	middleware.GinEngine.PUT(RestUserEndpoint, updateUserFun)
	middleware.GinEngine.GET(RestUserEndpoint, getUserFun)
	middleware.GinEngine.DELETE(RestUserEndpoint, deleteUserFun)
	middleware.GinEngine.GET(BookableEndpoint, getBookable)
	middleware.GinEngine.GET(RedirectEndpoint, redirectFun)
	middleware.GinEngine.GET(SetCookie, setCookieFun)
	middleware.GinEngine.POST(MultipartEndpoint, multipartEndpointFun)
	middleware.GinEngine.POST(SingleFileEndpoint, singleFileEndpoint)
	middleware.GinEngine.GET(PageEndpoint, pageFun)
}

func pageFun(context *gin.Context) {
	var users = new([]*model.User)
	service.UserServiceImpl.PageUser(users)
	context.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": *users,
	})
}

func multipartEndpointFun(context *gin.Context) {
	form, err := context.MultipartForm()
	if err != nil {
		context.AbortWithStatusJSON(500, gin.H{
			"msg": "get form error",
		})
	}
	files := form.File["upload[]"]
	for _, file := range files {
		log.Println(file.Filename)

	}
	context.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))

}

func singleFileEndpoint(context *gin.Context) {
	file, err := context.FormFile("upload")
	if err != nil {
		context.AbortWithStatusJSON(500, gin.H{
			"msg": "get form error",
		})
	}
	log.Println(file.Filename)

	context.String(http.StatusOK, fmt.Sprintf("%s %d size file uploaded!", file.Filename, file.Size))

}

func setCookieFun(context *gin.Context) {

}

func redirectFun(context *gin.Context) {
	context.Request.URL.Path = "/users/3"
	middleware.GinEngine.HandleContext(context)

}

func getBookable(context *gin.Context) {
	m := &model.User{}
	if err := context.ShouldBindWith(&m, binding.Query); err == nil {
		context.JSON(http.StatusOK, gin.H{"message": "Booking dates are valid!"})
	} else {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func createUserFun(context *gin.Context) {

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

func getUserFun(context *gin.Context) {
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
		"code":   200,
		"data":   user,
		"extend": "<!@啊啊>",
	})
}

func deleteUserFun(context *gin.Context) {
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

func updateUserFun(context *gin.Context) {
	log.Println("update in user controller")
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
