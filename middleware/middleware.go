package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"io/fs"
	"log"
	"os"
	"strconv"
	"syscall"
)

var DbClient *gorm.DB
var RedisClient *redis.Client
var GinEngine *gin.Engine

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
	config, err := getConfig()
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
}
func SetUp() *gin.Engine {
	GinEngine.Use(func(context *gin.Context) {
		log.Println("start")
		context.Handler()
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

func getConfig() (*Config, error) {
	file, err := os.OpenFile("../config.json", syscall.O_RDONLY, fs.ModePerm)
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
