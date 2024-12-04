# GO-IXML

GO-IXML is Go bindings for the **IluvatarCorex Management Library API(ixml)**.

**Note:** The runtime environment requires the library of **libixml.so**, please install **IluvatarCorex** SDK firstly.

## Build

The installation of GO-IXML is very simple, just execute the following command in the command line：
```shell
$ make bindings
```

## Samples

The following is a simple example of how to use GO-IXML to call the ixml API:

```go
package main

import (
	"fmt"
	"log"

	"gitee.com/deep-spark/go-ixml/pkg/ixml"
)

func main() {
	// ret := ixml.Init()
	ret := ixml.AbsInit("/usr/local/corex/lib/libixml.so")
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

	for i := 0; i < count; i++ {
		var device ixml.Device
		ret = ixml.DeviceGetHandleByIndex(i, &device)

		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to get device at index %d: %v", i, ret)
		}

		ID, ret := device.GetMinorNumber()
		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to get id of device at index %v: %v", ID, ret)
		}

		fmt.Printf("ID: %v\n", ID)

		uuid, ret := device.GetUUID()
		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to get uuid of device at uuid %d: %v", i, ret)
		}
		fmt.Printf("len(uuid): %v, uuid: %v\n", len(uuid), uuid)

		name, ret := device.GetName()
		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to get name of device %d: %v", i, ret)
		}
		fmt.Printf("Device Name: %v\n", name)

		hardWareVersion, ret := device.GetHardWareVersion()
		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to get hardware version of device %d: %v", i, ret)
		}
		fmt.Printf("hardWare Version: %v\n", hardWareVersion)

		integer, decimal, ret := device.GetGPUVoltage()
		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to get GPU voltage of device %d: %v", i, ret)
		}
		fmt.Printf("GPU Voltage: %v.%v\n", integer, decimal)

		pn, ret := device.GetBoardPn()
		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to get Board Pn of device %d: %v", i, ret)
		}
		fmt.Printf("pn: %s\n", pn)

		fmt.Printf("========================================\n")
	}
}
```

## License

Copyright (c) 2024 Iluvatar CoreX. All rights reserved. This project has an Apache-2.0 license, as
found in the [LICENSE](LICENSE) file.
