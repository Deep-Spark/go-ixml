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
	device1 := flag.Uint("device1", 0, "first GPU index")
	device2 := flag.Uint("device2", 1, "second GPU index")
	allPairs := flag.Bool("all-pairs", false, "run checks for all GPU pairs")
	flag.Parse()

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
	if count < 2 {
		log.Fatalf("Not enough GPUs for P2P test, need >= 2, got: %d", count)
	}
	fmt.Printf("GPU Count: %d\n", count)

	if *allPairs {
		for i := uint(0); i < count; i++ {
			for j := i + 1; j < count; j++ {
				runPairChecks(i, j)
			}
		}
		return
	}

	if *device1 >= count || *device2 >= count {
		log.Fatalf("device index out of range: device1=%d device2=%d count=%d", *device1, *device2, count)
	}
	if *device1 == *device2 {
		log.Fatalf("device1 and device2 must be different")
	}

	runPairChecks(*device1, *device2)
}

func runPairChecks(devIdx1, devIdx2 uint) {
	var device1, device2 ixml.Device

	ret := ixml.DeviceGetHandleByIndex(devIdx1, &device1)
	if ret != ixml.SUCCESS {
		fmt.Printf("GPU %d: DeviceGetHandleByIndex failed, ret=%v\n", devIdx1, ret)
		return
	}
	ret = ixml.DeviceGetHandleByIndex(devIdx2, &device2)
	if ret != ixml.SUCCESS {
		fmt.Printf("GPU %d: DeviceGetHandleByIndex failed, ret=%v\n", devIdx2, ret)
		return
	}

	name1, _ := ixml.DeviceGetName(device1)
	name2, _ := ixml.DeviceGetName(device2)

	fmt.Printf("\n===== Pair GPU %d(%s) <-> GPU %d(%s) =====\n", devIdx1, name1, devIdx2, name2)

	checkP2PStatus(device1, device2)
	checkIxLinkInfo(device1, device2)
}

func checkP2PStatus(device1, device2 ixml.Device) {
	caps := []struct {
		name  string
		index ixml.GpuP2PCapsIndex
	}{
		{name: "READ", index: ixml.P2P_CAPS_INDEX_READ},
		{name: "WRITE", index: ixml.P2P_CAPS_INDEX_WRITE},
		{name: "NVLINK", index: ixml.P2P_CAPS_INDEX_NVLINK},
		{name: "ATOMICS", index: ixml.P2P_CAPS_INDEX_ATOMICS},
		{name: "PROP", index: ixml.P2P_CAPS_INDEX_PROP},
	}

	for _, cap := range caps {
		status, ret := ixml.DeviceGetP2PStatus(device1, device2, cap.index)
		fmt.Printf("DeviceGetP2PStatus(%s): status=%v ret=%v\n", cap.name, status, ret)
	}
}

func checkIxLinkInfo(device1, device2 ixml.Device) {
	count, port, remotePort, ret := ixml.DeviceGetIxLinkInfo(device1, device2)
	fmt.Printf("DeviceGetIxLinkInfo: count=%d port=%v remote_port=%v ret=%v\n", count, port, remotePort, ret)
}
