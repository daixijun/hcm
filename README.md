# Helm Chart Mirror(hcm)

helm charts 同步工具

支持 本地目录、Aliyun OSS存储


## 编译

> 需要 go1.11及以上

```
$ git clone https://github.com/daixijun/hcm.git
$ cd hcm
$ go build -o hcm cmd/hcm/main.go
```

## TODO

* [ ] 存储及下载包的操作需要单独启用goroutine
