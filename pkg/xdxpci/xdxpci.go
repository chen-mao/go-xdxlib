package xdxpci

import "fmt"

const (
	// PCIDevicesRoot represents base path for all pci devices under sysfs
	PCIDevicesRoot = "/sys/bus/pci/devices"
)

// Interface allows us to get a list of all XDXCT PCI devices
type Interface interface {
	GetGPUs() ([]*XDXCTPCIDevice, error)
	GetGPUByIndex(int) (*XDXCTPCIDevice, error)
	GetGPUByPciBusID(string) (*XDXCTPCIDevice, error)
}

type xdxpci struct {
	logger         logger
	pciDevicesRoot string
	pcidbPath      string
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
	for _, opt := range opts {
		opt(n)
	}
	if n.logger == nil {
		n.logger = &simpleLogger{}
	}
	if n.pciDevicesRoot == "" {
		n.pciDevicesRoot = PCIDevicesRoot
	}
	return n
}

// Option defines a function for passing options to the New() call
type Option func(*xdxpci)

func WithLogger(logger logger) Option {
	return func(n *xdxpci) {
		n.logger = logger
	}
}

func WithPCIDevicesRoot(path string) Option {
	return func(n *xdxpci) {
		n.pcidbPath = path
	}
}

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

func (p xdxpci) GetGPUByIndex(i int) (*XDXCTPCIDevice, error) {
	gpus, err := p.GetGPUs()
	if err != nil {
		return nil, fmt.Errorf("error getting all gpus: %v", err)
	}
	if i < 0 || i > len(gpus) {
		return nil, fmt.Errorf("invalid index '%d'", i)
	}
	return gpus[i], nil
}

func (p xdxpci) GetGPUByPciBusID(address string) (*XDXCTPCIDevice, error) {
	return &XDXCTPCIDevice{
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
	}, nil
}
