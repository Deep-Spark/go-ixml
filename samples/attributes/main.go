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

const defalutGpu = "GPU-6d2ec5fa-f293-57a3-9f2c-335f78120578"

func main() {
	var device ixml.Device

	ret := ixml.Init()
	if ret != ixml.SUCCESS {
		log.Fatalf("Unable to initialize IXML: %v", ret)
	}
	defer func() {
		ret := ixml.Shutdown()
		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to shutdown IXML: %v", ret)
		}
	}()

	device, ret = ixml.GetHandleByUUID(defalutGpu)
	if ret != ixml.SUCCESS {
		log.Fatalf("Unable to get Handle by uuid %v", ret)
	}

	name, ret := device.GetName()
	if ret != ixml.SUCCESS {
		log.Fatalf("Unable to get name %v", ret)
	}
	fmt.Printf("name:%s, len(name): %d\n", name, len(name))

	index, ret := device.GetIndex()
	if ret != ixml.SUCCESS {
		log.Fatalf("Unable to get index %v", ret)
	}
	fmt.Printf("index: %d\n", index)

	Integer, Decimal, ret := device.GetGPUVoltage()
	if ret != ixml.SUCCESS {
		log.Fatalf("Unable to get GPU Voltage: %v", ret)
	}
	fmt.Printf("GPU Voltage: %v.%v\n", Integer, Decimal)

	pos, ret := device.GetBoardPosition()
	if ret == ixml.ERROR_NOT_SUPPORTED {
		fmt.Printf("GetBoardPosition interface is not supported\n")
	} else if ret != ixml.SUCCESS {
		log.Fatalf("Unable to get BoardPosition %v", ret)
	} else {
		fmt.Printf("position: %d\n", pos)
	}

	usage, ret := device.GetPowerUsage()
	if ret != ixml.SUCCESS {
		log.Fatalf("Unable to get usage %v", ret)
	}
	fmt.Printf("usage: %d\n", usage)

	fmt.Println("========================================")
}
