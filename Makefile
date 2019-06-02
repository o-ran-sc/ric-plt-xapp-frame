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

ROOT_DIR:=.
BUILD_DIR:=$(ROOT_DIR)/build

COVEROUT := $(BUILD_DIR)/cover.out
COVERHTML := $(BUILD_DIR)/cover.html

GOOS=$(shell go env GOOS)
GOCMD=go
GOBUILD=$(GOCMD) build -a -installsuffix cgo
GOTEST=$(GOCMD) test -v -coverprofile $(COVEROUT)

GOFILES := $(shell find $(ROOT_DIR) -name '*.go' -not -name '*_test.go') go.mod go.sum
GOFILES_NO_VENDOR := $(shell find $(ROOT_DIR) -path ./vendor -prune -o -name "*.go" -not -name '*_test.go' -print)

APP:=$(BUILD_DIR)/xapp-sim
APPTST:=$(APP)_test

.PHONY: FORCE
 
.DEFAULT: build

default: build

$(APP): $(GOFILES)
	GO111MODULE=on GO_ENABLED=0 GOOS=linux $(GOBUILD) -o $@ ./test/xapp

$(APPTST): $(GOFILES)
	GO111MODULE=on GO_ENABLED=0 GOOS=linux $(GOTEST) -c -o $@ ./pkg/xapp 
	RMR_SEED_RT=config/uta_rtg.rt $@ -f config/config-file.yaml -test.coverprofile $(COVEROUT)
	go tool cover -html=$(COVEROUT) -o $(COVERHTML)

build: $(APP)

test: $(APPTST)

fmt: $(GOFILES_NO_VENDOR)
	gofmt -w -s $^
	@(RESULT="$$(gofmt -l $^)"; test -z "$${RESULT}" || (echo -e "gofmt failed:\n$${RESULT}" && false) )

clean:
	@echo "  >  Cleaning build cache"
	@-rm -rf $(APP) $(APPTST) 2> /dev/null
	go clean 2> /dev/null
