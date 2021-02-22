package test

import (
	"encoding/csv"
	"io"
	"net/http"
	"net/http/httptest"
	"term-frequency/config"
	"term-frequency/restapi"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	// use db no 1 for testing api
	redisPool := config.GetRedisPool(":6379", 1)
	// flush db 1 data
	conn := redisPool.Get()
	conn.Do("FlushDB")
	conn.Close()

	gin.SetMode(gin.ReleaseMode)
	r := restapi.RunApi(redisPool)
	w := httptest.NewRecorder()

	// Mock HTTP Request and it's return
	req, err := http.NewRequest("POST", "/api/cache-query/insert?query=Please, email john.doe@foo.com by 03-09, re: m37-xq.", nil)

	// make sure request was executed
	assert.NoError(t, err)

	// Serve Request and recorded data
	r.ServeHTTP(w, req)

	// Test results
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"msg":"query successfully cached !"}`, w.Body.String())

}

func TestGetReport(t *testing.T) {
	// use db no 1 for testing api
	redisPool := config.GetRedisPool(":6379", 1)

	gin.SetMode(gin.ReleaseMode)
	r := restapi.RunApi(redisPool)
	w := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/api/cache-query/get-report?n=11&t=1", nil)

	// make sure request was executed
	assert.NoError(t, err)

	// Serve Request and recorded data
	r.ServeHTTP(w, req)

	// Test results
	assert.Equal(t, 200, w.Code)

	var mockItems [][]string
	mockItems = append(mockItems, []string{"term", "count"})
	mockItems = append(mockItems, []string{"xq", "1"})
	mockItems = append(mockItems, []string{"re", "1"})
	mockItems = append(mockItems, []string{"please, email john.doe@foo.com by 03-09, re: m37-xq.", "1"})
	mockItems = append(mockItems, []string{"please", "1"})
	mockItems = append(mockItems, []string{"m37", "1"})
	mockItems = append(mockItems, []string{"john.doe", "1"})
	mockItems = append(mockItems, []string{"foo.com", "1"})
	mockItems = append(mockItems, []string{"email", "1"})
	mockItems = append(mockItems, []string{"by", "1"})
	mockItems = append(mockItems, []string{"09", "1"})
	mockItems = append(mockItems, []string{"03", "1"})

	reader := csv.NewReader(w.Body)
	var results [][]string
	for {
		// read one row from csv
		record, err := reader.Read()
		if err == io.EOF {
			//t.Fail()
			break
		}
		if err != nil {
			t.Fail()
		}

		// add record to result set
		if record != nil {
			results = append(results, record)
		}

	}
	// make sure data from csv file match with request query
	assert.Equal(t, mockItems, results)
}
