# term-frequency
Sample term-frequency service using Golang and Redis

## API ENDPOINTS

### Insert Query
- Path : `/Insert`
- Method: `POST`
- Params: `query`
- Response: `200`

### Get Report
- Path : `/GetReport`
- Method: `GET`
- Params: `n, t`
- Response: `200`

## Required Packages
- Database
    * [Redigo](https://github.com/gomodule/redigo)
- Routing
    * [gin](https://github.com/gin-gonic/gin)

## Quick Run Project

```
git clone https://github.com/elmiramakvand/term-frequency.git

docker-compose up
```
