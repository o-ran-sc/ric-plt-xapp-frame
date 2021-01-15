#!/bin/bash

#==================================================================================
#   Copyright (c) 2020 AT&T Intellectual Property.
#   Copyright (c) 2020 Nokia
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.
#==================================================================================

set -eux

echo "--> build_example_xapp.sh starts"

# Install RMR from deb packages at packagecloud.io
rmr=rmr_4.5.0_amd64.deb
wget --content-disposition  https://packagecloud.io/o-ran-sc/staging/packages/debian/stretch/$rmr/download.deb
sudo dpkg -i $rmr
rm $rmr
rmrdev=rmr-dev_4.5.0_amd64.deb
wget --content-disposition https://packagecloud.io/o-ran-sc/staging/packages/debian/stretch/$rmrdev/download.deb
sudo dpkg -i $rmrdev
rm $rmrdev

# Required to find nng and rmr libs
export LD_LIBRARY_PATH=/usr/local/lib

# Go install, build, etc
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:$PATH

# xApp-framework stuff
export CFG_FILE=config/config-file.json
export RMR_SEED_RT=config/uta_rtg.rt

GO111MODULE=on GO_ENABLED=0 GOOS=linux


# Build
go build -a -installsuffix cgo -o example_xapp cmd/example-xapp.go

echo "--> build_example_xapp.sh ends"
