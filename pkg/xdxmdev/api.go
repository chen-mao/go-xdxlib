package xdxmdev

// Interface allows us to get a list of XDXCT MDEV (vGPU) and parent devices
type Interface interface {
	GetAllMediatedDevices() ([]*Device, error)
	GetAllParentMediatedDevices() ([]*ParentDevice, error)
}
