# Copyright 2018 Tendermint. All rights reserved.
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

PACKAGES=$(shell go list ./... | grep -Ev 'vendor|importer|rpc/tester')
COMMIT_HASH := $(shell git rev-parse --short HEAD)
BUILD_FLAGS = -tags netgo -ldflags "-X github.com/cosmos/ethermint/version.GitCommit=${COMMIT_HASH}"
DOCKER_TAG = unstable
DOCKER_IMAGE = cosmos/ethermint
ETHERMINT_DAEMON_BINARY = emintd
ETHERMINT_CLI_BINARY = emintcli
GO_MOD=GO111MODULE=on

all: tools verify install

#######################
### Build / Install ###
#######################

build:
ifeq ($(OS),Windows_NT)
	${GO_MOD} go build $(BUILD_FLAGS) -o build/$(ETHERMINT_DAEMON_BINARY).exe ./cmd/emintd
	${GO_MOD} go build $(BUILD_FLAGS) -o build/$(ETHERMINT_CLI_BINARY).exe ./cmd/emintcli
else
	${GO_MOD} go build $(BUILD_FLAGS) -o build/$(ETHERMINT_DAEMON_BINARY) ./cmd/emintd/
	${GO_MOD} go build $(BUILD_FLAGS) -o build/$(ETHERMINT_CLI_BINARY) ./cmd/emintcli/
endif

install:
	${GO_MOD} go install $(BUILD_FLAGS) ./cmd/emintd
	${GO_MOD} go install $(BUILD_FLAGS) ./cmd/emintcli

clean:
	@rm -rf ./build ./vendor

update-tools:
	@echo "--> Installing golangci-lint..."
	@wget -O - -q https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b ./bin v1.23.8

verify:
	@echo "--> Verifying dependencies have not been modified"
	${GO_MOD} go mod verify


#######################
### Testing / Misc. ###
#######################

test: test-unit

test-unit:
	@${GO_MOD} go test -v --vet=off $(PACKAGES)

test-race:
	@${GO_MOD} go test -v --vet=off -race $(PACKAGES)

test-cli:
	@echo "NO CLI TESTS"

lint:
	@echo "--> Running ci lint..."
	GOBIN=$(PWD)/bin go run scripts/ci.go lint

test-import:
	@${GO_MOD} go test ./importer -v --vet=off --run=TestImportBlocks --datadir tmp \
	--blockchain blockchain --timeout=5m
	# TODO: remove tmp directory after test run to avoid subsequent errors

test-rpc:
	@${GO_MOD} go test -v --vet=off ./rpc/tester

it-tests:
	./scripts/integration-test-all.sh -q 1 -z 1 -s 10

godocs:
	@echo "--> Wait a few seconds and visit http://localhost:6060/pkg/github.com/cosmos/ethermint"
	godoc -http=:6060

docker:
	docker build -t ${DOCKER_IMAGE}:${DOCKER_TAG} .
	docker tag ${DOCKER_IMAGE}:${DOCKER_TAG} ${DOCKER_IMAGE}:latest
	docker tag ${DOCKER_IMAGE}:${DOCKER_TAG} ${DOCKER_IMAGE}:${COMMIT_HASH}

format:
	@echo "--> Formatting go files"
	@find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -w -s
	@find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs misspell -w


.PHONY: build install update-tools tools godocs clean format lint \
test-cli test-race test-unit test test-import

###############################################################################
###                                Protobuf                                 ###
###############################################################################

proto-all: proto-gen proto-lint proto-check-breaking

proto-gen:
	@./scripts/protocgen.sh

proto-lint:
	@buf check lint --error-format=json

# NOTE: should match the default repo branch 
proto-check-breaking:
	@buf check breaking --against-input '.git#branch=development'


TM_URL           = https://raw.githubusercontent.com/tendermint/tendermint/v0.33.3
GOGO_PROTO_URL   = https://raw.githubusercontent.com/regen-network/protobuf/cosmos
COSMOS_PROTO_URL = https://raw.githubusercontent.com/regen-network/cosmos-proto/master
SDK_PROTO_URL 	 = https://raw.githubusercontent.com/cosmos/cosmos-sdk/master

