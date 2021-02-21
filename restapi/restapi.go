package restapi

import (
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

func RunApi(redisPool *redis.Pool) *gin.Engine {
	r := gin.Default()
	RunApiOnRouter(r, redisPool)
	return r
}

func RunApiOnRouter(r *gin.Engine, redisPool *redis.Pool) {
	Handler := NewCacheHandler(redisPool)
	cacheQueryGroup := r.Group("/api/cache-query")
	{
		cacheQueryGroup.GET("get-report", Handler.GetReport)
		cacheQueryGroup.POST("insert", Handler.Insert)
	}
}
