/*
Copyright (c) 2024, Shanghai Iluvatar CoreX Semiconductor Co., Ltd.
All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License"); you may
not use this file except in compliance with the License. You may obtain
a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"log"

	"gitee.com/deep-spark/go-ixml/pkg/ixml"
)

func main() {
	ret := ixml.Init()
	if ret != ixml.SUCCESS {
		log.Fatalf("Unable to initialize IXML, ret: %v", ret)
	}
	defer func() {
		ret := ixml.Shutdown()
		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to shutdown IXML, ret: %v", ret)
		}
	}()

	count, ret := ixml.DeviceGetCount()
	if ret != ixml.SUCCESS {
		log.Fatalf("Unable to get device count, ret: %v", ret)
	}
	fmt.Printf("Device count: %d\n", count)
	if count <= 0 {
		log.Fatalf("No GPU found")
	}

	var device ixml.Device
	if ret := ixml.DeviceGetHandleByIndex(0, &device); ret != ixml.SUCCESS {
		log.Fatalf("Unable to get device handle by index, ret: %v", ret)
	}
	uuid, ret := ixml.DeviceGetUUID(device)
	if ret != ixml.SUCCESS {
		log.Fatalf("Unable to get device uuid, ret: %v", ret)
	}
	fmt.Printf("Device UUID: %s\n", uuid)

	fmt.Printf("Start to get attributes of device: %s\n", uuid)
	device, ret = ixml.GetHandleByUUID(uuid)
	if ret != ixml.SUCCESS {
		log.Fatalf("Unable to get Handle by uuid, ret: %v", ret)
	}

	name, ret := device.GetName()
	if ret != ixml.SUCCESS {
		log.Fatalf("Unable to get name, ret: %v", ret)
	}
	fmt.Printf("Device Name: %s\n", name)

	index, ret := device.GetIndex()
	if ret != ixml.SUCCESS {
		log.Fatalf("Unable to get index, ret: %v", ret)
	}
	fmt.Printf("Device Index: %d\n", index)

	serialNumber, ret := device.GetSerial()
	if ret != ixml.SUCCESS {
		fmt.Printf("Unable to get GPU Serial Number , ret: %v\n", ret)
	} else {
		fmt.Printf("Device Serial Number: %s\n", serialNumber)
	}

	minorNumber, ret := device.GetMinorNumber()
	if ret != ixml.SUCCESS {
		fmt.Printf("Unable to get GPU MinorNumber, ret: %v\n", ret)
	} else {
		fmt.Printf("Device MinorNumber: %d\n", minorNumber)
	}

	currentEccMode, pendingEccMode, ret := device.GetEccMode()
	if ret != ixml.SUCCESS {
		log.Fatalf("Unable to get ECC Mode, ret: %v", ret)
	}
	fmt.Printf("Current ECC Mode: %d\n", currentEccMode)
	fmt.Printf("Pending ECC Mode: %d\n", pendingEccMode)

	fmt.Println("========================================")
}
