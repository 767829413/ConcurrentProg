package request

import (
	"ConcurrentProg/sshmysql"
	"log"
	"strconv"
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
		BatchExecOpTask(url, 1*time.Second, 0)
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

func TestSyncOp(t *testing.T) {
	spaceIds := []string{
		"bf4d63cbcdca441c8f0498a71c1effa9",
		"adaca618a99a41f0bff173f84fdc812f",
		"6f2a964514054f0d9f223950f02cfd14",
		"7d815fade7ce4037bd4a577a5e708bdd",
		"509f369d0e0d4cb69eeb38b264505fcf"}
	issuer := defaultConfig.issuer
	orgKey := defaultConfig.orgKey
	subOrgKey := defaultConfig.subOrgKey
	fromAppid := defaultConfig.fromAppid
	appid := defaultConfig.appid
	ucenterAlias := defaultConfig.ucenterAlias
	var records sshmysql.Records
	records, err := sshmysql.GetDeployRecordListBySpaceIds(spaceIds)
	if err != nil {
		t.Log(err)
	} else {
		for _, v := range records {
			channel := strconv.Itoa(v.Channel)
			//构建token
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
			r, err := send(
				map[string]string{"appid": v.Appid, "version": v.Version, "runtime": v.Runtime, "wer": "100"},
				map[string]string{"Authorization": "Bearer " + curToken},
				"http://dp-a36vf7jujm9x7.gw002.oneitfarm.com/topology/topo/sync",
				"POST")
			if err != nil {
				t.Log(err)
			} else {
				t.Log(v.SpaceId, "SUCCESS", string(r))
			}
			time.Sleep(1 * time.Second)
		}
	}
}
