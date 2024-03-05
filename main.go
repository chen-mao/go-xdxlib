package main

import (
	"fmt"

	"github.com/chen-mao/go-xdxlib/pkg/xdxpci"
)

// import (
// 	"fmt"

// 	"github.com/chen-mao/go-xdxlib/pkg/xdxmdev"
// )

// func main() {
// 	xdxmdev := xdxmdev.New()
// 	res, err := xdxmdev.GetAllDevices()
// 	if err != nil {
// 		return
// 	}
// 	for _, v := range res {
// 		fmt.Printf("Path --->: %v\n", v.Path)
// 		fmt.Printf("UUID --->: %v\n", v.UUID)
// 		fmt.Printf("MDEVType --->: %v\n", v.MDEVType)
// 		fmt.Printf("iommu_group-->: %v\n", v.IommuGroup)
// 	}
// }

func main() {
	xdxpci := xdxpci.New()
	res, err := xdxpci.GetGPUByIndex(0)
	if err != nil {
		return
	}
	fmt.Printf("Path --->: %v\n", res.Path)
	fmt.Printf("Address --->: %v\n", res.Address)
	fmt.Printf("Vendor --->: %v\n", res.Vendor)
	fmt.Printf("Class --->: %v\n", res.Class)
	fmt.Printf("DeviceId --->: %v\n", res.Device)
	fmt.Printf("Driver --->: %v\n", res.Driver)
	fmt.Printf("IommuGroup --->: %v\n", res.IommuGroup)
	fmt.Printf("numa_node --->: %v\n", res.NumaNode)
	fmt.Printf("Config Path --->: %v\n", res.Config.Path)
}
