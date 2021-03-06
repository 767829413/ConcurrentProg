package request

import (
	"github.com/spf13/viper"
	"log"
	"testing"
	"time"
)

func TestHandleExec(t *testing.T) {
	HandleExec()
}

func TestOpData(t *testing.T) {
	var tn int
	for {
		BatchExecOpTask(2*time.Second, 3)
		tn += 1
		log.Println("执行次数: ", tn)
		time.Sleep(30 * time.Second)
	}
}

func TestGetToken(t *testing.T) {
	appkey := viper.Get("appkey").(string)
	channel := viper.Get("channel").(string)
	accountId := viper.Get("accountId").(string)
	issuer := viper.Get("issuer").(string)
	orgKey := viper.Get("orgKey").(string)
	subOrgKey := viper.Get("subOrgKey").(string)
	fromAppid := viper.Get("fromAppid").(string)
	appid := viper.Get("appid").(string)
	ucenterAlias := viper.Get("ucenterAlias").(string)
	curToken, _ := createToken(
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
	t.Log(curToken)
}
