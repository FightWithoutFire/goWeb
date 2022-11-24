package middleware

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
var PublicKey *rsa.PublicKey
var PrivateKey *rsa.PrivateKey
var ProjectDir string
var DisableCache bool

func initPostgreCon(config *SqlConfig) (db *gorm.DB, err error) {
	host := config.Ip
	p := config.Port
	user := config.User
	pass := config.Password
	database := config.Database
	port, _ := strconv.Atoi(p)
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, database)
	return gorm.Open(postgres.Open(conn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
}

func init() {
	ProjectName = "goWeb"
	DisableCache = true
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
	SetUpEncryptKey()
	log.Println("finish init")
}

func SetUpEncryptKey() {
	GenerateKey(2048)

	privateKey := ProjectDir + "/private.key"
	publicKey := ProjectDir + "/public.key"
	privateKeyFile, err := os.ReadFile(privateKey)
	if err != nil {
		log.Fatalln("not read private key")
	}

	publicKeyFile, err := os.ReadFile(publicKey)
	if err != nil {
		log.Fatalln("not read public key")
	}

	publicKeyPem, _ := pem.Decode(publicKeyFile)
	PublicKeyTemp, err := x509.ParsePKIXPublicKey(publicKeyPem.Bytes)
	PublicKey = PublicKeyTemp.(*rsa.PublicKey)
	if err != nil {
		log.Fatalln("fail to parse Public Key:", err)
	}

	privateKeyPem, _ := pem.Decode(privateKeyFile)
	PrivateKeyTemp, err := x509.ParsePKCS1PrivateKey(privateKeyPem.Bytes)
	PrivateKey = PrivateKeyTemp
	if err != nil {
		log.Fatalln("fail to parse Private Key", err)
	}
}

func GenerateKey(bits int) {

	privateKeyPath := ProjectDir + "/private.key"
	privateKeyFile, err := os.ReadFile(privateKeyPath)
	if os.IsNotExist(err) {

		log.Println("helo")
		err := generatePrivateKey(bits, privateKeyPath)
		if err != nil {
			log.Fatalln("generate private key error, ", err)
		}

		privateKeyFile, err = os.ReadFile(privateKeyPath)
		if err != nil {
			log.Fatalln("read after generate private key error, ", err)
		}
	}

	bock, _ := pem.Decode(privateKeyFile)
	privateKey, err := x509.ParsePKCS1PrivateKey(bock.Bytes)
	if err != nil {
		log.Fatalln("parse private key error:", err)
	}
	err = generatePublicKey(privateKey)

	if err != nil {
		log.Fatalln("generate public key error, ", err)
	}

}

func generatePublicKey(privateKey *rsa.PrivateKey) error {
	/*
		生成公钥
	*/
	publicKeyPath := ProjectDir + "/public.key"
	_, err := os.ReadFile(publicKeyPath)
	if os.IsNotExist(err) {
		publicKey := privateKey.PublicKey
		publicStream, err := x509.MarshalPKIXPublicKey(&publicKey)
		//publicStream:=x509.MarshalPKCS1PublicKey(&publicKey)
		block2 := pem.Block{
			Type:  "public key",
			Bytes: publicStream,
		}
		fPublic, err := os.Create(publicKeyPath)
		if err != nil {
			return err
		}
		defer fPublic.Close()
		pem.Encode(fPublic, &block2)
	}
	return nil
}

func generatePrivateKey(bits int, privateKeyPath string) error {
	/*
		生成私钥
	*/
	//1、使用RSA中的GenerateKey方法生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	//2、通过X509标准将得到的RAS私钥序列化为：ASN.1 的DER编码字符串
	privateStream := x509.MarshalPKCS1PrivateKey(privateKey)
	//3、将私钥字符串设置到pem格式块中
	block1 := pem.Block{
		Type:  "private key",
		Bytes: privateStream,
	}
	//4、通过pem将设置的数据进行编码，并写入磁盘文件
	fPrivate, err := os.Create(privateKeyPath)
	if err != nil {
		return err
	}
	err = pem.Encode(fPrivate, &block1)
	defer fPrivate.Close()
	if err != nil {
		return err
	}
	return nil
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
	ProjectDir = split[0] + ProjectName
	configDir := ProjectDir + "/config.json"
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
