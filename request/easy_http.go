package request

import (
	"ConcurrentProg/sshmysql"
	"ConcurrentProg/util"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
)

type config struct {
	DeployRecordId string `json:"deploy_record_id"`
	Appkey         string `json:"appkey"`
	Channel        string `json:"channel"`
	AccountId      string `json:"account_id"`
	Issuer         string `json:"issuer"`
	OrgKey         string `json:"org_key"`
	SubOrgKey      string `json:"sub_org_key"`
	FromAppid      string `json:"from_appid"`
	Appid          string `json:"appid"`
	UcenterAlias   string `json:"ucenter_alias"`
	Host           string `json:"host"`
	BakHost        string `json:"bak_host"`
	Url            string `json:"url"`
	TaskUrl        string `json:"task_url"`
	DelUrl         string `json:"del_url"`
	Method         string `json:"method"`
}

var defaultConfig config

func init() {
	err := util.InitConfig("test", "json", "B:/study/ConcurrentProg/request")
	defaultConfig = config{
		DeployRecordId: viper.Get("deployRecordId").(string),
		Appkey:         viper.Get("appkey").(string),
		Channel:        viper.Get("channel").(string),
		AccountId:      viper.Get("accountId").(string),
		Issuer:         viper.Get("issuer").(string),
		OrgKey:         viper.Get("orgKey").(string),
		SubOrgKey:      viper.Get("subOrgKey").(string),
		FromAppid:      viper.Get("fromAppid").(string),
		Appid:          viper.Get("appid").(string),
		UcenterAlias:   viper.Get("ucenterAlias").(string),
		Host:           viper.Get("host").(string),
		BakHost:        viper.Get("bak_host").(string),
		Url:            viper.Get("url").(string),
		TaskUrl:        viper.Get("task_url").(string),
		DelUrl:         viper.Get("del_url").(string),
		Method:         viper.Get("method").(string),
	}
	if err != nil {
		panic(err)
	}
}

func HandleExec(url string) {
	curToken, err := createToken(
		[]byte(""),
		defaultConfig.Issuer,
		defaultConfig.Appkey,
		defaultConfig.Channel,
		defaultConfig.AccountId,
		defaultConfig.OrgKey,
		defaultConfig.SubOrgKey,
		defaultConfig.FromAppid,
		defaultConfig.Appid,
		defaultConfig.UcenterAlias,
		"",
		[]map[string]string{
			{
				"appid":   defaultConfig.Appid,
				"appkey":  defaultConfig.Appkey,
				"channel": defaultConfig.Channel,
				"alias":   "default",
				"version": "0.0.0",
			},
		},
	)
	log.Println("Bearer " + curToken)
	if err != nil {
		log.Println(err)
	} else {
		for {
			b, err := send(
				map[string]string{"deploy_record_id": defaultConfig.DeployRecordId},
				map[string]string{"Authorization": "Bearer " + curToken},
				url,
				defaultConfig.Method)
			if err != nil {
				log.Println("请求失败:", err.Error())
				return
			}
			ooo := &ddd{}
			_ = json.Unmarshal(b, ooo)
			log.Println("result: ", string(b))
			log.Println("Pending")
			if ooo.State != 1 || ooo.Data.Code == 1 {
				break
			}
			time.Sleep(1 * time.Second)
		}
		log.Println("end")
	}
}

func BatchExecOpTask(url string, waitSecond time.Duration, retry int) {
	records, err := sshmysql.GetDeployOpRecordList()
	if err != nil {
		fmt.Println(err)
	} else {
		for _, v := range records {
			start := 0
			channel := strconv.Itoa(v.Channel)
			deployRecordId := strconv.Itoa(v.DeployRecordId)
			curToken, _ := createToken(
				[]byte(""),
				defaultConfig.Issuer,
				v.Appkey,
				channel,
				v.SpaceDeployId,
				defaultConfig.OrgKey,
				defaultConfig.SubOrgKey,
				defaultConfig.FromAppid,
				defaultConfig.Appid,
				defaultConfig.UcenterAlias,
				"",
				[]map[string]string{
					{
						"appid":   defaultConfig.Appid,
						"appkey":  v.Appkey,
						"channel": channel,
						"alias":   "default",
						"version": "0.0.0",
					},
				},
			)
			exec(deployRecordId, curToken, url, defaultConfig.Method, start, retry, waitSecond)
		}
	}

}

func initHttpRequest(url, method string, postArgs, header map[string]string) (*fasthttp.Request, *fasthttp.Response) {
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

func getOpData(url, method string) sshmysql.PendingList {
	var ooo sshmysql.PendingList
	b, err := send(map[string]string{}, map[string]string{}, url, method)
	if err != nil {
		fmt.Println("请求失败:", err.Error())
		return ooo
	}
	_ = json.Unmarshal(b, &ooo)
	return ooo
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
	for {
		start++
		b, err := send(map[string]string{"deploy_record_id": deployRecordId},
			map[string]string{"Authorization": "Bearer " + curToken},
			url, method)
		if err != nil {
			log.Println("请求失败:", err.Error())
			return
		}
		ooo := &ddd{}
		_ = json.Unmarshal(b, ooo)
		log.Println("deploy record id: ", deployRecordId)
		log.Println("result: ", string(b))
		if ooo.Data.Code == 1 || start > retry {
			break
		}
		time.Sleep(waitSecond)
	}
}

func send(data, header map[string]string, url, method string) (r []byte, err error) {
	req, resp := initHttpRequest(url, method, data, header)
	defer func() {
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()
	err = fasthttp.Do(req, resp)
	r = resp.Body()
	return r, err
}
