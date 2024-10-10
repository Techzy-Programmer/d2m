package db

func WriteLog(log Log) {
	dbi.Save(log)
}

func GetAllLogsForDeploy(deployID uint) []Log {
	var logs []Log
	dbi.Where("deploy_id = ?", deployID).Find(&logs)
	return logs
}
