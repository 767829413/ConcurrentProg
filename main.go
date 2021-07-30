package main

import (
	"ConcurrentProg/sshmysql"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

func main() {
	setChannelTmpOp()
}

func setChannelTmpOp() {
	sshHost := "localhost"
	sshPort := "22222"
	sshUser := "root"
	sshPubKeyPath := "B:/study/ConcurrentProg/sshmysql/test_sql_link.ssh-rsa"
	dbUser1 := "dbdktqs8tapadhc7"
	dbPass1 := "ndn8bnxf1pmaeq5k8z186wi1ssl46wlx"
	dbHost1 := "rm-uf6gz6k88ao353q8m.mysql.rds.aliyuncs.com:3306"
	dbName1 := "dbdktqs8tapadhc7"

	dbUser2 := "dbhxlpe8xiug7f2d"
	dbPass2 := "lei2pqlkjylhrdyqvjdiige6xulm9zmi"
	dbHost2 := "rm-uf6gz6k88ao353q8m.mysql.rds.aliyuncs.com:3306"
	dbName2 := "dbhxlpe8xiug7f2d"

	Db1, _ := sshmysql.GetDbHandler(
		sshHost,
		sshPort,
		sshUser,
		sshPubKeyPath,
		dbUser1,
		dbPass1,
		dbHost1,
		dbName1)
	Db2, _ := sshmysql.GetDbHandler(
		sshHost,
		sshPort,
		sshUser,
		sshPubKeyPath,
		dbUser2,
		dbPass2,
		dbHost2,
		dbName2)
	rows, err := Db1.Query("SELECT appid,version,channel_data,`group`,user_appkey,user_channel FROM `channels` WHERE `space_id` = 'e980f26223804eeb9e0e16fa2654f439' AND `appid` = 'zlapficvjthds407hdbxlmgownj2qvxe' AND `state` = '2'")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var nodeStatus = make(map[string]*NodeStatu)
	for rows.Next() {
		data := &Channel{}
		_ = rows.Scan(&data.Appid, &data.Version, &data.ChannelData, &data.Group, &data.UserAppkey, &data.UserChannel)
		str, _ := JsonEncode(data)
		strChannel := strconv.Itoa(data.UserChannel)
		nodeStatu := &NodeStatu{
			Appkey:           "topogli44vvvgbvkohoy2kutdyvpli5s",
			Channel:          321,
			DeployRecordId:   2970,
			Appid:            data.Appid,
			Version:          data.Version,
			ChannelData:      data.ChannelData,
			DeployAppkey:     data.UserAppkey,
			DeployChannel:    data.UserChannel,
			RawNodeReturn:    str,
			Group:            data.Group,
			State:            0,
			IsDelete:         0,
			ChildChannelFlag: Md5(data.UserAppkey + strChannel),
		}
		nodeStatus[nodeStatu.ChildChannelFlag] = nodeStatu
	}
	num := 0
	sql := "INSERT INTO nodes_status (appkey, channel, deploy_record_id, appid, version, channel_data,deploy_appkey, deploy_channel, raw_node_return, `group`, state, is_delete, child_channel_flag) VALUES "
	total := len(nodeStatus) - 1
	for _, status := range nodeStatus {
		strSQL := "('%s', %d, %d, '%s', '%s', '%s', '%s', %d, '%s', '%s', %d, %d, '%s')"
		if num < total {
			strSQL += ","
		}
		sql += fmt.Sprintf(strSQL,
			status.Appkey,
			status.Channel,
			status.DeployRecordId,
			status.Appid,
			status.Version,
			status.ChannelData,
			status.DeployAppkey,
			status.DeployChannel,
			status.RawNodeReturn,
			status.Group,
			status.State,
			status.IsDelete,
			status.ChildChannelFlag,
		)
		num++
	}
	sql += ";"
	fmt.Println(sql)
	fmt.Println(Db2.Query(sql))
	//fmt.Println(len(nodeStatus))
}

type Channel struct {
	Appid       string `json:"appid"`
	Version     string `json:"version"`
	ChannelData string `json:"channel_data"`
	Group       string `json:"group"`
	UserAppkey  string `json:"user_appkey"`
	UserChannel int    `json:"user_channel"`
}

type NodeStatu struct {
	Appkey           string `json:"appkey"`
	Channel          int    `json:"channel"`
	DeployRecordId   int    `json:"deploy_record_id"`
	Appid            string `json:"appid"`
	Version          string `json:"version"`
	ChannelData      string `json:"channel_data"`
	DeployAppkey     string `json:"deploy_appkey"`
	DeployChannel    int    `json:"deploy_channel"`
	RawNodeReturn    string `json:"raw_node_return"`
	Group            string `json:"group"`
	State            int    `json:"state"`
	IsDelete         int    `json:"is_delete"`
	ChildChannelFlag string `json:"child_channel_flag"`
}

func Md5(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func JsonEncode(data interface{}) (string, error) {
	jsons, err := json.Marshal(data)
	return string(jsons), err
}
