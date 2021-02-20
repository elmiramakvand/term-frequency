package restapi

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/gin-gonic/gin"
)

var wg sync.WaitGroup

type RedisPool struct {
	redisPool *redis.Pool
}

func NewCacheQueryModel(redisPool *redis.Pool) *RedisPool {
	return &RedisPool{
		redisPool: redisPool,
	}
}

func (redisPool RedisPool) GetReport(c *gin.Context) {
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

	keys := getKeysForReport(timeInt)

	conn := redisPool.redisPool.Get()
	defer conn.Close()

	now := time.Now()
	keyTop := "TOP_" + now.Format("20060102") + "_" + t + "h"

	var args []interface{}
	args = append(args, keyTop)
	args = append(args, t)
	for _, k := range keys {
		args = append(args, k)
	}
	conn.Do("ZUNIONSTORE", args...)

	// get count of all keys in given time
	totalTokenCount, err := redis.Int(conn.Do("ZCOUNT", keyTop, "-inf", "+inf"))
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
	values, err := redis.Strings(conn.Do("ZREVRANGEBYSCORE", keyTop, "+inf", "-inf", "LIMIT", "0", n, "withscores"))
	if err != nil {
		fmt.Println(err)
		return
	}

	conn.Do("EXPIRE", keyTop, 1)

	// generate Csv File
	headers := []string{"term", "count"}
	generateCsvFile(headers, values, c)
	return
}

func generateCsvFile(headers []string, values []string, c *gin.Context) {
	c.Header("Content-Disposition", "attachment; filename=export.csv")
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

func (redisPool RedisPool) Insert(c *gin.Context) {
	//queryString, ok := c.GetQuery("query")
	queryStrings, ok := c.GetQueryArray("query")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "query not found!"})
		return
	}
	var tokens []string

	for _, queryString := range queryStrings {
		// standard tokenizer : This tokenizer splits the text field into tokens, treating whitespace and punctuation as delimiters. Delimiter characters are discarded
		standardTokens := standardTokenizer(queryString, ":@,-")
		tokens = append(tokens, standardTokens...)
		// keyword tokenizer : This tokenizer treats the entire text field as a single token.
		found := Find(tokens, strings.ToLower(queryString))
		if !found {
			//Value not found in slice
			tokens = append(tokens, strings.ToLower(queryString))
		}
	}

	for _, token := range tokens {
		wg.Add(1)
		go cacheTokensInRedis(token, redisPool)
	}

	wg.Wait()
	c.JSON(http.StatusOK, gin.H{"msg": "query successfully cached !"})
	return
}

func Find(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func cacheTokensInRedis(token string, redisPool RedisPool) {
	defer wg.Done()
	c := redisPool.redisPool.Get()
	now := time.Now()
	keySet := now.Format("20060102_15")
	c.Do("ZINCRBY", keySet, 1, token)
	// expire key after 168 hours or 1 week
	c.Do("EXPIRE", keySet, 604800)
	c.Close()
}

func standardTokenizer(s string, seps string) []string {
	step1 := strings.ToLower(s)
	//	var re = regexp.MustCompile(`(^\.*)| \.| *\. |@|,|-|:|\.*$`)
	var re = regexp.MustCompile(`(^\.*)| \.| *\. |@|'|\?|\(|\)|"|“|”|,|-|:|\.*$`)
	step2 := re.ReplaceAllString(step1, " ")
	return strings.Fields(step2)
}
