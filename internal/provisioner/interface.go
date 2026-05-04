package provisioner

type Provisioner interface {
	Provision() error
	Deprovision() error
}
