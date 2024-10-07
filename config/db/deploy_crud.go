package db

func SaveDeployment(req *Deployment) {
	dbi.Save(req)
}

func GetAllDeployments() []Deployment {
	var deployments []Deployment
	dbi.Find(&deployments)
	return deployments
}
