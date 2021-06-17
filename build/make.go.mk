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


#------------------------------------------------------------------------------
#
#------------------------------------------------------------------------------
ifndef ROOT_DIR
$(error ROOT_DIR NOT DEFINED)
endif
ifndef CACHE_DIR
$(error CACHE_DIR NOT DEFINED)
endif

#------------------------------------------------------------------------------
#
#------------------------------------------------------------------------------

GO_CACHE_DIR?=$(abspath $(CACHE_DIR)/go)

#------------------------------------------------------------------------------
#
#------------------------------------------------------------------------------
ifndef MAKE_GO_TARGETS
MAKE_GO_TARGETS:=1


.PHONY: FORCE go-build go-test go-test-fmt go-fmt go-clean
 
FORCE:


GOOS=$(shell go env GOOS)
GOCMD=go
GOBUILD=$(GOCMD) build -a -installsuffix cgo
GORUN=$(GOCMD) run -a -installsuffix cgo
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test -v
GOGET=$(GOCMD) get

GOFILES:=$(shell find $(ROOT_DIR) -name '*.go' -not -name '*_test.go')
GOALLFILES:=$(shell find $(ROOT_DIR) -name '*.go')
GOMODFILES:=go.mod go.sum


.SECONDEXPANSION:
$(GO_CACHE_DIR)/%: $(GOFILES) $(GOMODFILES) $$(BUILDDEPS)
	@echo "Building:\t$*"
	@eval GO111MODULE=on GOSUMDB=off GO_ENABLED=0 GOOS=linux $(BUILDARGS) $(GOBUILD) -o $@ ./$*


.SECONDEXPANSION:
$(GO_CACHE_DIR)/%_test: $(GOALLFILES) $(GOMODFILES) $$(BUILDDEPS) FORCE
	@echo "Testing:\t$*"
	@eval GO111MODULE=on GOSUMDB=off GO_ENABLED=0 GOOS=linux $(BUILDARGS) $(GOTEST) -cover -coverprofile=coverage.out -c -o $@ ./$*
	@if test -e $@ ; then eval $(TESTENV) $@ -test.v -test.coverprofile $(COVEROUT); else true ; fi
	@if test -e $@ ; then go tool cover -html=$(COVEROUT) -o $(COVERHTML); else true ; fi


.SECONDEXPANSION:
go-build: GO_TARGETS:=
go-build: $$(GO_TARGETS)

.SECONDEXPANSION:
go-test: GO_TARGETS:=
go-test: $$(GO_TARGETS)

go-test-fmt: $(GOFILES)
	@(RESULT="$$(gofmt -l $^)"; test -z "$${RESULT}" || (echo -e "gofmt failed:\n$${RESULT}" && false) )

go-fmt: $(GOFILES)
	gofmt -w -s $^

go-mod-tidy: FORCE
	GO111MODULE=on GOSUMDB=off go mod tidy

go-mod-download: FORCE
	GO111MODULE=on GOSUMDB=off go mod download

go-clean: GO_TARGETS:=
go-clean:
	@echo "  >  Cleaning build cache"
	@-rm -rf $(GO_TARGETS)* 2> /dev/null
	go clean 2> /dev/null


endif

#------------------------------------------------------------------------------
#
#-------------------------------------------------------------------- ----------

$(GO_CACHE_DIR)/$(XAPP_ROOT)/$(XAPP_NAME): BUILDDEPS:=$(XAPP_BUILDDEPS)
$(GO_CACHE_DIR)/$(XAPP_ROOT)/$(XAPP_NAME): BUILDARGS:=$(XAPP_BUILDARGS)


$(GO_CACHE_DIR)/$(XAPP_ROOT)/$(XAPP_NAME)_test: BUILDDEPS:=$(XAPP_BUILDDEPS)
$(GO_CACHE_DIR)/$(XAPP_ROOT)/$(XAPP_NAME)_test: BUILDARGS:=$(XAPP_BUILDARGS)
$(GO_CACHE_DIR)/$(XAPP_ROOT)/$(XAPP_NAME)_test: COVEROUT:=$(abspath $(GO_CACHE_DIR)/$(XAPP_ROOT)/$(XAPP_NAME)_cover.out)
$(GO_CACHE_DIR)/$(XAPP_ROOT)/$(XAPP_NAME)_test: COVERHTML:=$(abspath $(GO_CACHE_DIR)/$(XAPP_ROOT)/$(XAPP_NAME)_cover.html)
$(GO_CACHE_DIR)/$(XAPP_ROOT)/$(XAPP_NAME)_test: TESTENV:=$(XAPP_TESTENV)

go-build: GO_TARGETS+=$(GO_CACHE_DIR)/$(XAPP_ROOT)/$(XAPP_NAME)
go-test: GO_TARGETS+=$(GO_CACHE_DIR)/$(XAPP_ROOT)/$(XAPP_NAME)_test
go-clean: GO_TARGETS+=$(GO_CACHE_DIR)/$(XAPP_ROOT)/$(XAPP_NAME) $(GO_CACHE_DIR)/$(XAPP_ROOT)/$(XAPP_NAME)_test
