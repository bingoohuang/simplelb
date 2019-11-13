# simplelb

simplelb is the simplest Load Balancer ever created.

It uses RoundRobin algorithm to send requests into set of backends and support
retries too.

It also performs active cleaning and passive recovery for unhealthy backends.

Since its simple it assume if / is reachable for any host its available


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
