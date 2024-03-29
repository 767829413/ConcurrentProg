package sshmysql

import "strings"

type DeployRecord struct {
	Id            int    `json:"id"`
	Appkey        string `json:"appkey"`
	Channel       int    `json:"channel"`
	Appid         string `json:"appid"`
	Version       string `json:"version"`
	Runtime       string `json:"runtime"`
	SpaceId       string `json:"space_id"`
	SpaceDeployId string `json:"space_deploy_id"`
	SubOrgKey     string `json:"sub_org_key"`
}

type Records []*DeployRecord

func GetDeployRecordListBySpaceIds(spaceIds []string) (records Records, err error) {
	str := "'" + strings.Join(spaceIds, "','") + "'"
	rows, err := Db.Query("SELECT `id`,`appkey`,`channel`,`appid`,`version`,`runtime`,`space_id`,`space_deploy_id` FROM `deploy_record` WHERE `space_id` IN (" + str + ")")
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		data := &DeployRecord{}
		_ = rows.Scan(&data.Id, &data.Appkey, &data.Channel, &data.Appid, &data.Version, &data.Runtime, &data.SpaceId, &data.SpaceDeployId)
		records = append(records, data)
	}
	return
}
