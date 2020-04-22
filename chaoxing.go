package main

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/valyala/fasthttp"
)

type SourceResult struct {
	Msg    interface{} `json:"msg"`
	Result int8        `json:"result"`
}

type Data struct {
	Puid       int    `bson:"puid"`
	Sex        int8   `json:"sex" bson:"sex"`
	PhoneS     string `json:"phone" bson:"-"`
	Phone      int    `bson:"phone,omitempty"`
	SchoolName string `json:"schoolname" bson:"schoolname,omitempty"`
	EMail      string `json:"email" bson:"email,omitempty"`
	Name       string `json:"name" bson:"name,omitempty"`
	PicUrl     string `json:"pic" bson:"picurl,omitempty"`
	Dept       string `json:"dept" bson:"dept,omitempty"`
	Pic        []byte `bson:"pic,omitempty"`
}

func (d *Data) GetPic() {
	if d.PicUrl == "" {
		log.Printf("[%d]用户无头像", d.Puid)
		return
	}
	status, resp, err := fasthttp.Get(nil, d.PicUrl)
	if err != nil {
		log.Printf("[%d]头像请求错误:%s\n", d.Puid, err.Error())
		return
	}
	if status != fasthttp.StatusOK {
		log.Printf("[%d]头像请求失败 HTTP Status:%d\n", d.Puid, status)
		return
	}
	d.Pic = resp
	log.Printf("[%d]头像获取成功", d.Puid)
}

func getUserInfo(puid int) []byte {
	req := &fasthttp.Request{}
	req.SetRequestURI("https://contactsyd.chaoxing.com/pc/user/getUserInfo?puid=" + strconv.Itoa(puid))
	req.Header.SetContentType("application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.SetReferer("http://i.chaoxing.com/pc/contacts/home")
	req.Header.SetUserAgent("Mozilla/5.0 (X11; Linux x86_64; rv:73.0) Gecko/20100101 Firefox/73.0")
	req.Header.SetContentType("application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", config.Account.Cookie)
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.SetMethod("GET")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")

	resp := &fasthttp.Response{}

	client := &fasthttp.Client{}

	if err := client.Do(req, resp); err != nil {
		log.Printf("[%d]用户信息请求错误:%s\n", puid, err.Error())
		return nil
	}
	log.Printf("[%d]用户信息请求成功\n", puid)

	return resp.Body()
}

func parseRes(res []byte, puid int) (m Data, result int8) {
	var msg json.RawMessage
	env := SourceResult{
		Msg: &msg,
	}
	if err := json.Unmarshal(res, &env); err != nil {
		log.Panic(err)
	}
	switch env.Result {
	case 0:
		result = env.Result

	case 1:
		if err := json.Unmarshal(msg, &m); err != nil {
			log.Panic(err)
		}
		m.Puid = puid

		if m.PhoneS != "" {
			var err error
			m.Phone, err = strconv.Atoi(m.PhoneS)
			if err != nil {
				log.Panic(err)
			}
		}

		m.GetPic()
		result = env.Result
	}

	return
}
