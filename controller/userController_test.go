package controller

import (
	"encoding/json"
	"fmt"
	"goWeb/middleware"
	"goWeb/model"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestGetPostgreCon(t *testing.T) {
	r := middleware.SetUp()
	w := httptest.NewRecorder()
	name := randStr(7)
	data, err := json.Marshal(&model.User{Name: name})
	if err != nil {
		t.Error("error", err)
	}
	log.Println(string(data))
	request := httptest.NewRequest(http.MethodPost, AddUserUri, strings.NewReader(string(data)))
	request.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, request)
	resbyte, err := io.ReadAll(w.Body)
	if err != nil {
		t.Error("error", err.Error())
	}
	log.Println(string(resbyte))

}

func randStr(len int) string {
	rand.Seed(time.Now().Unix())
	name := ""
	charset := "abcdefghijklmnopqrstuvwxyz"

	for i := 0; i < len; i++ {
		name += string(charset[rand.Intn(26)])
	}
	return name
}

func TestGetUser(t *testing.T) {
	r := middleware.SetUp()
	w := httptest.NewRecorder()
	id := "3"

	request := httptest.NewRequest(http.MethodGet, fmtUsersUri(id), nil)
	r.ServeHTTP(w, request)
	resbyte, err := io.ReadAll(w.Body)
	if err != nil {
		t.Error("error", err.Error())
	}
	log.Println(string(resbyte))

}

func fmtUsersUri(id string) string {
	return fmt.Sprintf("/users/%s", id)
}

func TestValidate(t *testing.T) {
	req := model.RequestParam{
		"abde",
		"qwer",
		110,
		"953@mail.com",
	}
	err := model.ValidateFuc(req)

	if err != nil {
		t.Error("error", err)
	}

	vt := reflect.TypeOf(req)
	vv := reflect.ValueOf(req)

	for i := 0; i < vv.NumField(); i++ {
		fieldValue := vv.Field(i)
		fieldType := vt.Field(i)
		kind := fieldValue.Kind()
		switch kind {
		case reflect.Int:
			fmt.Printf("%v  %v %v %v \n", fieldType.Name, fieldType.Tag, fieldType.Type, fieldValue.Int())
		case reflect.String:
			fmt.Printf("%v  %v %v %s \n", fieldType.Name, fieldType.Tag, fieldType.Type, fieldValue.String())
		case reflect.Struct:

		}
	}

}

func TestDeleteUser(t *testing.T) {

	s := middleware.SetUp()

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, fmtUsersUri("13"), nil)
	s.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Error("error delete")
	}
}

func TestUpdateUser(t *testing.T) {
	err, w := updateTest()

	if w.Code != http.StatusOK {
		t.Error("error update", err)
	}

}

func updateTest() (error, *httptest.ResponseRecorder) {
	s := middleware.SetUp()
	user, data, err := getUpdateUserJson("3")
	r := httptest.NewRequest(http.MethodPut, fmtUsersUri(user.ID), strings.NewReader(string(data)))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.ServeHTTP(w, r)
	return err, w
}

func getUpdateUserJson(id string) (*model.User, []byte, error) {
	user := &model.User{ID: id, Name: randStr(7)}
	data, err := json.Marshal(user)
	log.Println(string(data))
	return user, data, err
}

func TestUpdateWithoutId(t *testing.T) {
	s := middleware.SetUp()
	data, err := json.Marshal(&model.User{Name: randStr(7)})
	if err != nil {
		t.Error("error", err)
	}
	log.Println(string(data))
	r := httptest.NewRequest(http.MethodPut, "/users", strings.NewReader(string(data)))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Error("error delete")
	}

}

func BenchmarkUpdateUserFun(b *testing.B) {
	for i := b.N; i < 0; i++ {
		err, recorder := updateTest()
		if err != nil {
			b.Error(err, recorder)
		}
	}
}
