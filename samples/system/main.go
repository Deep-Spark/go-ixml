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
	"strconv"

	"gitee.com/deep-spark/go-ixml/pkg/ixml"
)

func main() {
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

	// Get the driver version
	version, ret := ixml.SystemGetDriverVersion()
	if ret != ixml.SUCCESS {
		log.Fatalf("Unable to get driver version: %v", ret)
	}
	fmt.Printf("Driver Version: %v\n", version)

	// Get the cuda version
	cudaVersion, ret := ixml.SystemGetCudaDriverVersion()
	if ret != ixml.SUCCESS {
		log.Fatalf("Unable to get cuda version: %v", ret)
	}
	cudaVersionInt, err := strconv.Atoi(cudaVersion)
	if err != nil {
		log.Fatalf("Failed to convert cuda version: %v", err)
	}
	major, minor := uint(cudaVersionInt)/1000, uint(cudaVersionInt)%1000/10
	fmt.Printf("Cuda Version: %d.%d\n", major, minor)

	// Get the IXML version
	ixmlVersion, ret := ixml.SystemGetNVMLVersion()
	if ret != ixml.SUCCESS {
		log.Fatalf("Unable to get IXML version: %v", ret)
	} else {
		fmt.Printf("IXML Version: %v\n", ixmlVersion)
	}

	fmt.Println("========================================")
}
