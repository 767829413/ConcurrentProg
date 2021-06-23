package sshmysql

type PendingList []*TaskRecord

type TaskRecord struct {
	Appkey         string `json:"appkey"`
	Channel        int    `json:"channel"`
	DeployRecordId int    `json:"deploy_record_id"`
	SpaceDeployId  string `json:"space_deploy_id"`
}

func GetDeployOpRecordList() (records PendingList, err error) {
	rows, err := Db.Query("SELECT `appkey`,`channel`,`deploy_record_id`,`space_deploy_id` FROM `deploy_op_record` WHERE `is_delete` = 0 AND `state` = 0")
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		data := &TaskRecord{}
		_ = rows.Scan(&data.Appkey, &data.Channel, &data.DeployRecordId, &data.SpaceDeployId)
		records = append(records, data)
	}
	return
}
