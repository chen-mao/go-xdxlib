package main

import (
	"fmt"

	"github.com/chenmao/go-xdxlib/pkg/xdxmdev"
)

func main() {
	xdxmdev := xdxmdev.New()
	res, err := xdxmdev.GetAllDevices()
	if err != nil {
		return
	}
	for _, v := range res {
		fmt.Printf("Path --->: %v\n", v.Path)
		fmt.Printf("UUID --->: %v\n", v.UUID)
		fmt.Printf("MDEVType --->: %v\n", v.MDEVType)
		fmt.Printf("iommu_group-->: %v\n", v.IommuGroup)
	}
}
