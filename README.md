# wechat_message(模板消息接口)


## 简介
开发目的：用于弥补PHP群发模板消息的性能硬伤。
利用golang能达到不可思议的速度。针对golang的结构体传输符合微信接口的数据。


## 编译

```bash
	go build main.go
	go run wechat_mesage
```

## 运行
```golang

	//接收
	api:"localhost:7767/getMessage"

```

```golang

	//停止
	api:"localhost:7767/stop"

```
