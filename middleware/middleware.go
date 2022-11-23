package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"goWeb/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"io/fs"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
)

var DbClient *gorm.DB
var RedisClient *redis.Client
var GinEngine *gin.Engine
var ProjectName string

func initPostgreCon(config *SqlConfig) (db *gorm.DB, err error) {
	host := config.Ip
	p := config.Port
	user := config.User
	pass := config.Password
	database := config.Database
	port, _ := strconv.Atoi(p)
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, database)
	return gorm.Open(postgres.Open(conn), &gorm.Config{})
}

func init() {
	ProjectName = "goWeb"

	config, err := GetConfig()
	if err != nil {
		log.Fatalln("init fail")
	}
	con, err := initPostgreCon(config.SqlConfig)
	if err != nil {
		log.Fatalln("error", err)
	}
	DbClient = con
	initRedis(config.RedisConfig)
	initGinEngine()
	SetUp()
	err = DbClient.Migrator().AutoMigrate(&model.User{})
	if err != nil {
		log.Fatalln("AutoMigrate error", err)
	}
}

func initRedis(config *RedisConfig) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Ip, config.Port),
		Password: config.Password,
		DB:       config.Database,
	})
}

func initGinEngine() {
	GinEngine = gin.Default()
	GinEngine.Use(gin.Logger())
	GinEngine.Use(gin.Recovery())
}
func SetUp() *gin.Engine {
	GinEngine.Use(func(context *gin.Context) {
		log.Println("start")
		//m := &model.User{}
		//cCp := context.Copy()
		//cCp.BindJSON(m)
		//hello := cCp.GetHeader("hello")
		//log.Println("hello:", hello)
		//if m.ID == "0" {
		//	log.Println("hello 0")
		//	context.AbortWithStatusJSON(500, gin.H{
		//		"mgs": "id not 0",
		//	})
		//} else {
		//
		//}
		context.Next()
		log.Println("end")
	})

	return GinEngine
}

type Config struct {
	RedisConfig *RedisConfig `json:"redis"`
	SqlConfig   *SqlConfig   `json:"sql"`
}

type RedisConfig struct {
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	Password string `json:"password"`
	Database int    `json:"database"`
}

type SqlConfig struct {
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

func GetConfig() (*Config, error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalln("error read file")
	}
	split := strings.Split(dir, "goWeb")
	configDir := split[0] + ProjectName + "/config.json"
	log.Println(configDir)

	file, err := os.OpenFile(configDir, syscall.O_RDONLY, fs.ModePerm)
	if err != nil {
		log.Fatalln("error read file")
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatalln("error read file")
	}

	config := &Config{
		RedisConfig: &RedisConfig{},
		SqlConfig:   &SqlConfig{},
	}
	err = json.Unmarshal(bytes, config)
	fmt.Println("config: ", config)
	if err != nil {
		log.Fatalln("error json translate")
	}

	return config, err

}
