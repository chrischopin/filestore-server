package handler

import (
	dblayer "filestore-server/db"
	"filestore-server/util"
	"fmt"
	"net/http"
	"time"
)

const (
	salt = ")()(*()&*IO"
)

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.Redirect(w, r, "/static/view/signup.html", http.StatusFound)
		return
	}
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	if len(username) < 3 || len(password) < 5 {
		w.Write([]byte("Invalid Parameter"))
		return
	}
	enc_password := util.Sha1([]byte(password + salt))
	suc := dblayer.UserSignUp(username, enc_password)
	if suc {
		w.Write([]byte("SUCCESS"))
	} else {
		w.Write([]byte("FAILED"))
	}
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.Redirect(w, r, "/static/view/signin.html", http.StatusFound)
		return
	}

	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	encPasswd := util.Sha1([]byte(password + salt))

	// 1. 校验用户名及密码
	pwdChecked := dblayer.UserSignIn(username, encPasswd)
	if !pwdChecked {
		w.Write([]byte("FAILED"))
		return
	}

	// 2. 生成访问凭证(token)
	token := GenToken(username)
	upRes := dblayer.UpdateToken(username, token)
	if !upRes {
		w.Write([]byte("FAILED"))
		return
	}

	// 3. 登录成功后重定向到首页
	//w.Write([]byte("http://" + r.Host + "/static/view/home.html"))
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "http://" + r.Host + "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	}
	w.Write(resp.JSONBytes())
}

func GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Get Parameters
	r.ParseForm()
	username := r.Form.Get("username")
	//token := r.Form.Get("token")

	// 2. Check if token is valid
	//isValidToken := IsTokenValid(token)
	//if !isValidToken {
	//	w.WriteHeader(http.StatusForbidden)
	//	return
	//}

	// 3. Query to get user information
	user, err := dblayer.GetUserInfo(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// 4. Format return json
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	w.Write(resp.JSONBytes())
}

func IsTokenValid(token string) bool {
	if len(token) != 40 {
		return false
	}

	// TODO: 判断token的时效性，是否过期
	// example，假设token的有效期为1天   (根据同学们反馈完善, 相对于视频有所更新)
	tokenTS := token[:8]
	if util.Hex2Dec(tokenTS) < time.Now().Unix()-86400 {
		return false
	}

	// TODO: 从数据库表tbl_user_token查询username对应的token信息
	// TODO: 对比两个token是否一致
	// example, IsTokenValid方法增加传入参数username

	// if dblayer.GetUserToken(username) != token {
	// 	return false
	// }

	return true
}

func GenToken(username string) string {
	// 40位字符:md5(username+timestamp+token_salt)+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}
