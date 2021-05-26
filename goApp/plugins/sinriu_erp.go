package plugins

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
)

const Host = "https://nowkin.com:9949"
const LoginUrl = "/api/auth/login/"
const ProfileUrl = "/api/auth/u/my/"

var once sync.Once
var cli *erpCli
var ErpLoginRequiredErr = errors.New("请先登陆")
var ErpLogInParamsErr = errors.New("账号或密码错误")

type ErpProfile struct {
	username string
	realName string
	roleName string
}

type erpCli struct {
	token   string
	httpCli *http.Client
	ErpProfile
}

func (obj *erpCli) Login(username, password string) (err error) {
	defer func() {
		if r := recover(); nil != r {
			err = errors.New(fmt.Sprint(r))
		}
	}()

	resp, err := obj.httpCli.PostForm(Host+LoginUrl, url.Values{
		"username": []string{username},
		"password": []string{password},
	})
	if nil != err {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return ErpLogInParamsErr
	}
	body, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		return
	}

	res := make(map[string]interface{})
	err = json.Unmarshal(body, &res)
	if nil != err {
		return
	}
	obj.token = res["data"].(map[string]interface{})["access_token"].(string)
	obj.username = username
	return err
}

func (obj *erpCli) FlushProfile() (err error) {
	defer func() {
		if r := recover(); nil != r {
			err = errors.New(fmt.Sprint(r))
		}
	}()
	if obj.token == "" {
		return ErpLoginRequiredErr
	}
	req, err := http.NewRequest("GET", Host+ProfileUrl, nil)
	if nil != err {
		return
	}
	req.AddCookie(&http.Cookie{
		Name:  "access_token",
		Value: obj.token,
	})
	resp, err := obj.httpCli.Do(req)
	if nil != err {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return ErpLoginRequiredErr
	}
	body, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		return
	}
	res := make(map[string]interface{})
	err = json.Unmarshal(body, &res)
	if nil != err {
		return
	}
	data := res["data"].(map[string]interface{})
	obj.realName = data["real_name"].(string)
	obj.roleName = data["role_name"].(string)
	return err
}

func NewSinRiuErpCli() *erpCli {
	once.Do(func() {
		cli = &erpCli{httpCli: &http.Client{}}
	})
	return cli
}
