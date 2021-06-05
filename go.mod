module github.com/cosmos/ethermint

go 1.16

require (
	github.com/aristanetworks/goarista v0.0.0-20201012165903-2cb20defcd66 // indirect
	github.com/armon/go-metrics v0.3.8
	github.com/btcsuite/btcd v0.21.0-beta
	github.com/btcsuite/btcutil v1.0.2
	github.com/cespare/cp v1.1.1 // indirect
	github.com/cosmos/cosmos-sdk v0.42.5
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/ethereum/go-ethereum v1.9.25
	github.com/fjl/memsize v0.0.0-20190710130421-bcb5799ab5e5 // indirect
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.1.4 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/miguelmota/go-ethereum-hdwallet v0.0.0-20200123000308-a60dcd172b4c
	github.com/pkg/errors v0.9.1
	github.com/prometheus/tsdb v0.10.0 // indirect
	github.com/regen-network/cosmos-proto v0.3.1
	github.com/rs/zerolog v1.22.0
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	github.com/status-im/keycard-go v0.0.0-20200402102358-957c09536969
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.10
	github.com/tendermint/tm-db v0.6.4
	github.com/tyler-smith/go-bip39 v1.1.0
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
	google.golang.org/genproto v0.0.0-20210114201628-6edceaf6022f
	google.golang.org/grpc v1.38.0
	gopkg.in/yaml.v2 v2.4.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
