# micro-leaderboard

Leader board microservice written in Go/Gin using redis

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

### Mac/Linux

```
./micro-leaderboard
```

### Docker

```
docker pull maguec/micro-leaderboard:latest
docker run -i -t -p 8080:8080 maguec/micro-leaderboard
```

## Testing

run either the docker container or the raw application binary

```
curl http://localhost:8080/health
```

---
Copyright © 2018, Chris Mague
