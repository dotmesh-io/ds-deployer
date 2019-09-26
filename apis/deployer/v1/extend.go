package deployer_v1

func (d *Deployment) ModelProxyEnabled() bool {
	return d.Metrics.GetImage() != ""
}
