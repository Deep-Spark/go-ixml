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
	} else if count == 0 {
		log.Fatalf("No GPUs found.")
	}
	fmt.Printf("GPU Chips Count: %v\n", count)

	boardGroups := make(map[uint32]struct{})
	for i := uint(0); i < count; i++ {
		var device ixml.Device
		ret = ixml.DeviceGetHandleByIndex(i, &device)
		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to get Handle by index, ret: %v", ret)
		}
		boardId, ret := device.GetBoardId()
		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to get Board Part Number, ret: %v", ret)
		}
		boardGroups[boardId] = struct{}{}
	}

	fmt.Printf("GPU Count: %d\n", len(boardGroups))
}
