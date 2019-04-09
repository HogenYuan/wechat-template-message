package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/braintree/manners"
	"github.com/gin-gonic/gin"
	// "github.com/goinggo/mapstructure"
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
	Touser  string  `json:"touser"`
	Msgtype string  `json:"msgtype"`
	News    NewsMsg `json:"news"`
}
type NewsMsg struct {
	Articles [1]ArticlesMsg `json:"articles"`
}
type ArticlesMsg struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
	Picurl      string `json:"picurl"`
}
type TemplateMsg struct {
	Touser      string                  `json:"touser"`                //接收者的OpenID
	Template_id string                  `json:"template_id"`           //模板消息ID
	Url         string                  `json:"url"`                   //点击后跳转链接
	Miniprogram MiniprogramMsg          `json:"miniprogram,omitempty"` //点击跳转小程序
	Data        map[string]*KeyWordData `json:"data"`
}
type MiniprogramMsg struct {
	Appid    string `json:"appid"`
	Pagepath string `json:"pagepath"`
}

type DataMsg struct {
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
		if c.Request.Form == nil {
			c.Request.ParseMultipartForm(32 << 20)
		}
		// for k, v := range c.Request.Form {
		// 	fmt.Println(k, v)
		// }
		ac_token := c.PostForm("ac_token")
		openid := c.PostForm("openid")
		mess_type := c.PostForm("mess_type")
		var msg interface{}
		post_url := ""
		if mess_type == "1" {
			//文字消息
			post_url = "https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=" + ac_token
			msg = &Message{
				Touser:  openid,
				Msgtype: "text",
				Text:    TextMsg{Content: c.DefaultPostForm("content", "")},
			}
		} else if mess_type == "2" {
			//图文消息
			post_url = "https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=" + ac_token
			var msgs [1]ArticlesMsg
			msgs[0] = ArticlesMsg{
				Title:       c.DefaultPostForm("title", ""),
				Description: c.DefaultPostForm("description", ""),
				Url:         c.DefaultPostForm("url", ""),
				Picurl:      c.DefaultPostForm("picurl", ""),
			}
			msg = &PicMessage{
				Touser:  openid,
				Msgtype: "news",
				News: NewsMsg{
					Articles: msgs,
				},
			}
		} else if mess_type == "0" {
			msg = 0
			tempMsg_json := c.PostForm("example")
			map1 := make(map[string]*KeyWordData)
			err := json.Unmarshal([]byte(tempMsg_json), &map1)
			if err != nil {
				fmt.Println("error", err)
			}
			fmt.Printf("dataMsg:%+v\n", map1)
			// var dataMsg map[string]*KeyWordData
			// for k, v := range map1 {
			// 	fmt.Printf("k值:%s,v值:%s\n", k, v)
			// 	// keyword := KeyWordData{v["Value"], v["Color"]}
			// 	keyword := KeyWordData{}
			// 	fmt.Printf("vvvv值%s\n", v["value"])

			// 	// err := mapstructure.Decode(v, &keyword) //map转struct
			// 	// if err != nil {
			// 	// 	fmt.Println(err)
			// 	// }
			// 	fmt.Printf("keyword值%+v\n", keyword)
			// }
			//模板消息
			post_url = "https://api.weixin.qq.com/cgi-bin/message/wxopen/template/send?access_token=" + ac_token
			msg = &TemplateMsg{
				Touser:      openid,
				Template_id: c.DefaultPostForm("template_id", ""),
				Url:         c.DefaultPostForm("url", ""),
				Miniprogram: MiniprogramMsg{
					Appid:    c.DefaultPostForm("appid", ""),
					Pagepath: c.DefaultPostForm("pagepath", ""),
				},
				Data: map1,
			}
		} else {
			msg = 0
		}
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

	pid := fmt.Sprint(os.Getpid())
	ioutil.WriteFile("./pid", []byte(pid), 0666)

	manners.ListenAndServe(":7767", r)

}
