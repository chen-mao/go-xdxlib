package xdxmdev

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/chen-mao/go-xdxlib/pkg/xdxpci"
)

type ParentDevice struct {
	*xdxpci.XDXCTPCIDevice
	mdevPaths map[string]string
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
