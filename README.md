# Golang에서의 eureka client 

golang에서 spring cloud의 eureka 서버와의 연동을 위해 테스트



## Build

```bash
go build --race -o eureka-client main.go
```

## Run

```bash
go run --race main.go
```



## Docker

### Build
```bash
docker build --tag eureka-client -f docker/Dockerfile .
```

### Run
```bash
docker run -it --rm eureka-client
```



## 테스트 결과

golang의 eureka client 구현 라이브러리를 정리한 사이트: https://libs.garden/go/search?q=eureka&page=1



### 이슈

- DeRegister시에 status를 `DOWN`으로 변경하고 DeRegister(DELETE Method)를 수행해야 한다. (spring cloud client의 REST 패킷 분석 결과)
- instance id 필드를 사용하면 hearbeat시에 404에러 발생. spring cloud client도 동일 현상 발생됨.



### 기능 동작 확인
- github.com/hudl/fargo
- github.com/xuanbo/eureka-client
- github.com/pineda89/golang-springboot/eureka : 내부적으로 fargo 사용
- 직접 REST로 구현



### 문제 있음

- github.com/ArthurHlt/go-eureka-client : port enable 필드가 bool로 정의 되어 있어서 eureka 서버에서 에러 발생

- github.com/spectre013/fairway : port enable 필드가 bool로 정의 되어 있어서 eureka 서버에서 에러 발생
- github.com/HikoQiu/go-eureka-client/eureka : instance 정보에 meta 필드가 없음.