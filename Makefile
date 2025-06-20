# Copyright (c) 2024, Shanghai Iluvatar CoreX Semiconductor Co., Ltd.
# All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

MODULE := gitee.com/deep-spark/go-ixml

GEN_DIR = $(PWD)/gen
PKG_DIR = $(PWD)/pkg
GEN_BINDINGS_DIR = $(GEN_DIR)/ixml
PKG_BINDINGS_DIR = $(PKG_DIR)/ixml

SOURCES = $(shell find $(GEN_BINDINGS_DIR) -type f)

.PHONY: all test clean
.PHONY: bindings test-bindings clean-bindings

all: bindings

bindings: $(SOURCES)
	rm -rf $(PKG_BINDINGS_DIR)/{ixml,doc,const,cgo_helpers,types,types_gen}.go
	c-for-go -nostamp -out $(PKG_DIR) $(GEN_BINDINGS_DIR)/ixml.yml
	cp -f $(GEN_BINDINGS_DIR)/*.h $(PKG_BINDINGS_DIR)
	cp -f $(GEN_BINDINGS_DIR)/cgo_helpers.go $(PKG_BINDINGS_DIR)
	cd $(PKG_BINDINGS_DIR); \
		go tool cgo -godefs types.go > types_gen.go; \
		go fmt types_gen.go; \
	cd -> /dev/null
	rm -rf $(PKG_BINDINGS_DIR)/types.go $(PKG_BINDINGS_DIR)/_obj

COVERAGE_FILE := coverage.out
test: bindings
	go test -v -coverprofile=$(COVERAGE_FILE) $(MODULE)/pkg/...

coverage: test
	cat $(COVERAGE_FILE) | grep -v "_mock.go" > $(COVERAGE_FILE).no-mocks
	go tool cover -func=$(COVERAGE_FILE).no-mocks

clean:
	rm -rf $(PKG_BINDINGS_DIR)/{ixml,doc,const,cgo_helpers,types,types_gen}.go
	rm -rf $(PKG_BINDINGS_DIR)/types.go $(PKG_BINDINGS_DIR)/_obj