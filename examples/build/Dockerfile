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
#
#   This source code is part of the near-RT RIC (RAN Intelligent Controller)
#   platform project (RICP).
#==================================================================================

FROM nexus3.o-ran-sc.org:10002/o-ran-sc/bldr-ubuntu18-c-go:1.9.0 as ubuntu-example-xapp

# Install utilities
RUN apt update && apt install -y iputils-ping net-tools curl sudo

# Install RMr shared library & development header files
RUN wget --content-disposition https://packagecloud.io/o-ran-sc/staging/packages/debian/stretch/rmr_3.5.0_amd64.deb/download.deb && dpkg -i rmr_3.5.0_amd64.deb && rm -rf rmr_3.5.0_amd64.deb
RUN wget --content-disposition https://packagecloud.io/o-ran-sc/staging/packages/debian/stretch/rmr-dev_3.5.0_amd64.deb/download.deb && dpkg -i rmr-dev_3.5.0_amd64.deb && rm -rf rmr-dev_3.5.0_amd64.deb

# Install dependencies, compile and test the module
RUN mkdir -p /go/src/example-xapp
COPY . /go/src/example-xapp

WORKDIR "/go/src/example-xapp"
#RUN ./build/build_example_xapp.sh

ENV CFG_FILE=config/config-file.json
ENV RMR_SEED_RT=config/uta_rtg.rt
ENV GO111MODULE=on GO_ENABLED=0 GOOS=linux

RUN go build -a -installsuffix cgo -o example_xapp cmd/example-xapp.go

# Final, executable and deployable container
FROM ubuntu:18.04

COPY --from=ubuntu-example-xapp /go/src/example-xapp/build/run_example_xapp.sh /
COPY --from=ubuntu-example-xapp /go/src/example-xapp/example_xapp /
COPY --from=ubuntu-example-xapp /go/src/example-xapp/config/* /
COPY --from=ubuntu-example-xapp /usr/local/lib /usr/local/lib

RUN ldconfig

RUN chmod 755 /run_example_xapp.sh
CMD /run_example_xapp.sh
