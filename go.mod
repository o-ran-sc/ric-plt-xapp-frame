module gerrit.o-ran-sc.org/r/ric-plt/xapp-frame

go 1.12

require (
	gerrit.o-ran-sc.org/r/com/golog v0.0.2
	gerrit.o-ran-sc.org/r/ric-plt/alarm-go.git/alarm v0.5.0
	gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common v1.0.35
	gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities v1.0.35
	gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader v1.0.35
	gerrit.o-ran-sc.org/r/ric-plt/sdlgo v0.5.0
	github.com/fsnotify/fsnotify v1.4.7
	github.com/go-openapi/errors v0.19.3
	github.com/go-openapi/loads v0.19.4
	github.com/go-openapi/runtime v0.19.4
	github.com/go-openapi/spec v0.19.3
	github.com/go-openapi/strfmt v0.19.4
	github.com/go-openapi/swag v0.19.7
	github.com/go-openapi/validate v0.19.6
	github.com/golang/protobuf v1.3.4
	github.com/gorilla/mux v1.7.1
	github.com/jessevdk/go-flags v1.4.0
	github.com/prometheus/client_golang v0.9.3
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.5.1
	golang.org/x/net v0.0.0-20190827160401-ba9fcec4b297
	k8s.io/utils v0.0.0-20201110183641-67b214c5f920 // indirect
)

replace gerrit.o-ran-sc.org/r/ric-plt/sdlgo => gerrit.o-ran-sc.org/r/ric-plt/sdlgo.git v0.5.2

replace gerrit.o-ran-sc.org/r/com/golog => gerrit.o-ran-sc.org/r/com/golog.git v0.0.2

replace gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common => gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common v1.0.35

replace gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities => gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities v1.0.35

replace gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader => gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader v1.0.35

replace gerrit.o-ran-sc.org/r/ric-plt/alarm-go.git/alarm => gerrit.o-ran-sc.org/r/ric-plt/alarm-go.git/alarm v0.5.0
