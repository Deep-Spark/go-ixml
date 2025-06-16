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
			log.Fatalf("Unable to get device at index %d, ret: %v", i, ret)
		}

		// Supported Event Types:
		//   ixml.EventTypeSingleBitEccError   (1)
		//   ixml.EventTypeDoubleBitEccError   (2)
		//   ixml.EventTypeXidCriticalError    (8)
		//   ixml.EventTypeClock               (16)
		supportTypes, ret := device.GetSupportedEventTypes()
		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to get supported event types, ret: %v", ret)
		} else {
			log.Printf("Successfully retrieved supported event types: %v\n", supportTypes)
		}

		set, ret := ixml.EventSetCreate()
		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to create event set, ret: %v", ret)
		} else {
			log.Printf("Successfully created event set: %v\n", set)
		}

		// Register event types to the event set, supported types can be found above.
		eventTypes := uint64(ixml.EventTypeXidCriticalError | ixml.EventTypeClock)
		ret = device.RegisterEvents(eventTypes, set)
		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to register events, ret: %v", ret)
		} else {
			log.Printf("Successfully registered events: %v\n", ret)
		}

		timeoutms := uint32(10)
		eventData, ret := ixml.EventSetWait(set, timeoutms) // or set.Wait(timeoutms)
		if ret != ixml.SUCCESS && ret != ixml.ERROR_TIMEOUT {
			log.Fatalf("EventSetWait failed, ret: %v", ret)
		} else if ret == ixml.ERROR_TIMEOUT {
			log.Printf("EventSetWait timed out after %d ms.\n", timeoutms)
		} else {
			log.Printf("Successfully received the event data: %+v \n", eventData)
		}

		ret = ixml.EventSetFree(set) // or set.Free()
		log.Printf("EventSetFree: %v\n", ret)

		fmt.Println("------------------------------------")
	}

}
