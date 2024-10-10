package db

func SaveDeployment(req *Deployment) {
	dbi.Save(req)
}

func GetDeploymentByID(id uint) (Deployment, error) {
	var deployment Deployment
	res := dbi.First(&deployment, id)
	if res.Error != nil {
		return Deployment{}, res.Error
	}

	return deployment, nil
}

func GetAllDeployments() []Deployment {
	var deployments []Deployment
	res := dbi.Find(&deployments)
	if res.Error != nil {
		return nil
	}

	return deployments
}
