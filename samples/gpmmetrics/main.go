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
	"time"

	"gitee.com/deep-spark/go-ixml/pkg/ixml"
)

func main() {
	ret := ixml.AbsInit("/usr/local/corex/lib/libixml.so")
	if ret != ixml.SUCCESS {
		log.Fatalf("failed to initialize IXML: %v", ret)
	}
	defer func() {
		ret := ixml.Shutdown()
		if ret != ixml.SUCCESS {
			log.Fatalf("failed to shutdown IXML: %v", ret)
		}
	}()

	count, ret := ixml.DeviceGetCount()
	if ret != ixml.SUCCESS {
		log.Fatalf("failed to get device count: %v", ret)
	}
	for i := uint(0); i < count; i++ {
		if err := collectGPMMetrics(i); err != nil {
			log.Printf("failed to collect metrics for device %d: %v", i, err)
		}
	}

}

func collectGPMMetrics(i uint) error {
	var device ixml.Device
	ret := ixml.DeviceGetHandleByIndex(i, &device)
	if ret != ixml.SUCCESS {
		return fmt.Errorf("failed to get device at index %d: %v", i, ret)
	}

	gpuQuerySupport, ret := device.GpmQueryDeviceSupport()
	if ret != ixml.SUCCESS {
		return fmt.Errorf("failed to query GPM support: %v", ret)
	}
	if gpuQuerySupport.IsSupportedDevice == 0 {
		return fmt.Errorf("GPM queries are not supported")
	}
	fmt.Printf("GPM queries are supported for device %d\n", i)

	sample1, ret := ixml.GpmSampleAlloc()
	if ret != ixml.SUCCESS {
		return fmt.Errorf("failed to allocate GPM sample1: %v", ret)
	}
	defer func() {
		_ = sample1.Free()
	}()
	sample2, ret := ixml.GpmSampleAlloc()
	if ret != ixml.SUCCESS {
		return fmt.Errorf("failed to allocate GPM sample2: %v", ret)
	}
	defer func() {
		_ = sample2.Free()
	}()
	if ret := device.GpmSampleGet(sample1); ret != ixml.SUCCESS {
		return fmt.Errorf("failed to get GPM sample1: %v", ret)
	}
	time.Sleep(1 * time.Second)
	if ret := device.GpmSampleGet(sample2); ret != ixml.SUCCESS {
		return fmt.Errorf("failed to get GPM sample2: %v", ret)
	}

	gpmMetric := ixml.GpmMetricsGetType{
		NumMetrics: 3,
		Sample1:    sample1,
		Sample2:    sample2,
		Metrics: [98]ixml.GpmMetric{
			{
				MetricId: uint32(ixml.GPM_METRIC_SM_UTIL),
			},
			{
				MetricId: uint32(ixml.GPM_METRIC_SM_OCCUPANCY),
			},
			{
				MetricId: uint32(ixml.GPM_METRIC_DRAM_BW_UTIL),
			},
		},
	}

	ret = ixml.GpmMetricsGet(&gpmMetric)
	if ret != ixml.SUCCESS {
		return fmt.Errorf("failed to get gpm metric: %v", ret)
	}

	fmt.Printf("success to get gpm metric for device %d\n", i)
	for i := 0; i < int(gpmMetric.NumMetrics); i++ {
		if gpmMetric.Metrics[i].MetricId > 0 {
			fmt.Printf("gpm metric id: %+4v, value: %v\n", gpmMetric.Metrics[i].MetricId, gpmMetric.Metrics[i].Value)
		}
	}

	return nil
}
