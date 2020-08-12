package handler

import (
	dblayer "../db"
	"../util"
	"fmt"
	"net/http"
	"time"
)

const (
	pwdSalt = "#890"
)

// SignupHandler : 处理用户注册请求
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.Redirect(w, r, "/static/view/signup.html", http.StatusFound)
		return
	}
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	fmt.Println(username)
	//传入参数验证
	if len(username) < 3 || len(password) < 5 {
		w.Write([]byte("Invalid parameter"))
		return
	}

	// 对密码进行加盐及取Sha1值加密
	encPassword := util.Sha1([]byte(password + pwdSalt))
	// 将用户信息注册到用户表中
	suc := dblayer.UserSignup(username, encPassword)
	if suc {
		w.Write([]byte("SUCCESS"))
	} else {
		w.Write([]byte("FAILED"))
	}
}

//登录接口
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	userName := r.Form.Get("username")
	password := r.Form.Get("password")

	encPasswd := util.Sha1([]byte(password + pwdSalt))
	
	//1. 校验用户名及密码
	pwdChecked := dblayer.UserSignIn(userName, encPasswd)
	if !pwdChecked {
		w.Write([]byte("FAILED"))
		return
	}

	//2. 生成访问凭证（token）
	token := GenToken(userName)
	UpRes := dblayer.UpdateToken(userName,token)
	if !UpRes{
		w.Write([]byte("FAILED"))
	}


	//3. 登录成功后重定向到首页
	w.Write([]byte("http://" + r.Host + "/static/view/home.html"))
	//resp := util.RespMsg{
	//	Code: 0,
	//	Msg:  "OK",
	//	Data: struct {
	//		Location string
	//		Username string
	//		Token    string
	//	}{
	//		Location: "http://" + r.Host + "/static/view/home.html",
	//		Username: username,
	//		Token:    token,
	//	},
	//}
	//w.Write(resp.JSONBytes())
}

//生成Token
func GenToken(username string) string {
	//当前时间戳
	ts := fmt.Sprintf("%x", time.Now().Unix())
	//md5加密拼接
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}

