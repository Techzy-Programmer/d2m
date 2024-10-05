package db

func SaveDeployment(req *Deployment) {
	dbi.Save(req)
}
