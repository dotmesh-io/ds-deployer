package controller

const AnnControllerIdentifier = "deployer.dotscience.com/identifier"

func getDeployerID(annotations map[string]string) string {
	identifier, ok := annotations[AnnControllerIdentifier]
	if ok {
		return identifier
	}
	return ""
}
