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
	ret := ixml.AbsInit("/usr/local/corex/lib/libixml.so")
	if ret != ixml.SUCCESS {
		log.Fatalf("Unable to initialize IXML: %v", ret)
		return
	}
	defer func() {
		ret := ixml.Shutdown()
		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to shutdown IXML: %v", ret)
		}
	}()

	count, ret := ixml.DeviceGetCount()
	if ret != ixml.SUCCESS {
		log.Fatalf("Unable to get device count %v", ret)
	}
	for i := uint(0); i < count; i++ {
		var device ixml.Device
		ret := ixml.DeviceGetHandleByIndex(i, &device)
		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to get device at index %d: %v", i, ret)
		}

		processInfos, ret := device.GetComputeRunningProcesses()
		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to get processInfos %d: %v", i, ret)
		}

		fmt.Printf("found processInfos %d on device %d\n", len(processInfos), i)
		for _, processInfo := range processInfos {
			fmt.Printf("processInfo.Pid: %d\n", processInfo.Pid)
			fmt.Printf("processInfo.Name: %s\n", processInfo.Name)
			fmt.Printf("processInfo.UsedGpuMemory: %d\n", processInfo.UsedGpuMemory)
		}
	}

	fmt.Println("========================================")
}
