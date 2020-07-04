package interfaces

type EurekaClient interface {
	Register()
	DeRegister()
}
