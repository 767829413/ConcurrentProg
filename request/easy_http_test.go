package request

import (
	"log"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func TestHandleExec(t *testing.T) {
	host := defaultConfig.host
	//host := defaultConfig.bakHost
	//url := host + defaultConfig.url
	url := host + defaultConfig.delUrl
	HandleExec(url)
}

func TestOpData(t *testing.T) {
	for {
		host := defaultConfig.host
		//host := defaultConfig.bakHost
		url := host + defaultConfig.url
		taskUrl := host + defaultConfig.taskUrl
		BatchExecOpTask(url, taskUrl, 1*time.Second, 0)
		log.Println("end")
		time.Sleep(720 * time.Second)
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

func TestPrint(t *testing.T) {
	for i := 0; i < 10; i++ {
		log.Println(i)
	}
}
