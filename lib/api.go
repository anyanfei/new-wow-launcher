package lib

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/tidwall/gjson"
)

const (
	API = "https://www.laghaim.cn/api/home/"
)

func GetCode() (string, string, error) {
	api := API + "getVerifyCode"
	resp, _ := http.Get(api)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	//解析json
	if gjson.GetBytes(body, "code").Int() == 0 {
		return gjson.GetBytes(body, "data.base64_info").String(), gjson.GetBytes(body, "data.captcha_id").String(), nil
	}
	return "", "", errors.New("无法获取注册验证码")
}
func Register(account, password, confirmPassword, email, captchaId, answer string) error {
	api := API + "填写你的注册地址"
	params := make(map[string]interface{})
	params["username"] = account
	params["user_pass"] = password
	params["user_pass_repeat"] = confirmPassword
	params["user_email"] = email
	params["captcha_id"] = captchaId
	params["answer"] = answer
	body, _ := json.Marshal(params)
	req, err := http.NewRequest("POST", api, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if gjson.GetBytes(responseBody, "code").Int() == 0 {
		return nil
	} else {
		return errors.New(gjson.GetBytes(responseBody, "message").String())
	}

}

func GetNotice() (string, string) {
	api := API + "getServerNews?page=1&page_size=1"
	resp, _ := http.Get(api)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", ""
	}
	//解析json
	if gjson.GetBytes(body, "code").Int() == 0 {
		list := gjson.GetBytes(body, "data.lists").Array()
		return list[0].Get("title").String(), list[0].Get("content").String()
	}
	return "", ""
}
