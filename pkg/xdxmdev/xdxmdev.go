package xdxmdev

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/chen-mao/go-xdxlib/pkg/xdxpci"
)

const (
	mdevParentRoot = "/sys/class/mdev_bus"
	mdevDeviceRoot = "/sys/bus/mdev/devices"
)

type ParentDevice struct {
	*xdxpci.XDXCTPCIDevice
	mdevPaths map[string]string
}

type Device struct {
	Path       string
	UUID       string
	MDEVType   string
	Driver     string
	IommuGroup int
	Parent     ParentDevice
}

type Interface interface {
	GetAllDevices() ([]*Device, error)
}

type xdxmdev struct {
	mdevParentRoot string
	mdevDeviceRoot string
}

var _ Interface = (*xdxmdev)(nil)

func New() Interface {
	return &xdxmdev{
		mdevParentRoot: mdevParentRoot,
		mdevDeviceRoot: mdevDeviceRoot,
	}
}

// GetAllDevices returns all XDXCT mdev (vGPU) devices on the system
func (m *xdxmdev) GetAllDevices() ([]*Device, error) {
	deviceDirs, err := os.ReadDir(m.mdevDeviceRoot)
	if err != nil {
		return nil, fmt.Errorf("unable to read PCI bus devices: %v", err)
	}
	var xdxdevices []*Device
	for _, deviceDir := range deviceDirs {
		xdxdevice, err := NewDevice(m.mdevDeviceRoot, deviceDir.Name())
		if err != nil {
			return nil, fmt.Errorf("error constructing xdxct MDEV device: %v", err)
		}
		if xdxdevice == nil {
			continue
		}
		xdxdevices = append(xdxdevices, xdxdevice)
	}
	return xdxdevices, nil
}

func NewDevice(root string, uuid string) (*Device, error) {
	path := path.Join(root, uuid)

	m, err := newMdev(path)
	if err != nil {
		return nil, err
	}

	parent, err := NewParentDevice(m.parentDevicePath())
	if err != nil {
		return nil, fmt.Errorf("error getting mdev type: %v", err)
	}
	if parent == nil {
		return nil, nil
	}

	mdevType, err := m.Type()
	if err != nil {
		return nil, fmt.Errorf("error get mdev type: %v", err)
	}

	driver, err := m.Driver()
	if err != nil {
		return nil, fmt.Errorf("error detecting driver: %v", err)
	}

	iommu_group, err := m.IommuGroup()
	if err != nil {
		return nil, fmt.Errorf("error detecting Iommu Group: %v", err)
	}

	device := Device{
		Path:       path,
		UUID:       uuid,
		MDEVType:   mdevType,
		Driver:     driver,
		IommuGroup: int(iommu_group),
		Parent:     *parent,
	}
	return &device, nil
}

type mdev string

func newMdev(devicePath string) (mdev, error) {
	mdevDir, err := filepath.EvalSymlinks(devicePath)
	if err != nil {
		return "", fmt.Errorf("error resolving symlink for %s: %v", devicePath, err)
	}
	return mdev(mdevDir), nil
}

// parentDevicePath() return "/sys/devices/pci0000:00/0000:00:01.1/<pcu-id>"
func (m *mdev) parentDevicePath() string {
	return path.Dir(string(*m))
}

func (m *mdev) resolve(target string) (string, error) {
	resolved, err := filepath.EvalSymlinks(path.Join(string(*m), target))
	if err != nil {
		return "", fmt.Errorf("error resolving %q: %v", target, err)
	}
	return resolved, nil
}

func (m *mdev) Type() (string, error) {
	reg := regexp.MustCompile(`Type Name: (\w+)`)

	mdevTypeDir, err := m.resolve(string("mdev_type"))
	if err != nil {
		return "", err
	}
	mdevTypeName, err := os.ReadFile(path.Join(mdevTypeDir, "name"))
	if err != nil {
		return "", fmt.Errorf("unable to read mdev_type name for mdev %s: %v", mdevTypeName, err)
	}

	mdevTypeStr := strings.TrimSpace(string(mdevTypeName))

	matches := reg.FindStringSubmatch(mdevTypeStr)
	if len(matches) > 1 {
		extracted := matches[1]
		return extracted, nil
	} else {
		return "", fmt.Errorf("unable to parse mdev_type name for mdev %s", mdevTypeName)
	}
}

func (m *mdev) Driver() (string, error) {
	driver, err := m.resolve(string("driver"))
	if err != nil {
		return "", err
	}

	return filepath.Base(driver), nil
}

func (m *mdev) IommuGroup() (int64, error) {
	IommuGroup, err := m.resolve(string("iommu_group"))
	if err != nil {
		return -1, err
	}

	IommuGroupStr := strings.TrimSpace(filepath.Base(IommuGroup))
	IommuGroupInt, err := strconv.ParseInt(IommuGroupStr, 0, 64)
	if err != nil {
		return -1, fmt.Errorf("unable to convert iommu_group string to int64: %v", err)
	}

	return IommuGroupInt, nil
}

func NewParentDevice(devicePath string) (*ParentDevice, error) {
	reg := regexp.MustCompile(`Type Name: (\w+)`)
	xdxDevice, err := newXDXCTPCIDeviceFromPath(devicePath)
	if err != nil {
		return nil, fmt.Errorf("failed to construct XDXCT PCI device: %v", err)
	}
	if xdxDevice == nil {
		// Not a XDXCT device
		return nil, err
	}

	paths, err := filepath.Glob(fmt.Sprintf("%s/mdev_supported_types/xgv-XGV_V0_*/name", xdxDevice.Path))
	if err != nil {
		return nil, fmt.Errorf("unable to get files in mdev_supported_types directory: %v", err)
	}
	mdevTypesMap := make(map[string]string)
	for _, path := range paths {
		name, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("unable to read file %s: %v", path, err)
		}
		mdevTypeStr := strings.TrimSpace(string(name))
		matches := reg.FindStringSubmatch(mdevTypeStr)
		if len(matches) > 1 {
			extracted := matches[1]
			mdevTypesMap[extracted] = filepath.Dir(path)
		} else {
			return nil, fmt.Errorf("unable to parse mdev_type name for mdev %s", mdevTypeStr)
		}
	}
	fmt.Printf("---> %v", mdevTypesMap)
	return &ParentDevice{
		xdxDevice,
		mdevTypesMap,
	}, nil
}

func newXDXCTPCIDeviceFromPath(devicePath string) (*xdxpci.XDXCTPCIDevice, error) {
	root := filepath.Dir(devicePath)
	address := filepath.Base(devicePath)
	return xdxpci.New(xdxpci.WithPCIDevicesRoot(root)).
		GetGPUByPciBusID(address)
}
