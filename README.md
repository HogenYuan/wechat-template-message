#wechat_message(模板消息接口)
---
##简介
开发目的：用于弥补PHP群发模板消息的性能硬伤。
利用golang能达到不可思议的速度。针对golang的结构体传输数据。

---
##编译
//需先安装go环境
cmd
(```)
	//项目根目录
	go build main.go
	go run wechat_mesage
(```)

##运行
(```)
	//接收
	api:"localhost:7767/getMessage"
	
	//停止
	api:"localhost:7767/stop"
(```)
