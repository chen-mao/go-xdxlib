package xdxpci

// Interface allows us to get a list of all XDXCT PCI devices
type Interface interface {
	GetGPUs() ([]*XDXCTPCIDevice, error)
}

type xdxpci struct {
	// logger         logger
	// pciDevicesRoot string
	// pcidbPath      string
}

var _ Interface = (*xdxpci)(nil)

// XDXCTPCIDevice represents a PCI device for an XDXCT product
type XDXCTPCIDevice struct {
	Path       string
	Address    string
	Vendor     uint16
	Class      uint32
	ClassName  string
	Device     uint16
	DeviceName string
	Driver     string
	IommuGroup string
	NumaNode   int
	ISVF       bool
}

func New(opts ...Option) Interface {
	n := &xdxpci{}
	return n
}

// Option defines a function for passing options to the New() call
type Option func(*xdxpci)

// GetGPUs returns all XDXCT GPU devices on the system
func (p *xdxpci) GetGPUs() ([]*XDXCTPCIDevice, error) {
	return []*XDXCTPCIDevice{
		{
			Path:       "/sys/bus/pci/devices/0000:01:00.0",
			Address:    "0000:01:00.0",
			Vendor:     7917,
			Class:      196608,
			ClassName:  "Display controller",
			Device:     4912,
			DeviceName: "UNKNOWN_DEVICE",
			Driver:     "xgv",
			IommuGroup: "12",
			NumaNode:   -1,
			ISVF:       false,
		},
	}, nil
}
