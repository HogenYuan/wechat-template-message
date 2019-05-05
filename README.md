# wechat_message(公众号模板消息接口)


## Introduction
- 开发目的：用于弥补PHP群发模板消息的性能硬伤。
- 利用golang能达到不可思议的速度。针对golang的结构体传输符合微信接口的数据。


## Installation

```bash
	go build main.go
	go run wechat_mesage
```

## Run
```php
	/**
	 * 接收数据
	 * api:"post:localhost:7767/getMessage"
	 */
	
	example:{
		openid:['111','fd2dw','gree']，		//openid数组,需提前切好
		ac_token : string,					//公众号access_token
		mess_type: int,						// 0:模板消息，1:文字消息，2:图文消息

		//模板消息
		template_id:string,					//模板id
		url:"http://www.baidu.com",			//可选,跳转url	
		appid:string,						//可选,小程序appid
		pagepath:string,					//可选,小程序page
		example:{
			first:{
				value:"你好",
				color:"#e6212a"
			},
			keyword1:{
				value:"你好",
				color:"#e6212a"				
			}
		}

		//文字消息
		content:string

		//图文消息
		title:"标题",
		description:"图文描述",
		picurl:"http://www.baidu.com",		//图文链接
		url:"http://www.baidu.com"			//可选,跳转链接
	}
	 
	return [
		'total', 							//发送总人数
		'suc', 								//成功人数
	];
```

```golang
	/**
	 * 停止任务
	 */
	api:"get:localhost:7767/stop"

```

```golang
	/**
	 * 查看任务PID
	 */
	api:"get:localhost:7767/pid"

```
