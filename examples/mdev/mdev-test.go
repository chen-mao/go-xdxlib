package main

import (
	"fmt"

	"github.com/chen-mao/go-xdxlib/pkg/xdxmdev"
)

func main() {
	xdxmdev := xdxmdev.New()
	fmt.Println("------- Mediated Devices --------")
	resMdev, err := xdxmdev.GetAllMediatedDevices()
	if err != nil {
		return
	}
	for _, v := range resMdev {
		fmt.Printf("Path --->: %v\n", v.Path)
		fmt.Printf("UUID --->: %v\n", v.UUID)
		fmt.Printf("MDEVType --->: %v\n", v.MDEVType)
		fmt.Printf("iommu_group-->: %v\n", v.IommuGroup)
	}

	fmt.Println("------- Parent Mediated Devices --------")
	resParentMD, err := xdxmdev.GetAllParentMediatedDevices()
	if err != nil {
		return
	}
	for _, v := range resParentMD {
		fmt.Printf("Address --->: %v\n", v.Address)
	}

}
