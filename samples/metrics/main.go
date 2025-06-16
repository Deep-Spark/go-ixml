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
		log.Fatalf("Unable to initialize IXML, ret : %v", ret)
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
	fmt.Printf("GPU Count: %v\n", count)

	for i := uint(0); i < count; i++ {
		var device ixml.Device
		ret = ixml.DeviceGetHandleByIndex(i, &device)
		if ret != ixml.SUCCESS {
			fmt.Printf("Unable to get device at index %d, ret: %v\n", i, ret)
		} else {
			fmt.Printf("Get device at index %d\n", i)
		}

		integer, decimal, ret := device.GetGPUVoltage()
		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to get GPU Voltage, ret: %v", ret)
		}
		fmt.Printf("GPU Voltage of device %d: %v.%v\n", i, integer, decimal)

		temperature, ret := device.GetTemperature()
		if ret != ixml.SUCCESS {
			fmt.Printf("Unable to get GPU temperature of device %d, ret: %v\n", i, ret)
		} else {
			fmt.Printf("Temperature of device %d: %d\n", i, temperature)
		}

		fanSpeed, ret := device.GetFanSpeed()
		if ret != ixml.SUCCESS {
			fmt.Printf("Unable to get GPU FanSpeed of device %d, ret: %v\n", i, ret)
		} else {
			fmt.Printf("FanSpeed of device %d: %d\n", i, fanSpeed)
		}

		clockInfo, ret := device.GetClockInfo()
		if ret != ixml.SUCCESS {
			fmt.Printf("Unable to get GPU MemClock of device %d, ret: %v\n", i, ret)
		} else {
			fmt.Printf("MemClock of device %d: %d\n", i, clockInfo.Mem)
		}

		memoryInfo, ret := device.GetMemoryInfo()
		if ret != ixml.SUCCESS {
			fmt.Printf("Unable to get MemoryInfo of device %d, ret: %v\n", i, ret)
		} else {
			fmt.Printf("MemoryInfo totalMem of device %d: %d (MiB)\n", i, memoryInfo.Total)
			fmt.Printf("MemoryInfo usedMem of device %d: %d (MiB)\n", i, memoryInfo.Used)
			fmt.Printf("MemoryInfo freeMem of device %d: %v (MiB)\n", i, memoryInfo.Free)
		}

		utilizationRates, ret := device.GetUtilizationRates()
		if ret != ixml.SUCCESS {
			fmt.Printf("Unable to get UtilizationRates of device %d, ret: %v\n", i, ret)
		} else {
			fmt.Printf("Mem utilizationRates of device %d: %d\n", i, utilizationRates.Memory)
			fmt.Printf("GPU utilizationRates of device %d: %d\n", i, utilizationRates.Gpu)
		}

		pciInfo, ret := device.GetPciInfo()
		if ret != ixml.SUCCESS {
			fmt.Printf("Unable to get PciInfo of device %d, ret: %v\n", i, ret)
		} else {
			fmt.Printf("PciInfo of device %d: %v\n", i, pciInfo.BusId)
		}

		pcieGeneration, ret := device.GetCurrPcieLinkGeneration()
		if ret != ixml.SUCCESS {
			fmt.Printf("Unable to get PcieGeneration of device %d, ret: %v\n", i, ret)
		} else {
			fmt.Printf("PcieGeneration of device %d: %d\n", i, pcieGeneration)
		}

		pcieWidth, ret := device.GetCurrPcieLinkWidth()
		if ret != ixml.SUCCESS {
			fmt.Printf("Unable to get PcieWidth of device %d, ret: %v\n", i, ret)
		} else {
			fmt.Printf("PcieWidth of device %d: %d\n", i, pcieWidth)
		}

		pcieReplyCnt, ret := device.GetPcieReplayCounter()
		if ret != ixml.SUCCESS {
			fmt.Printf("Unable to get PcieReplayCounter of device %d, ret: %v\n", i, ret)
		} else {
			fmt.Printf("PcieReplayCounter of device %d: %v\n", i, pcieReplyCnt)
		}

		clocksThrottleReasons, ret := device.GetCurrentClocksThrottleReasons()
		if ret != ixml.SUCCESS {
			fmt.Printf("Unable to get GPU ClocksThrottleReasons of device %d, ret: %v\n", i, ret)
		} else {
			fmt.Printf("ClocksThrottleReasons of device %d: %v\n", i, clocksThrottleReasons)
		}

		singleErr, doubleErr, ret := device.GetEccErros()
		if ret != ixml.SUCCESS {
			fmt.Printf("Unable to get ecc errors %d, ret: %v\n", i, ret)
		} else {
			fmt.Printf("SingleEccErr: %d, DoubleEccErr: %d\n", singleErr, doubleErr)
		}

		usage, ret := device.GetPowerUsage()
		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to get Power Usage, ret: %v", ret)
		}
		fmt.Printf("Power Usage: %d\n", usage)

		limit, ret := device.GetPowerManagementLimit()
		if ret != ixml.SUCCESS {
			fmt.Printf("Unable to get PowerManagementLimit of device %d, ret: %v\n", i, ret)
		} else {
			fmt.Printf("Power Management Limit: %d\n", limit)
		}

		defaultLimit, ret := device.GetPowerManagementDefaultLimit()
		if ret != ixml.SUCCESS {
			fmt.Printf("Unable to get PowerManagementDefaultLimit of device %d, ret: %v\n", i, ret)
		} else {
			fmt.Printf("Default Power Management Limit: %d\n", defaultLimit)
		}

		minLimit, maxLimit, ret := device.GetPowerManagementLimitConstraints()
		if ret != ixml.SUCCESS {
			fmt.Printf("Unable to get PowerManagementLimitConstraints of device %d, ret: %v\n", i, ret)
		} else {
			fmt.Printf("MinPowerMgmtLimit: %d, MaxPowerMgmtLimit: %d\n", minLimit, maxLimit)
		}

		threshType := ixml.TEMPERATURE_THRESHOLD_SLOWDOWN
		threshVal, ret := device.GetTemperatureThreshold(threshType)
		if ret != ixml.SUCCESS {
			fmt.Printf("Unable to get TemperatureThreshold of device %d with type %d, ret: %v\n", i, threshType, ret)
		} else {
			fmt.Printf("The temperature threshold with type %d is: %d\n", threshType, threshVal)
		}

		fmt.Println("------------------------------------")
	}
}
