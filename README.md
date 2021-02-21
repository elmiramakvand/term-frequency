# term-frequency
Sample term-frequency service using Golang and Redis

## API ENDPOINTS

### Insert Query
- Path : `/api/cacheQuery/Insert`
- Method: `POST`
- Params: `query`
- Response: `200`

### Get Report
- Path : `/api/cacheQuery/GetReport`
- Method: `GET`
- Params: `n, t`
- Response: `200`

### URL examples
* Insert Query:
    * POST   localhost:8080/api/cache-query/insert?query=Please, email john.doe@foo.com by 03-09, re: m37-xq.
* Get Report:
    * GET    localhost:8080/api/cache-query/get-report?n=10&t=1

## Required Packages
- Database
    * [Redigo](https://github.com/gomodule/redigo)
- Routing
    * [gin](https://github.com/gin-gonic/gin)

## Run Project

```
git clone https://github.com/elmiramakvand/term-frequency.git

docker-compose up
```
