package db

var DeployLogStorage = make(map[uint][]Log)

func WriteLog(log Log) {
	dbi.Save(log)
}

func AppendDeployLog(deployId uint, log Log) {
	DeployLogStorage[deployId] = append(DeployLogStorage[deployId], log)
}

func FlushDeployLogs(deployId uint) {
	delete(DeployLogStorage, deployId)
}
