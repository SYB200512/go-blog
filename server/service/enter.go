package service

// ServiceGroup 服务组
type ServiceGroup struct {
	EsService
}

// ServiceGroupApp 服务组应用实例
var ServiceGroupApp = new(ServiceGroup)
