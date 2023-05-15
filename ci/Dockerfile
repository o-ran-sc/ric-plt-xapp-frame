#   Copyright (c) 2019 AT&T Intellectual Property.
#   Copyright (c) 2019 Nokia.
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

#----------------------------------------------------------
#
#----------------------------------------------------------

FROM nexus3.o-ran-sc.org:10002/o-ran-sc/bldr-ubuntu20-c-go:1.1.0 as xapp-base


RUN apt-get update -y \
    &&apt-get install -y \
    apt-utils \
    cmake \
    gawk \
    sudo \
    nano \
    jq \
    gettext-base \
    bison \
    flex \
    curl \
    tree

RUN curl -s https://packagecloud.io/install/repositories/o-ran-sc/master/script.deb.sh | bash

# RMR
ARG RMRVERSION=4.9.0
#RUN apt-get install -y rmr=${RMRVERSION} rmr-dev=${RMRVERSION}
RUN wget --content-disposition https://packagecloud.io/o-ran-sc/release/packages/debian/stretch/rmr_${RMRVERSION}_amd64.deb/download.deb && dpkg -i rmr_${RMRVERSION}_amd64.deb
RUN wget --content-disposition https://packagecloud.io/o-ran-sc/release/packages/debian/stretch/rmr-dev_${RMRVERSION}_amd64.deb/download.deb && dpkg -i rmr-dev_${RMRVERSION}_amd64.deb
RUN rm -f rmr_${RMRVERSION}_amd64.deb rmr-dev_${RMRVERSION}_amd64.deb

#
RUN ldconfig

#
ENV PATH="/usr/local/go/bin:${PATH}"

WORKDIR "/"
CMD ["/bin/bash"]

#----------------------------------------------------------
#
#----------------------------------------------------------
FROM xapp-base as xapp-base-testbuild

RUN mkdir -p /ws
WORKDIR "/ws"

# Module prepare (if go.mod/go.sum updated)
COPY go.mod /ws
COPY go.sum /ws
RUN go mod download

# Copy
COPY . /ws

# Build
RUN make -C /ws go-build

# Test fmt
RUN make -C /ws go-test-fmt

# Alternative way to run tests:
#   If there is ipv6 in use, then there might be problems if
#   "localhost" is resolved as ipv6 address
#   --> Comment out ipv6 hosts file entries before running tests
#RUN sed -r  "s/^(::1.*)/#\1/" /etc/hosts  > /etc/hosts.new \
#    && cat /etc/hosts.new > /etc/hosts \
#    && cat /etc/hosts  \
#    && make -C /ws go-test

RUN  make -C /ws go-test
