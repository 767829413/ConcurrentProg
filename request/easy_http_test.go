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
	//host := defaultConfig.bakHost
	//url := host + defaultConfig.url
	url := defaultConfig.Host + defaultConfig.DelUrl
	HandleExec(defaultConfig.DeployRecordId, defaultConfig.Appkey, defaultConfig.Channel, defaultConfig.AccountId, url)
}

func TestOpData(t *testing.T) {
	for {
		host := defaultConfig.Host
		//host := defaultConfig.bakHost
		url := host + defaultConfig.Url
		BatchExecOpTask(url, 1*time.Second, 0)
		log.Println("end")
		time.Sleep(3600 * time.Second)
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
		"0ddeaa42c04d4a9184afa61861917d7b", "04070c8138e541a798b9b794175772b1", "55012dba5c104f3dae0f19bcd89effbd",
	}
	issuer := defaultConfig.Issuer
	orgKey := defaultConfig.OrgKey
	subOrgKey := defaultConfig.SubOrgKey
	fromAppid := defaultConfig.FromAppid
	appid := defaultConfig.Appid
	ucenterAlias := defaultConfig.UcenterAlias
	var records sshmysql.Records
	records, err := sshmysql.GetDeployRecordListBySpaceIds(spaceIds)
	if err != nil {
		t.Log(err)
	} else {
		for _, v := range records {
			channel := strconv.Itoa(v.Channel)
			id := strconv.Itoa(v.Id)
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
				map[string]string{"deploy_record_id": id},
				map[string]string{"Authorization": "Bearer " + curToken},
				defaultConfig.Host+defaultConfig.DelUrl,
				"POST")
			//r, err := send(
			//	map[string]string{"appid": v.Appid, "version": v.Version, "runtime": v.Runtime, "wer": "100"},
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
	}
}