TM_KV_TYPES         = third_party/proto/tendermint/libs/kv
TM_MERKLE_TYPES     = third_party/proto/tendermint/crypto/merkle
TM_ABCI_TYPES       = third_party/proto/tendermint/abci/types
GOGO_PROTO_TYPES    = third_party/proto/gogoproto
COSMOS_PROTO_TYPES  = third_party/proto/cosmos-proto
SDK_PROTO_TYPES     = third_party/proto/cosmos-sdk/types
AUTH_PROTO_TYPES    = third_party/proto/cosmos-sdk/x/auth/types
VESTING_PROTO_TYPES = third_party/proto/cosmos-sdk/x/auth/vesting/types
SUPPLY_PROTO_TYPES  = third_party/proto/cosmos-sdk/x/supply/types

proto-update-deps:
	@mkdir -p $(GOGO_PROTO_TYPES)
	@curl -sSL $(GOGO_PROTO_URL)/gogoproto/gogo.proto > $(GOGO_PROTO_TYPES)/gogo.proto

	@mkdir -p $(COSMOS_PROTO_TYPES)
	@curl -sSL $(COSMOS_PROTO_URL)/cosmos.proto > $(COSMOS_PROTO_TYPES)/cosmos.proto

	@mkdir -p $(TM_ABCI_TYPES)
	@curl -sSL $(TM_URL)/abci/types/types.proto > $(TM_ABCI_TYPES)/types.proto
	@sed -i '' '8 s|crypto/merkle/merkle.proto|third_party/proto/tendermint/crypto/merkle/merkle.proto|g' $(TM_ABCI_TYPES)/types.proto
	@sed -i '' '9 s|libs/kv/types.proto|third_party/proto/tendermint/libs/kv/types.proto|g' $(TM_ABCI_TYPES)/types.proto

	@mkdir -p $(TM_KV_TYPES)
	@curl -sSL $(TM_URL)/libs/kv/types.proto > $(TM_KV_TYPES)/types.proto

	@mkdir -p $(TM_MERKLE_TYPES)
	@curl -sSL $(TM_URL)/crypto/merkle/merkle.proto > $(TM_MERKLE_TYPES)/merkle.proto

	@mkdir -p $(SDK_PROTO_TYPES)
	@curl -sSL $(SDK_PROTO_URL)/types/types.proto > $(SDK_PROTO_TYPES)/types.proto

	@mkdir -p $(AUTH_PROTO_TYPES)
	@curl -sSL $(SDK_PROTO_URL)/x/auth/types/types.proto > $(AUTH_PROTO_TYPES)/types.proto
	@sed -i '' '5 s|types/types.proto|third_party/proto/cosmos-sdk/types/types.proto|g' $(AUTH_PROTO_TYPES)/types.proto

	@mkdir -p $(VESTING_PROTO_TYPES)
	curl -sSL $(SDK_PROTO_URL)/x/auth/vesting/types/types.proto > $(VESTING_PROTO_TYPES)/types.proto
	@sed -i '' '5 s|types/types.proto|third_party/proto/cosmos-sdk/types/types.proto|g' $(VESTING_PROTO_TYPES)/types.proto
	@sed -i '' '6 s|x/auth/types/types.proto|third_party/proto/cosmos-sdk/x/auth/types/types.proto|g' $(VESTING_PROTO_TYPES)/types.proto

	@mkdir -p $(SUPPLY_PROTO_TYPES)
	curl -sSL $(SDK_PROTO_URL)/x/supply/types/types.proto > $(SUPPLY_PROTO_TYPES)/types.proto
	@sed -i '' '5 s|types/types.proto|third_party/proto/cosmos-sdk/types/types.proto|g' $(SUPPLY_PROTO_TYPES)/types.proto
	@sed -i '' '6 s|x/auth/types/types.proto|third_party/proto/cosmos-sdk/x/auth/types/types.proto|g' $(SUPPLY_PROTO_TYPES)/types.proto


.PHONY: proto-all proto-gen proto-lint proto-check-breaking proto-update-deps