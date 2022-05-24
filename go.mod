module gerrit.o-ran-sc.org/r/ric-plt/xapp-frame

go 1.12

require (
	gerrit.o-ran-sc.org/r/com/golog v0.0.2
	gerrit.o-ran-sc.org/r/ric-plt/alarm-go.git/alarm v0.5.1-0.20211223104552-f7d2cf80e85c
	gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common v1.2.1
	gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities v1.2.1
	gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader v1.2.1
	gerrit.o-ran-sc.org/r/ric-plt/sdlgo v0.7.0
	github.com/fsnotify/fsnotify v1.4.9
	github.com/go-openapi/errors v0.19.3
	github.com/go-openapi/loads v0.19.4
	github.com/go-openapi/runtime v0.19.4
	github.com/go-openapi/spec v0.19.3
	github.com/go-openapi/strfmt v0.19.4
	github.com/go-openapi/swag v0.19.7
	github.com/go-openapi/validate v0.19.6
	github.com/golang/protobuf v1.4.2
	github.com/gorilla/mux v1.7.1
	github.com/jessevdk/go-flags v1.4.0
	github.com/prometheus/client_golang v0.9.3
	github.com/prometheus/common v0.4.0
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.5.1
	golang.org/x/net v0.0.0-20200520004742-59133d7f0dd7
	k8s.io/utils v0.0.0-20201110183641-67b214c5f920 // indirect
)

replace gerrit.o-ran-sc.org/r/ric-plt/sdlgo => gerrit.o-ran-sc.org/r/ric-plt/sdlgo.git v0.7.0

replace gerrit.o-ran-sc.org/r/com/golog => gerrit.o-ran-sc.org/r/com/golog.git v0.0.2

replace gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common => gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common v1.2.1

replace gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities => gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities v1.2.1

replace gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader => gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader v1.2.1
