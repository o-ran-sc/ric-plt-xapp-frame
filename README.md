# XAPP-FRAME

## Introduction
**xapp-frame** is a simple framework for rapid development of RIC xapps, and supports various services essential for RIC xapps such as RESTful APIs, RMR (RIC Message Routing), database backend services and watching and populating config-map changes in K8S environment.

## Architecture

![Architecture](assets/xappframe-arch.png)

## Features and Components

* RESTful support
* Health check/probes (readiness and liveliness)
* Reading and watching config-map
* RMR messaging
* SDL
* Loggind and tracing
* Encoding and decoding of commonly used RIC ASN.1 messages
* And more to come

## Quick Start

#### Below is a simple example xapp. For more information, see the sample code in the xapp/examples folder:
```go
package main

import "gerrit.o-ran-sc.org/r/ric-plt/xapp-frame/pkg/xapp"

type ExampleXapp struct {
}

func (m *ExampleXapp) Consume(rp *xapp.RMRParams) (err error) {
        xapp.Logger.Debug("Message received - type=%d len=%d", rp.Mtype, rp.PayloadLen)

        xapp.Sdl.Store("myKey", rp.Payload)
        xapp.Rmr.Send(r)
        return nil
}

func main() {
       xapp.Run(ExampleXapp{})
}
```
#### Building and running the example xapp

    git clone https://gerrit.o-ran-sc.org/r/ric-plt/xapp-frame
    cd xapp-frame

#### To generate an executable binary for our sample xApp application, run the following command:

    GO111MODULE=on GO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o example-xapp examples/example-xapp.go

#### To run the generated executable binary locally, run the following command:

    RMR_SEED_RT=config/uta_rtg.rt ./example-xapp -f config/config-file.yaml

Congratulations! You've just built your first **xapp** application.

## API
#### API List
 * TBD

#### API Usage and Examples
* Setting logging level and writing to log
    ```
    xapp.Logger.SetLevel(4)
    xapp.Logger.Info("Status inquiry ...")
    ```
* Storing key-value data to SDL
    ```
    xapp.Sdl.Store("myKey", payload)
    ```
* Sending RMR messages
    ```
    mid := Rmr.GetRicMessageId("RIC_SUB_RESP")
    xapp.Rmr.Send(mid, 1234, len, payload)
    ```
* Injecting REST API resources (URL)
    ```
    xapp.Resource.InjectRoute("/ric/v1/health/stat", statisticsHandler, "GET")
    Resource.InjectQueryRoute("/ric/v1/user", handler, "GET", "foo", "bar", "id", "mykey")
    ```

## Documentation

## Community

## License
This project is licensed under the Apache License 2.0 - see the [LICENSE.md](LICENSE.md) file for details

