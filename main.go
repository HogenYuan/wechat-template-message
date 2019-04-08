package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/braintree/manners"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	// "net/url"
	"os"
)

type Message struct {
	Touser  string  `json:"touser"`
	Msgtype string  `json:"msgtype"`
	Text    TextMsg `json:"text"`
}
type TextMsg struct {
	Content string `json:"content"`
}
type PicMessage struct {
	Touser  string `json:"touser"`
	Msgtype string `json:"msgtype"`
	Image   PicMsg `json:"media_id"`
}
type PicMsg struct {
	Media_id string `json:"media_id"`
}
type TemplateMsg struct {
	Touser      string                   `json:"touser"`      //接收者的OpenID
	TemplateID  string                   `json:"template_id"` //模板消息ID
	URL         string                   `json:"url"`         //点击后跳转链接
	Miniprogram Miniprogram              `json:"miniprogram"` //点击跳转小程序
	Data        map[string]*TemplateData `json:"data"`
}
type Miniprogram struct {
	AppID    string `json:"appid"`
	Pagepath string `json:"pagepath"`
}

type TemplateData struct {
	First    KeyWordData `json:"first,omitempty"`
	Keyword1 KeyWordData `json:"keyword1,omitempty"`
	Keyword2 KeyWordData `json:"keyword2,omitempty"`
	Keyword3 KeyWordData `json:"keyword3,omitempty"`
	Keyword4 KeyWordData `json:"keyword4,omitempty"`
	Keyword5 KeyWordData `json:"keyword5,omitempty"`
}

type KeyWordData struct {
	Value string `json:"value"`
	Color string `json:"color,omitempty"`
}

func main() {
	r := gin.Default()

	r.POST("/getMessage/", func(c *gin.Context) {
		// log_id := c.Param("log_id")
		if c.Request.Form == nil {
			c.Request.ParseMultipartForm(32 << 20)
		}
		for k, v := range c.Request.Form {
			fmt.Println(k, v)
		}
		ac_token := c.PostForm("ac_token")
		openid := c.PostForm("openid")
		// log := make(map[string]interface{})
		// err := json.Unmarshal([]byte(logJson), &log)
		// fmt.Printf("读取转换数据 %v \n", log)
		// if err != nil {
		// fmt.Println("Can't decode json message", err)
		// }
		// fmt.Println("获取tempid", log["temp"].(map[string]interface{})["example"].([]interface{})[3].(map[string]interface{})["keyword"])
		// fmt.Printf("看一哈 %s \n", log["uniacid"])

		// post_url := "https://api.weixin.qq.com/cgi-bin/message/wxopen/template/send?access_token=" + ac_token

		// mess_type := c.PostForm("mess_type")
		//格式转换
		post_url := "https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=" + ac_token
		msg := &Message{
			Touser:  openid,
			Msgtype: "text",
			Text:    TextMsg{Content: c.DefaultPostForm("content", "")},
		}
		// msg := struct{}{}
		// if mess_type == "1" {
		// 	//文字消息
		// 	post_url = "https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=" + ac_token
		// 	msg = &Message{
		// 		Touser:  openid,
		// 		Msgtype: "text",
		// 		Text:    TextMsg{Content: content},
		// 	}
		// } else if mess_type == "2" {
		// 	msg = 2
		// 	//图文消息
		// 	// post_url := "https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=" + ac_token
		// 	// var msg = &PicMessage{
		// 	// 	Touser:  openid,
		// 	// 	Msgtype: "image",
		// 	// 	Image:    PicMsg{Media_id:['media_id']},
		// 	// }
		// } else if mess_type == "0" {
		// 	msg = 1
		// 	//模板消息
		// 	// post_url := "https://api.weixin.qq.com/cgi-bin/message/wxopen/template/send?access_token=" + ac_token
		// 	// var msg = &TemplateMsg{
		// 	// 	Touser: openid,
		// 	// 	Template_id:log["tempid"],
		// 	// }
		// } else {
		// 	msg = 0
		// }
		fmt.Printf("msg:%+v\n", msg)

		body, err := json.MarshalIndent(msg, " ", "  ") //struct转->返回[]byte字符串
		if err != nil {
			fmt.Println("json转换错误", err)
		} else {
			fmt.Printf("转换%+v\n", body)
			fmt.Printf("转换str%s\n", string(body))
		}
		//发送请求
		req, err := http.NewRequest("POST", post_url, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json;encoding=utf-8")
		client := &http.Client{}
		res, err := client.Do(req)
		//解析数据
		if err != nil {
			fmt.Printf("请求失败%v\n", err)
			return
		} else {
			bts, err := ioutil.ReadAll(res.Body)
			if err != nil {
				fmt.Printf("错误:读取body%v\n", err)
				return
			} else {
				fmt.Printf("解析结果%v\n", string(bts))
			}
		}
		defer res.Body.Close()
	})

	//获取pid
	r.GET("/pid", func(c *gin.Context) {
		c.String(200, "PID:  %d", os.Getpid())
	})

	//停止任务
	r.GET("/stop", func(c *gin.Context) {
		c.String(200, "STOP success")
		manners.Close()
	})

	manners.ListenAndServe(":7767", r)
}
