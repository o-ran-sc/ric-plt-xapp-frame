module gerrit.o-ran-sc.org/r/ric-plt/xapp-frame

go 1.12

require (
	gerrit.o-ran-sc.org/r/com/golog v0.0.0
	gerrit.o-ran-sc.org/r/ric-plt/sdlgo v0.1.1
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/fsnotify/fsnotify v1.4.7
	github.com/gorilla/mux v1.7.1
	github.com/prometheus/client_golang v0.9.3
	github.com/spf13/viper v1.3.2
	gitlabe1.ext.net.nokia.com/ric_dev/ue-nib v0.0.0-20190531121637-5379da45e235
)

replace gerrit.o-ran-sc.org/r/ric-plt/sdlgo => gerrit.o-ran-sc.org/r/ric-plt/sdlgo.git v0.1.1

replace gerrit.o-ran-sc.org/r/com/golog => gerrit.o-ran-sc.org/r/com/golog.git v0.0.0-20190604083303-aaffc8ebe3f1
