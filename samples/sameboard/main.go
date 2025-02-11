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

const defalutGpu1 = "GPU-6d2ec5fa-f293-57a3-9f2c-335f78120578"
const defalutGpu2 = "GPU-7edb0dc9-9291-5e13-9e1c-ad92672bdfec"

type Chip struct {
	name       string
	uuid       string
	index      uint
	Operations ixml.Device
}

func main() {
	// TODO: Add your code here.
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

	device1, ret := ixml.GetHandleByUUID(defalutGpu1)
	if ret != ixml.SUCCESS {
		log.Fatalf("Unable to get handle by uuid %v", ret)
	}

	name1, ret := device1.GetName()
	if ret != ixml.SUCCESS {
		log.Fatalf("Unable to get name of Device1: %v", ret)
	}
	fmt.Printf("name1: %s\n", name1)

	device2, ret := ixml.GetHandleByUUID(defalutGpu2)
	if ret != ixml.SUCCESS {
		log.Fatalf("Unable to get Handle by uuid %v", ret)
	}

	name2, ret := device2.GetName()
	if ret != ixml.SUCCESS {
		log.Fatalf("Unable to get name of Device2: %v", ret)
	}
	fmt.Printf("name2: %s\n", name2)

	OnSameBoard, ret := ixml.GetOnSameBoard(device1, device2)
	if ret == ixml.ERROR_NOT_SUPPORTED {
		fmt.Printf("GetOnSameBoard: Not supported\n")
	} else if ret != ixml.SUCCESS {
		log.Printf("Device1 and Device2 Not On Same Board: %v\n", ret)
	} else {
		fmt.Printf("Device1 and Device2 On Same Board: %d\n", OnSameBoard)
	}

	fmt.Println("========================================")
}
