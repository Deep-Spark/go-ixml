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
	"flag"
	"fmt"
	"log"

	"gitee.com/deep-spark/go-ixml/pkg/ixml"
)

func main() {
	libPath := flag.String("lib", "", "optional path to libixml.so; if empty, try common default paths")
	deviceIndex := flag.Uint("device", 0, "device index to test")
	allDevices := flag.Bool("all", false, "test all devices")
	runFastClear := flag.Bool("fast-clear", false, "also call DeviceFastClearDevice")
	runReset := flag.Bool("reset", false, "also call DeviceReset")
	flag.Parse()

	usedLibPath, ret := initIXML(*libPath)
	if ret != ixml.SUCCESS {
		log.Fatalf("Unable to initialize IXML: %v", ret)
	}
	fmt.Printf("Initialized IXML with: %s\n", usedLibPath)
	defer func() {
		ret := ixml.Shutdown()
		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to shutdown IXML: %v", ret)
		}
	}()

	count, ret := ixml.DeviceGetCount()
	if ret != ixml.SUCCESS {
		log.Fatalf("Unable to get device count: %v", ret)
	}
	if count == 0 {
		log.Fatal("No IXML devices found")
	}

	if *allDevices {
		for i := uint(0); i < count; i++ {
			testDevice(i, *runFastClear, *runReset)
		}
		return
	}

	if *deviceIndex >= count {
		log.Fatalf("Device index %d out of range, device count: %d", *deviceIndex, count)
	}
	testDevice(*deviceIndex, *runFastClear, *runReset)
}

func initIXML(libPath string) (string, ixml.Return) {
	candidates := []string{}
	if libPath != "" {
		candidates = append(candidates, libPath)
	}
	candidates = append(candidates,
		"/usr/local/corex/lib/libixml.so",
		"/usr/local/corex/lib64/libixml.so",
	)

	visited := make(map[string]struct{}, len(candidates))
	var lastRet ixml.Return
	for _, candidate := range candidates {
		if _, ok := visited[candidate]; ok {
			continue
		}
		visited[candidate] = struct{}{}

		ret := ixml.AbsInit(candidate)
		if ret == ixml.SUCCESS {
			return candidate, ret
		}
		lastRet = ret
	}
	return "", lastRet
}

func testDevice(index uint, runFastClear, runReset bool) {
	var device ixml.Device
	ret := ixml.DeviceGetHandleByIndex(index, &device)
	if ret != ixml.SUCCESS {
		fmt.Printf("Device %d: DeviceGetHandleByIndex ret=%v\n", index, ret)
		return
	}

	name, nameRet := ixml.DeviceGetName(device)
	if nameRet != ixml.SUCCESS {
		name = "unknown"
	}

	fmt.Printf("===== Device %d (%s) =====\n", index, name)

	busyStatus, ret := ixml.DeviceGetGpuBusyStatus(device)
	fmt.Printf("DeviceGetGpuBusyStatus: busy_status=%d ret=%v\n", busyStatus, ret)

	processInfos, ret := device.GetComputeRunningProcesses()
	fmt.Printf("device.GetComputeRunningProcesses: info_count=%d ret=%v\n", len(processInfos), ret)
	if ret == ixml.SUCCESS {
		for _, processInfo := range processInfos {
			fmt.Printf("  pid=%d name=%s used_gpu_memory=%d\n", processInfo.Pid, processInfo.Name, processInfo.UsedGpuMemory)
		}
	}

	if runFastClear {
		ret = ixml.DeviceFastClearDevice(device)
		fmt.Printf("DeviceFastClearDevice: ret=%v\n", ret)
	} else {
		fmt.Println("DeviceFastClearDevice: skipped, pass -fast-clear to run")
	}

	if runReset {
		ret = ixml.DeviceReset(device)
		fmt.Printf("DeviceReset: ret=%v\n", ret)
	} else {
		fmt.Println("DeviceReset: skipped, pass -reset to run")
	}
}
