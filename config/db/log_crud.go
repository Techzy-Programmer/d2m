package db

func SaveDeploymentLog(log DeploymentLog) {
	dbi.Save(log)
}

func GetAllLogsForDeployment(id uint) []DeploymentLog {
	var logs []DeploymentLog
	dbi.Where("deploy_id = ?", id).Find(&logs)
	return logs
}
