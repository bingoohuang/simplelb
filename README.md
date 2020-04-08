# simplelb

simplelb is the simplest Load Balancer ever created.

It uses RoundRobin algorithm to send requests into set of backends and support
retries too.

It also performs active cleaning and passive recovery for unhealthy backends.

Since its simple it assume if / is reachable for any host its available

已经实现的功能：

1. 后端的RR负载均衡算法（见下面示意图）
1. 后端每20秒的健康检查
1. 单个后端的3次重试(retries)和整个后端集群的3次尝试(attempts)

![image](https://user-images.githubusercontent.com/1940588/68740133-5a7e3a80-0625-11ea-8faa-dcca0df04b7b.png)

![image](https://user-images.githubusercontent.com/1940588/68740165-6e29a100-0625-11ea-8ddd-3854735e0ae1.png)

## How to build

```bash
 go fmt ./...&&goimports -w .&&golint ./...&&golangci-lint run --enable-all&& go install ./...
```

## How to use

```bash
$ simplelb -h
Usage of simplelb:
  -b string
    	Load balanced backends, use , to separate
  -p int
    	Port to serve (default 3030)
```

Example:

```bash
$ simplelb -b http://127.0.0.1:9001,http://127.0.0.1:9002
2019/11/13 13:22:07 Configured server: http://127.0.0.1:9001
2019/11/13 13:22:07 Configured server: http://127.0.0.1:9002
2019/11/13 13:22:07 Load Balancer started at :3030
```

## Thanks

1. [Let's Create a Simple Load Balancer With Go](https://kasvith.github.io/posts/lets-create-a-simple-lb-go/)
1. [Reverse Proxy in Go](https://blog.charmes.net/post/reverse-proxy-go/)
1. [reverse http / websocket proxy based on fasthttp](https://github.com/yeqown/fasthttp-reverse-proxy)
