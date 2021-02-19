package main

import (
	"term-frequency/config"
	"term-frequency/restapi"
)

func main() {
	redisPool := config.GetRedisPool(":6379", 0)
	r := restapi.RunApi(redisPool)
	//running
	r.Run(":8080")

}
