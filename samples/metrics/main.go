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
		log.Fatalf("Unable to get device count: %v", ret)
	}
	fmt.Printf("GPU Count: %v\n", count)

	for i := uint(0); i < count; i++ {
		var device ixml.Device
		ret = ixml.DeviceGetHandleByIndex(i, &device)
		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to get device at index %d: %v", i, ret)
		}

		Uuid, ret := device.GetUUID()
		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to get GPU Uuid of device %d: %v", i, ret)
		}
		fmt.Printf("Uuid of device %d: %s\n", i, Uuid)

		MinorNumber, ret := device.GetMinorNumber()
		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to get GPU MinorNumber of device %d: %v", i, ret)
		}
		fmt.Printf("MinorNumber of device %d: %d\n", i, MinorNumber)

		temperature, ret := device.GetTemperature()
		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to get GPU temperature of device %d: %v", i, ret)
		}
		fmt.Printf("temperature of device %d: %d\n", i, temperature)

		FanSpeed, ret := device.GetFanSpeed()
		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to get GPU FanSpeed of device %d: %v", i, ret)
		}
		fmt.Printf("FanSpeed of device %d: %d\n", i, FanSpeed)

		ClockInfo, ret := device.GetClockInfo()
		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to get GPU MemClock of device %d: %v", i, ret)
		}
		fmt.Printf("MemClock of device %d: %d\n", i, ClockInfo.Mem)

		MemoryInfo, ret := device.GetMemoryInfo()
		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to get GPU MemoryInfo of device %d: %v", i, ret)
		}
		fmt.Printf("MemoryInfo totalMem of device %d: %d (MiB)\n", i, MemoryInfo.Total)
		fmt.Printf("MemoryInfo usedMem of device %d: %d (MiB)\n", i, MemoryInfo.Used)
		fmt.Printf("MemoryInfo freeMem of device %d: %v (MiB)\n", i, MemoryInfo.Free)

		utilizationRates, ret := device.GetUtilizationRates()
		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to get GPU utilizationRates of device %d: %v", i, ret)
		}
		fmt.Printf("Mem utilizationRates of device %d: %d\n", i, utilizationRates.Memory)
		fmt.Printf("GPU utilizationRates of device %d: %d\n", i, utilizationRates.Gpu)

		PciInfo, ret := device.GetPciInfo()
		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to get GPU PciInfo of device %d: %v", i, ret)
		}
		fmt.Printf("PciInfo of device %d: %v\n", i, PciInfo.BusId)

		clocksThrottleReasons, ret := device.GetCurrentClocksThrottleReasons()
		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to get GPU clocksThrottleReasons of device %d: %v", i, ret)
		}
		fmt.Printf("clocksThrottleReasons of device %d: %v\n", i, clocksThrottleReasons)

		singleErr, doubleErr, ret := device.GetEccErros()
		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to get ecc errors %d: %v", i, ret)
		}
		fmt.Printf("singleErr: %d, doubleErr: %d\n", singleErr, doubleErr)

		fmt.Println("========================================")
	}
}