package request

import (
	"log"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func TestHandleExec(t *testing.T) {
	url := viper.Get("del_url").(string)
	HandleExec(url)
}

func TestOpData(t *testing.T) {
	BatchExecOpTask(1*time.Second, 2)
	log.Println("end")
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

func TestPrint(t *testing.T) {
	for i := 0; i < 10; i++ {
		log.Println(i)
	}
}
