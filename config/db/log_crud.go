package db

func WriteLog(log Log) {
	dbi.Save(log)
}
