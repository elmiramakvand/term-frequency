package restapi

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"term-frequency/repository"
	cacherepo "term-frequency/repository/cacherepository"
	"term-frequency/tokenizer"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/gin-gonic/gin"
)

func NewCacheHandler(redisPool *redis.Pool) *Cache {
	return &Cache{
		repo: cacherepo.NewCacheRepository(redisPool),
	}
}

type Cache struct {
	repo repository.ICacheRepository
}

func (cache *Cache) Insert(c *gin.Context) {
	queryStrings, ok := c.GetQueryArray("query")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "query not found!"})
		return
	}

	var tokens []string

	standardTokens := tokenizer.StandardTokenizer(queryStrings)
	tokens = append(tokens, standardTokens...)

	keywordTokens := tokenizer.KeywordTokenizer(queryStrings, tokens)
	tokens = append(tokens, keywordTokens...)

	cache.repo.InsertTokens(tokens)
	c.JSON(http.StatusOK, gin.H{"msg": "query successfully cached !"})
	return
}

func (cache *Cache) GetReport(c *gin.Context) {
	t, ok := c.GetQuery("t")
	if !ok {
		// Parameter does not exist : default value of t is 1
		t = "1"
	}
	timeInt, err := strconv.Atoi(t)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter t should be number"})
		return
	}

	if timeInt > 168 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter t should be less or equal to 168"})
		return
	}

	n, ok := c.GetQuery("n")
	if !ok {
		// Parameter does not exist : default value of n is 10
		n = "10"
	}

	numberOfTokensInt, err := strconv.Atoi(n)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter n should be number"})
		return
	}

	//generate keys for the last t hour
	keys := getKeysForReport(timeInt)

	now := time.Now()
	keyTop := "TOP_" + now.Format("20060102") + "_" + t + "h"

	err = cache.repo.StoreKeyUnionOfTokens(keyTop, t, keys)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// get count of all keys in given time
	totalTokenCount, err := cache.repo.GetCountOfTokensInSortedSet(keyTop)

	if err != nil {
		fmt.Println(err)
		return
	}
	if totalTokenCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no token inserted in the last " + t + " hours"})
		return
	}

	if numberOfTokensInt > totalTokenCount {
		c.JSON(http.StatusBadRequest, gin.H{"error": strconv.Itoa(totalTokenCount) + " token found in Db.parameter n should be less or equal to " + strconv.Itoa(totalTokenCount)})
		return
	}

	//get top n token in last t hours
	values, err := cache.repo.GetTopValuesOfSortedSet(keyTop, n)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = cache.repo.ExpireKey(keyTop, 1)
	if err != nil {
		fmt.Println(err)
		return
	}

	// generate Csv File
	headers := []string{"term", "count"}
	generateCsvFile(headers, values, c)
	return
}

func generateCsvFile(headers []string, values []string, c *gin.Context) {
	c.Header("Content-Disposition", "attachment; filename=report.csv")
	c.Header("Content-Type", "text/csv")
	c.Header("Transfer-Encoding", "chunked")

	writer := csv.NewWriter(c.Writer)
	err := writer.Write(headers)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	for i := 0; i < len(values); i += 2 {
		if values[i] != " " {
			err = writer.Write([]string{values[i], values[i+1]})
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}
	}

	writer.Flush()
	return
}

func getKeysForReport(n int) []string {
	var keys []string
	now := time.Now()

	for i := n - 1; i >= 0; i-- {
		keyDate := now.Add(time.Duration(-i) * time.Hour)
		keys = append(keys, keyDate.Format("20060102_15"))
	}
	return keys
}
