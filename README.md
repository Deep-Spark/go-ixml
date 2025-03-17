# GO-IXML

GO-IXML is Go bindings for the **IluvatarCorex Management Library API(IXML)**.

**Note:** 
- The runtime environment requires the library of **libixml.so**, please install IX SDK firstly.
- The current version of GO-IXML is compatible with IX driver version **4.2.0**.

## Build

The installation of GO-IXML is very simple, just execute the following command in the command line：
```shell
$ make bindings
```
## Sample
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
	fmt.Printf("GPU Count: %d\n", count)

	for i := 0; i < count; i++ {
		var device ixml.Device
		ret = ixml.DeviceGetHandleByIndex(i, &device)
		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to get device handle at index %d: %v", i, ret)
		}

		minor, ret := device.GetMinorNumber()
		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to get the minor of device at index %d: %v", i, ret)
		}
		fmt.Printf("Device Minor: %v\n", minor)

		uuid, ret := device.GetUUID()
		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to get the uuid of device at index %d: %v", i, ret)
		}
		fmt.Printf("Device Uuid: %s\n", uuid)

		name, ret := device.GetName()
		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to get the name of device %d: %v", i, ret)
		}
		fmt.Printf("Device Name: %v\n", name)

		driverVersion, ret := ixml.SystemGetDriverVersion()
		if ret != ixml.SUCCESS {
			log.Fatalf("Unable to get the system driver version: %v", ret)
		}
		fmt.Printf("System Driver Version: %v\n", driverVersion)
	}
}
```
## More Samples

The `samples` folder contains more simple examples of how to use GO-IXML to call the ixml API.

To get device attributes, run the following command:
```bash
go run samples/attributes/main.go
```

To get basic metrics of device, run the following command:
```bash
go run samples/metrics/main.go
```

To get gpm metrics of device, run the following command:
```bash
go run samples/gpmmetrics/main.go
```

To get running process information of device,run the following command:
```bash
go run samples/processinfo/main.go
```

To get system information such as driver version and CUDA version, run the following command:
```bash
go run samples/system/main.go
```

## License

Copyright (c) 2024 Iluvatar CoreX. All rights reserved. This project has an Apache-2.0 license, as
found in the [LICENSE](LICENSE) file.
