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

package ixml

import (
	"testing"

	"gitee.com/deep-spark/go-ixml/pkg/dl"
)

func requireLibIXML(t *testing.T) {
	lib := dl.New(ixmlLibraryName, ixmlLibraryLoadFlags)
	if err := lib.Open(); err != nil {
		t.Skipf("This test requires %v", ixmlLibraryName)
	}
	lib.Close()
}

func TestInit(t *testing.T) {
	requireLibIXML(t)

	ret := Init()
	if ret != SUCCESS {
		t.Errorf("Init: %v", ret)
	} else {
		t.Logf("Init: %v", ret)
	}

	ret = Shutdown()
	if ret != SUCCESS {
		t.Errorf("Shutdown: %v", ret)
	} else {
		t.Logf("Shutdown: %v", ret)
	}
}

func TestSystem(t *testing.T) {
	requireLibIXML(t)

	Init()
	defer Shutdown()

	driverVersion, ret := SystemGetDriverVersion()
	if ret != SUCCESS {
		t.Errorf("SystemGetDriverVersion: %v", ret)
	} else {
		t.Logf("SystemGetDriverVersion: %v", ret)
		t.Logf("Driver version: %v", driverVersion)
	}

	ixmlVersion, ret := SystemGetNVMLVersion()
	if ret != SUCCESS {
		t.Errorf("SystemGetNVMLVersion: %v", ret)
	} else {
		t.Logf("IXML version: %v", ixmlVersion)
	}

	cudaDriverVersion, ret := SystemGetCudaDriverVersion()
	if ret != SUCCESS {
		t.Errorf("SystemGetCudaDriverVersion: %v", ret)
	} else {
		t.Logf("Cuda driver version: %v", cudaDriverVersion)
	}

	cudaDriverVersionV2, ret := SystemGetCudaDriverVersion_v2()
	if ret != SUCCESS {
		t.Errorf("SystemGetCudaDriverVersion_v2: %v", ret)
	} else {
		t.Logf("SystemGetCudaDriverVersion_v2: %v", ret)
		t.Logf("Cuda driver version_v2: %v", cudaDriverVersionV2)
	}
}
