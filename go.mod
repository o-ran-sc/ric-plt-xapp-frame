module gerrit.o-ran-sc.org/r/ric-plt/xapp-frame

go 1.12

require (
	gerrit.o-ran-sc.org/r/com/golog v0.0.1
	gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common v1.0.21
	gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities v1.0.21
	gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader v1.0.21
	gerrit.o-ran-sc.org/r/ric-plt/sdlgo v0.3.1
	github.com/fsnotify/fsnotify v1.4.7
	github.com/golang/protobuf v1.3.2
	github.com/gorilla/mux v1.7.1
	github.com/prometheus/client_golang v0.9.3
	github.com/spf13/viper v1.4.0
)

replace gerrit.o-ran-sc.org/r/ric-plt/sdlgo => gerrit.o-ran-sc.org/r/ric-plt/sdlgo.git v0.3.1

replace gerrit.o-ran-sc.org/r/com/golog => gerrit.o-ran-sc.org/r/com/golog.git v0.0.0-20190604083303-aaffc8ebe3f1

replace gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common => gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common v1.0.21

replace gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities => gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities v1.0.21

replace gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader => gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader v1.0.21
