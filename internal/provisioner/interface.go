package provisioner

type Provisioner interface {
	Provision(...*GitFile) error
	Deprovision() error
}
