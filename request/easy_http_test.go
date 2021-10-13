package request

import (
	"ConcurrentProg/sshmysql"
	"fmt"
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
		BatchExecOpTask(url, 1*time.Second, 0, SkipRecordIds)
		log.Println("未执行记录: ", SkipRecordIds)
		log.Println("end")
		log.Println("url: ", url)
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

func TestSyncOp(t *testing.T) {
	spaceIds := []string{
		"1463a85457ec4fa39df0673825d89e70",

	}
	var records sshmysql.Records
	records = getRecordDataBySpaceIds(spaceIds, defaultConfig.Host+"/deploy/record/exec/list", "POST")
	for {
		for _, v := range records {
			channel := strconv.Itoa(v.Channel)
			var subOrgKey string
			if v.SubOrgKey != "" {
				subOrgKey = v.SubOrgKey
			}else{
				subOrgKey = defaultConfig.SubOrgKey
			}
			//构建token
			curToken, _ := createToken(
				[]byte(""),
				defaultConfig.Issuer,
				v.Appkey,
				channel,
				v.SpaceDeployId,
				defaultConfig.OrgKey,
				subOrgKey,
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
			//销毁
			reqUrl := defaultConfig.Host + defaultConfig.DelUrl
			//清理
			//reqUrl := defaultConfig.Host + defaultConfig.ClearUrl
			//重试
			//reqUrl := defaultConfig.Host + defaultConfig.RetryUrl
			//回滚
			//reqUrl := defaultConfig.Host + defaultConfig.RollbackUrl

			id := strconv.Itoa(v.Id)
			r, err := send(
				map[string]string{"deploy_record_id": id, "firm_destroy": "1"},
				map[string]string{"Authorization": "Bearer " + curToken},
				reqUrl,
				"POST")

			//r, err := send(
			//	map[string]string{"appid": v.Appid, "version": v.Version, "runtime": v.Runtime},
			//	map[string]string{"Authorization": "Bearer " + curToken},
			//	"https://dp-a36vf7jujm9x7.gw002.oneitfarm.com/topology/topo/sync",
			//	"POST")

			if err != nil {
				t.Log(err)
			} else {
				t.Log(v.SpaceId, "SUCCESS", string(r))
			}
			time.Sleep(1 * time.Second)
		}
		fmt.Println("结束一轮")
		time.Sleep(1 * time.Second)
	}
}
