package request

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"log"
	"strconv"
	"time"
)

func init() {
	err := initConfig()
	if err != nil {
		panic(err)
	}
}

func HandleExec() {
	deployRecordId := viper.Get("deployRecordId").(string)
	appkey := viper.Get("appkey").(string)
	channel := viper.Get("channel").(string)
	accountId := viper.Get("accountId").(string)
	issuer := viper.Get("issuer").(string)
	orgKey := viper.Get("orgKey").(string)
	subOrgKey := viper.Get("subOrgKey").(string)
	fromAppid := viper.Get("fromAppid").(string)
	appid := viper.Get("appid").(string)
	ucenterAlias := viper.Get("ucenterAlias").(string)
	url := viper.Get("url").(string)
	method := viper.Get("method").(string)
	curToken, err := createToken(
		[]byte(""),
		issuer,
		appkey,
		channel,
		accountId,
		orgKey,
		subOrgKey,
		fromAppid,
		appid,
		ucenterAlias,
		"",
		[]map[string]string{
			{
				"appid":   appid,
				"appkey":  appkey,
				"channel": channel,
				"alias":   "default",
				"version": "0.0.0",
			},
		},
	)
	log.Println("Bearer " + curToken)
	if err != nil {
		log.Println(err)
	} else {
		req, resp := initHttp(
			url,
			method,
			map[string]string{"deploy_record_id": deployRecordId},
			map[string]string{"Authorization": "Bearer " + curToken})
		defer func() {
			// 用完需要释放资源
			fasthttp.ReleaseResponse(resp)
			fasthttp.ReleaseRequest(req)
		}()
		for {
			if err := fasthttp.Do(req, resp); err != nil {
				log.Println("请求失败:", err.Error())
				return
			}
			b := resp.Body()
			ooo := &ddd{}
			_ = json.Unmarshal(b, ooo)
			log.Println("result: ", string(b))
			log.Println("Pending")
			if ooo.Data.Code == 1 {
				break
			}
			time.Sleep(1 * time.Second)
		}
		log.Println("OK")
	}
}

func BatchExecOpTask(waitSecond time.Duration, retry int) {
	url := viper.Get("url").(string)
	taskUrl := viper.Get("task_url").(string)
	method := viper.Get("method").(string)
	issuer := viper.Get("issuer").(string)
	orgKey := viper.Get("orgKey").(string)
	subOrgKey := viper.Get("subOrgKey").(string)
	fromAppid := viper.Get("fromAppid").(string)
	appid := viper.Get("appid").(string)
	ucenterAlias := viper.Get("ucenterAlias").(string)
	records := getOpData(taskUrl, method)
	for _, v := range records {
		start := 0
		channel := strconv.Itoa(v.Channel)
		deployRecordId := strconv.Itoa(v.DeployRecordId)
		curToken, _ := createToken(
			[]byte(""),
			issuer,
			v.Appkey,
			channel,
			v.SpaceDeployId,
			orgKey,
			subOrgKey,
			fromAppid,
			appid,
			ucenterAlias,
			"",
			[]map[string]string{
				{
					"appid":   appid,
					"appkey":  v.Appkey,
					"channel": channel,
					"alias":   "default",
					"version": "0.0.0",
				},
			},
		)
		exec(deployRecordId, curToken, url, method, start, retry, waitSecond)
	}
}

func initHttp(url, method string, postArgs, header map[string]string) (*fasthttp.Request, *fasthttp.Response) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	req.Header.SetMethod(method)
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
	viper.AddConfigPath("B:/study/ConcurrentProg/request")
	return viper.ReadInConfig()
}

func createToken(
	SecretKey []byte,
	issuer,
	appkey,
	channel,
	accountId,
	orgKey,
	subOrgKey,
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
		SubOrgKey:    subOrgKey,
		FromAppid:    fromAppid,
		Appid:        appid,
		UcenterAlias: ucenterAlias,
		AclAlias:     aclAlias,
		CallStack:    CallStack,
	}
	claims.ExpiresAt = time.Now().Add(time.Hour * 72).Unix()
	claims.Issuer = issuer
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(SecretKey)
	return
}

func getOpData(url, method string) PendingList {
	var ooo PendingList
	req, resp := initHttp(
		url,
		method,
		map[string]string{},
		map[string]string{})
	defer func() {
		// 用完需要释放资源
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()
	if err := fasthttp.Do(req, resp); err != nil {
		fmt.Println("请求失败:", err.Error())
		return ooo
	}

	b := resp.Body()
	_ = json.Unmarshal(b, &ooo)
	return ooo
}

type PendingList []TaskRecord

type TaskRecord struct {
	Appkey         string `json:"appkey"`
	Channel        int    `json:"channel"`
	DeployRecordId int    `json:"deploy_record_id"`
	SpaceDeployId  string `json:"space_deploy_id"`
}

type jwtCustomClaims struct {
	jwt.StandardClaims
	// 追加自己需要的信息
	Appkey       string              `json:"appkey"`
	Channel      string              `json:"channel"`
	AccountId    string              `json:"account_id"`
	OrgKey       string              `json:"org_key"`
	SubOrgKey    string              `json:"sub_org_key"`
	FromAppid    string              `json:"from_appid"`
	Appid        string              `json:"appid"`
	UcenterAlias string              `json:"ucenter_alias"`
	AclAlias     string              `json:"acl_alias"`
	CallStack    []map[string]string `json:"call_stack"`
}

type ddd struct {
	State int
	Data  data
}

type data struct {
	Code int
}

func exec(deployRecordId, curToken, url, method string, start, retry int, waitSecond time.Duration) {
	req, resp := initHttp(
		url,
		method,
		map[string]string{"deploy_record_id": deployRecordId},
		map[string]string{"Authorization": "Bearer " + curToken})
	for {
		time.Sleep(waitSecond)
		start++
		if err := fasthttp.Do(req, resp); err != nil {
			log.Println("请求失败:", err.Error())
			return
		}
		b := resp.Body()
		ooo := &ddd{}
		_ = json.Unmarshal(b, ooo)
		log.Println("deploy record id: ", deployRecordId)
		log.Println("result: ", string(b))
		log.Println("Pending")
		if ooo.Data.Code == 1 || start > retry {
			// 用完需要释放资源
			fasthttp.ReleaseResponse(resp)
			fasthttp.ReleaseRequest(req)
			break
		}
	}
}
