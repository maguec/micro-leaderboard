# micro-leaderboard

Leader board microservice written in Go/Gin using redis

## Usage

This is meant to be a fast and efficient microservice to create a leaderboard.


### Memory Usage
The leaderboard is stored in Redis using a sorted set which is very fast an easy on memory

You can see that 995K entries only use about about 20 mb of memory
```
127.0.0.1:6379>  zcount testset -inf +inf
(integer) 995128
127.0.0.1:6379> debug object testset
Value at:0x7f993c466850 refcount:1 encoding:skiplist serializedlength:19791998 lru:15849160 lru_seconds_idle:8
```

### Latency

Using the [wrk](https://github.com/wg/wrk) benchmarking tool and a script to randomly generate entries on an already populated sorted set we get the following results *without any OS or redis tuning* on my 4 core desktop.

```
$ wrk -c 100 -t 20 -d 30s -s testing/random-user.lua http://localhost:8080
Running 30s test @ http://localhost:8080
  20 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     5.07ms    1.99ms  38.57ms   79.54%
    Req/Sec     1.00k   101.67     1.17k    90.58%
  597811 requests in 30.03s, 109.20MB read
Requests/sec:  19910.25
Transfer/sec:      3.64MB
```

For information on [tuning redis for prodution](http://shokunin.co/blog/2014/11/11/operational_redis.html) see my guide that while dated is still accurate.

## Building

### Mac/Linux

0) set GOROOT environment variable
1) Install Go and Make
2) make

### Docker

0) set GOROOT environment variable
1) Install Docker, Go and Make
2) make docker


## Running

### Routes
```
GET   /health                   --> Health Check
GET   /                         --> Fake root handler
GET   /inc/:set/:member         --> Increment the leaderboard for a member
GET   /inc/:set/:member/:count  --> Increment the leaderboard for a member by a specific count
GET   /member/:set/:member      --> Get Information for a member of the leaderboard
GET   /board/:set               --> Get the Full leaderboard returned - not really advisable
GET   /board/:set/:count        --> Get the Top X entries from the leaderboard
```

### Example Usage

#### Load Data:
```
curl -s localhost:8080/inc/myleaderboard/Reiko/7
curl -s localhost:8080/inc/myleaderboard/Alex/2
curl -s localhost:8080/inc/myleaderboard/Chris/3
curl -s localhost:8080/inc/myleaderboard/Sean
```

#### Get Alex's ranking

```
$ curl -s localhost:8080/member/myleaderboard/Alex |jq
{
  "board": "myleaderboard",
  "member": "Alex",
  "rank": 1,
  "score": 2
}
```

### Get the top 2 

```
$ curl -s localhost:8080/board/myleaderboard/2 |jq
{
  "board": "myleaderboard",
  "leaders": [
    {
      "Score": 7,
      "Member": "Reiko"
    },
    {
      "Score": 3,
      "Member": "Chris"
    }
  ]
}
```


### Mac/Linux

```
./micro-leaderboard
```

### Docker

```
docker pull maguec/micro-leaderboard:latest
docker run --rm -p 6379:6379 --name myredis redis
docker run --rm -i -t -p 8080:8080 -e REDIS_HOST=redis -e REDIS_PORT=6379 --link myredis:redis maguec/micro-leaderboard
```

## Testing

run either the docker container or the raw application binary

```
curl http://localhost:8080/health
```

---
Copyright © 2018, Chris Mague
