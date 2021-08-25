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
	url := defaultConfig.Host + defaultConfig.DelUrl
	HandleExec(url)
}

func TestOpData(t *testing.T) {
	SkipRecordIds := make(map[int]int)
	for {
		host := defaultConfig.Host
		//host := defaultConfig.bakHost
		url := host + defaultConfig.Url
		//SkipRecordIds[3898] = true
		BatchExecOpTask(url, 1*time.Second, 0, SkipRecordIds)
		log.Println("未执行记录: ", SkipRecordIds)
		log.Println("end")
		log.Println("url: ", url)
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

func TestSyncOp(t *testing.T) {
	spaceIds := []string{
		"e03a937f9b1a46c89b19a989f8e05ea3",
	}
	var records sshmysql.Records
	records, err := sshmysql.GetDeployRecordListBySpaceIds(spaceIds)
	if err != nil {
		t.Log(err)
	} else {
		for {
			for _, v := range records {
				channel := strconv.Itoa(v.Channel)
				//构建token
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
				id := strconv.Itoa(v.Id)
				r, err := send(
					map[string]string{"deploy_record_id": id, "firm_destroy": "1"},
					map[string]string{"Authorization": "Bearer " + curToken},
					defaultConfig.Host+defaultConfig.DelUrl,
					"POST")
				//r, err := send(
				//	map[string]string{"appid": v.Appid, "version": v.Version, "runtime": v.Runtime},
				//	map[string]string{"Authorization": "Bearer " + curToken},
				//	"http://dp-a36vf7jujm9x7.gw002.oneitfarm.com/topology/topo/sync",
				//	"POST")
				if err != nil {
					t.Log(err)
				} else {
					t.Log(v.SpaceId, "SUCCESS", string(r))
				}
				time.Sleep(1 * time.Second)
			}
			time.Sleep(1 * time.Second)
		}

	}
}
