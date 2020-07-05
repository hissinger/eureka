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


## 테스트 라이브러리
- github.com/hudl/fargo
- github.com/xuanbo/eureka-client
- github.com/pineda89/golang-springboot/eureka
- github.com/ArthurHlt/go-eureka-client
- github.com/spectre013/fairway
- github.com/HikoQiu/go-eureka-client/eureka (instance 정보에 meta 필드가 없음)
- 직접 REST로 구현


## 결과
- https://github.com/phantasmicmeans/spring-cloud-netflix-eureka-tutorial 를 참고해서 직접 eureka 서버를 실행하면 모든 라이브러리가 정상 작동하는 것을 확인함.
- 처음에는 docekr(https://hub.docker.com/r/springcloud/eureka)로 eureka 서버를 실행했는데 아래와 같은 이슈가 있었음.(삽질)
- docker(https://hub.docker.com/r/springcloud/eureka)로 eureka 서버를 실행했을 때 이슈
  - degegister시에도 eureka 서버에서 삭제되지 않음.
  - instance id 필드를 사용하면 hearbeat시에 404에러 발생
  - github.com/ArthurHlt/go-eureka-client : port enable 필드가 bool로 정의 되어 있어서 에러 발생
  - github.com/spectre013/fairway : port enable 필드가 bool로 정의 되어 있어서 에러 발생
  - github.com/HikoQiu/go-eureka-client/eureka : instance 정보에 meta 필드가 없음.


## References
 - https://github.com/phantasmicmeans/spring-cloud-netflix-eureka-tutorial
 - https://libs.garden/go/search?q=eureka&page=1 golang의 eureka client 구현 라이브러리를 정리한 사이트
