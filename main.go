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
	// "runtime"
	// "strconv"
	"sync"
	"time"
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
	Url         string                  `json:"url,omitempty"`         //点击后跳转链接
	Miniprogram MiniprogramMsg          `json:"miniprogram,omitempty"` //点击跳转小程序
	Data        map[string]*KeyWordData `json:"data"`
}
type MiniprogramMsg struct {
	Appid    string `json:"appid,omitempty"`
	Pagepath string `json:"pagepath,omitempty"`
}

type KeyWordData struct {
	Value string `json:"value"`
	Color string `json:"color,omitempty"`
}

func main() {
	r := gin.Default()

	r.POST("/getMessage/", func(c *gin.Context) {
		// if c.Request.Form == nil {
		// 	c.Request.ParseMultipartForm(32 << 20)
		// }
		// for k, v := range c.Request.Form {
		// 	fmt.Println(k, v)
		// }
		ac_token := c.PostForm("ac_token")
		mess_type := c.PostForm("mess_type")

		openid_100_json := c.PostForm("openid_100")
		var openid_100 map[string]string
		err := json.Unmarshal([]byte(openid_100_json), &openid_100)
		if err != nil {
			fmt.Println("openid_100 err: ", err)
			return
		}
		var msg interface{}
		var total = len(openid_100)
		var suc = 0
		post_url := ""
		var start_time = time.Now().Unix()
		//协程分发
		wg := sync.WaitGroup{}
		// max_process := c.DefaultPostForm("max_process", "0")
		// if max_process != "0" {
		// 	process, _ := strconv.Atoi(max_process)
		// }
		// runtime.GOMAXPROCS(runtime.NumCPU() - 2)

		for openid, _ := range openid_100 {
			wg.Add(1)
			go func(openid string) {
				defer func() {
					if err := recover(); err != nil {
						fmt.Println("don't worry, I can take care of myself.panic:", err)
					}
				}()
				content := c.DefaultPostForm("content", "")
				if mess_type == "1" {
					//文字消息
					post_url = "https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=" + ac_token
					msg = &Message{
						Touser:  openid,
						Msgtype: "text",
						Text:    TextMsg{Content: content},
					}
				} else if mess_type == "2" {
					//图文消息
					title := c.DefaultPostForm("title", "")
					description := c.DefaultPostForm("description", "")
					url := c.DefaultPostForm("url", "")
					picurl := c.DefaultPostForm("picurl", "")
					post_url = "https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=" + ac_token
					var msgs [1]ArticlesMsg
					msgs[0] = ArticlesMsg{
						Title:       title,
						Description: description,
						Url:         url,
						Picurl:      picurl,
					}
					msg = &PicMessage{
						Touser:  openid,
						Msgtype: "news",
						News: NewsMsg{
							Articles: msgs,
						},
					}
				} else if mess_type == "0" {
					tempMsg_json := c.PostForm("example")
					dataMsg := make(map[string]*KeyWordData)
					err := json.Unmarshal([]byte(tempMsg_json), &dataMsg)
					if err != nil {
						fmt.Println("error", err)
					}
					//模板消息
					post_url = "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=" + ac_token
					msg = &TemplateMsg{
						Touser:      openid,
						Template_id: c.DefaultPostForm("template_id", ""),
						Url:         c.DefaultPostForm("url", ""),
						Miniprogram: MiniprogramMsg{
							Appid:    c.DefaultPostForm("appid", ""),
							Pagepath: c.DefaultPostForm("pagepath", ""),
						},
						Data: dataMsg,
					}
				}

				body, err := json.MarshalIndent(msg, " ", "  ") //struct转->返回[]byte字符串
				if err != nil {
					fmt.Println("json转换错误", err)
					return
				}

				//发送请求
				req, err := http.NewRequest("POST", post_url, bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json;encoding=utf-8")
				client := &http.Client{}
				res, err := client.Do(req)
				//解析数据
				if err != nil {
					fmt.Printf("请求失败%v\n", err)
				} else {
					_, err := ioutil.ReadAll(res.Body)
					if err != nil {
					} else {
						suc++
					}
				}
				wg.Done()
				defer res.Body.Close()
			}(openid)
		}
		wg.Wait()
		core := c.DefaultPostForm("core", "0")
		fmt.Printf("线程:%s的开始时间:%d,结束时间:%d,发送人数为:%d个\n", core, start_time, time.Now().Unix(), suc)
		c.JSON(200, gin.H{
			"total": total,
			"suc":   suc,
		})
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
