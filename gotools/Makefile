# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
#

## GOTOOLS = golint govendor goimports protoc-gen-go ginkgo gomega
GOTOOLS = golint govendor goimports ginkgo gomega
GOTOOLS_BIN = $(patsubst %,$(GOPATH)/bin/%, $(GOTOOLS))

# go tool->path mapping
go.fqp.govendor  := github.com/kardianos/govendor
go.fqp.golint    := golang.org/x/lint/golint
go.fqp.goimports := golang.org/x/tools/cmd/goimports
go.fqp.ginkgo    := github.com/onsi/ginkgo/ginkgo
go.fqp.gomega    := github.com/onsi/gomega

go.benchdep.logrus := github.com/sirupsen/logrus
go.benchdep.logxi := github.com/mgutz/logxi/v1
# go.benchdep.zap := github.com/uber-go/zap

all: $(GOTOOLS_BIN) benchdep

# Special override for protoc-gen-go since we want to use the version vendored with the project
## gotool.protoc-gen-go:
## 	mkdir -p $(GOPATH)/src/github.com/golang/protobuf/
## 	cp -r $(GOPATH)/src/github.com/ewangplay/eventbus/vendor/github.com/golang/protobuf/ $(GOPATH)/src/github.com/golang/
## 	go install github.com/golang/protobuf/protoc-gen-go
## 	rm -rf $(GOPATH)/src/github.com/golang/protobuf

# Default rule for gotools uses the name->path map for a generic 'go get' style build
gotool.%:
	$(eval TOOL = ${subst gotool.,,${@}})
	@echo "Building $(TOOL)"
	go get ${go.fqp.${TOOL}}

$(GOPATH)/bin/%:
	$(eval TOOL = ${subst $(GOPATH)/bin/,,${@}})
	@$(MAKE) gotool.$(TOOL)

# get benchmak dependency
benchdep:
	@echo "Downloading benchmak dependency..."
	go get ${go.benchdep.logrus}
	go get ${go.benchdep.logxi}
	# go get ${go.benchdep.zap}

.PHONY: clean
clean:
	-@rm -f $(GOTOOLS_BIN) ||:
