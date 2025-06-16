/*
Copyright (c) 2020, NVIDIA CORPORATION.  All rights reserved.
Copyright (c) 2024, Shanghai Iluvatar CoreX Semiconductor Co., Ltd.
All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ixml

// EventSet represents the interface for the nvmlEventSet type.
type EventSet interface {
	Free() Return
	Wait(uint32) (EventData, Return)
}

// EventData includes an interface type for Device instead of nvmlDevice
type EventData struct {
	Device            Device
	EventType         uint64
	EventData         uint64
	GpuInstanceId     uint32
	ComputeInstanceId uint32
}

func (e nvmlEventData) convert() EventData {
	return EventData(e)
}

// ixml.EventSetCreate()
func EventSetCreate() (EventSet, Return) {
	var Set nvmlEventSet
	ret := nvmlEventSetCreate(&Set)
	return Set, ret
}

// ixml.EventSetWait()
func EventSetWait(set EventSet, timeoutms uint32) (EventData, Return) {
	return set.Wait(timeoutms)
}

// Wait waits for a registered event in the EventSet for up to timeoutms milliseconds.
func (set nvmlEventSet) Wait(timeoutms uint32) (EventData, Return) {
	var data nvmlEventData
	ret := nvmlEventSetWait(set, &data, timeoutms)
	return data.convert(), ret
}

// ixml.EventSetFree()
func EventSetFree(set EventSet) Return {
	return set.Free()
}

func (set nvmlEventSet) Free() Return {
	return nvmlEventSetFree(set)
}
