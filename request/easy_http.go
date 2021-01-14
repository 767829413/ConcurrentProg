package request

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"time"
)

func initHttp(url, method string, postArgs, header map[string]string) (*fasthttp.Request, *fasthttp.Response) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	req.Header.SetMethod("POST")
	req.SetRequestURI(url)
	for k, v := range header {
		req.Header.Add(k, v)
	}
	for k, v := range postArgs {
		req.PostArgs().Add(k, v)
	}
	return req, resp
}

func initConfig() error {
	viper.SetConfigName("test")
	viper.SetConfigType("json")
	viper.AddConfigPath("C:/Users/rtu/Documents/code/study/ConcurrentProg/request/")
	return viper.ReadInConfig()
}

func createToken(
	SecretKey []byte,
	issuer,
	appkey,
	channel,
	accountId,
	orgKey,
	fromAppid,
	appid,
	ucenterAlias,
	aclAlias string,
	CallStack []map[string]string) (tokenString string, err error) {
	claims := &jwtCustomClaims{
		Appkey:       appkey,
		Channel:      channel,
		AccountId:    accountId,
		OrgKey:       orgKey,
		FromAppid:    fromAppid,
		Appid:        appid,
		UcenterAlias: ucenterAlias,
		AclAlias:     aclAlias,
		CallStack:    CallStack,
	}
	claims.ExpiresAt = int64(time.Now().Add(time.Hour * 72).Unix())
	claims.Issuer = issuer
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(SecretKey)
	return
}

type jwtCustomClaims struct {
	jwt.StandardClaims
	// 追加自己需要的信息
	Appkey       string              `json:"appkey"`
	Channel      string              `json:"channel"`
	AccountId    string              `json:"account_id"`
	OrgKey       string              `json:"org_key"`
	FromAppid    string              `json:"from_appid"`
	Appid        string              `json:"appid"`
	UcenterAlias string              `json:"ucenter_alias"`
	AclAlias     string              `json:"acl_alias"`
	CallStack    []map[string]string `json:"call_stack"`
}
