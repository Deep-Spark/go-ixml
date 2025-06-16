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
	fmt.Printf("GPU Count: %v\n", count)

	if err := CheckOnSameBoard(count); err != nil {
		fmt.Println(err)
	}

	for i := uint(0); i < count; i++ {
		var device ixml.Device
		ret = ixml.DeviceGetHandleByIndex(i, &device) // Get the first GPU device
		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to get Handle by index, ret: %v", ret)
		}
		fmt.Printf("Device Index: %d\n", i)

		pos, ret := device.GetBoardPosition()
		if ret == ixml.ERROR_NOT_SUPPORTED {
			fmt.Printf("Unable to get Board Position, ret: %v\n", ret)
		} else if ret != ixml.SUCCESS {
			log.Fatalf("Unable to get Board Position, ret: %v", ret)
		} else {
			fmt.Printf("Board Position: %d\n", pos)
		}

		boardId, ret := device.GetBoardId()
		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to get Board ID, ret: %v", ret)
		}
		fmt.Printf("Board ID: %d\n", boardId)

		boardPartNumber, ret := device.GetBoardPartNumber()
		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to get Board Part Number, ret: %v", ret)
		}
		fmt.Printf("Board Part Number: %s\n", boardPartNumber)

		fmt.Println("----------------------------------------")
	}

}

func CheckOnSameBoard(gpuCount uint) error {
	fmt.Println("Check if the first two GPUs are on the same board...")
	if gpuCount < 2 {
		return fmt.Errorf("not enough GPUs to check on same board, gpu count: %d", gpuCount)
	}

	devIdx1, devIdx2 := uint(0), uint(1)
	var device1, device2 ixml.Device

	ret := ixml.DeviceGetHandleByIndex(devIdx1, &device1) // Get the first GPU device
	if ret != ixml.SUCCESS {
		return fmt.Errorf("failed to get handle by index, ret: %v", ret)
	}
	ret = ixml.DeviceGetHandleByIndex(devIdx2, &device2) // Get the second GPU device
	if ret != ixml.SUCCESS {
		return fmt.Errorf("failed to get Handle by index, ret: %v", ret)
	}

	onSameBoard, ret := ixml.GetOnSameBoard(device1, device2)
	if ret == ixml.ERROR_NOT_SUPPORTED {
		return fmt.Errorf("nvmlDeviceOnSameBoard: ERROR_NOT_SUPPORTED")
	} else if ret != ixml.SUCCESS {
		return fmt.Errorf("GPU %d and %d are NOT on same board: %v", devIdx1, devIdx2, ret)
	} else {
		fmt.Printf("GPU %d and %d on same board: %d\n", devIdx1, devIdx2, onSameBoard)
	}
	fmt.Println("----------------------------------------")
	return nil
}
