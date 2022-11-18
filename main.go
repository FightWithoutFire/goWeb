package main

import (
	"goWeb/middleware"
)

func main() {
	defer func() {
		middleware.RedisClient.Close()
	}()
	r := middleware.SetUp()
	r.Run(":8080")
}
