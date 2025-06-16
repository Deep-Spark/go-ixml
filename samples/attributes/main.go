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

// Replace with your actual GPU UUID
const defalutGpu = "GPU-6d2ec5fa-f293-57a3-9f2c-335f78120578"

func main() {
	var device ixml.Device

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

	fmt.Printf("Start to get attributes of device: %s\n", defalutGpu)
	device, ret = ixml.GetHandleByUUID(defalutGpu)
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

	uuid, ret := device.GetUUID()
	if ret != ixml.SUCCESS {
		fmt.Printf("Unable to get GPU Uuid , ret: %v\n", ret)
	} else {
		fmt.Printf("Device Uuid: %s\n", uuid)
	}

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
