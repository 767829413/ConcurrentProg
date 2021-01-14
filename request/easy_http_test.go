package request

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"testing"
	"time"
)

type ddd struct {
	State int
	Data  data
}

type data struct {
	Code int
}

func TestHandleExec(t *testing.T) {
	err := initConfig()
	if err != nil { // Handle errors reading the config file
		t.Log(fmt.Errorf("Fatal error config file: %s \n", err))
	} else {
		deployRecordId := viper.Get("deployRecordId").(string)
		appkey := viper.Get("appkey").(string)
		channel := viper.Get("channel").(string)
		accountId := viper.Get("accountId").(string)
		issuer := viper.Get("issuer").(string)
		orgKey := viper.Get("orgKey").(string)
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
		t.Log("Bearer " + curToken)
		if err != nil {
			t.Log(err)
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
					t.Log("请求失败:", err.Error())
					return
				}
				b := resp.Body()
				ooo := &ddd{}
				json.Unmarshal(b, ooo)
				t.Log("result: ", string(b))
				t.Log("Pending")
				if ooo.Data.Code == 1 {
					break
				}
				time.Sleep(2 * time.Second)
			}
			t.Log("OK")
		}
	}
}

func TestGetToken(t *testing.T) {
	t.Log(createToken(
		[]byte(""),
		"xxxxx",
		"xxxxx",
		"xxxxx",
		"xxxxx",
		"xxxxx",
		"xxxxx",
		"xxxxx",
		"xxxxx",
		"",
		[]map[string]string{
			{
				"appid":   "xxxxx",
				"appkey":  "xxxxx",
				"channel": "xxxxx",
				"alias":   "xxxxx",
				"version": "xxxxx",
			},
		},
	))
}
