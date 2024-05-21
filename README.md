# go-xdxlib
This repository contains a series of golang packages to manage xdx's GPU.

## Compile project
Make the code
```shell
go build ./...
```

## Test
1. Test for pci interface
```shell
go run examples/pci/pci-test.go
```
2. Prepare for installing the mdev driver before testing for mdev interface.
```shell
go run examples/mdev/mdev-test.go
```
