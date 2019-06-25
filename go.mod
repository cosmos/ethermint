module github.com/cosmos/ethermint

go 1.12

require (
	github.com/allegro/bigcache v1.2.1 // indirect
	github.com/aristanetworks/goarista v0.0.0-20181101003910-5bb443fba8e0
	github.com/beorn7/perks v0.0.0-20180321164747-3a771d992973
	github.com/btcsuite/btcd v0.0.0-20190115013929-ed77733ec07d
	github.com/btcsuite/btcutil v0.0.0-20180706230648-ab6388e0c60a
	github.com/cosmos/cosmos-sdk v0.0.0-20181218000439-ec9c4ea543b5
	github.com/davecgh/go-spew v1.1.1
	github.com/deckarep/golang-set v1.7.1
	github.com/edsrzf/mmap-go v0.0.0-20170320065105-0bce6a688712
	github.com/ethereum/go-ethereum v1.8.27
	github.com/fsnotify/fsnotify v1.4.7
	github.com/go-kit/kit v0.8.0
	github.com/go-logfmt/logfmt v0.4.0
	github.com/go-stack/stack v1.8.0
	github.com/gogo/protobuf v1.1.1
	github.com/golang/protobuf v1.3.0
	github.com/golang/snappy v0.0.1
	github.com/google/uuid v1.0.0
	github.com/gorilla/websocket v1.4.0
	github.com/hashicorp/golang-lru v0.5.0
	github.com/hashicorp/hcl v1.0.0
	github.com/huin/goupnp v1.0.0
	github.com/inconshreveable/mousetrap v1.0.0
	github.com/jackpal/go-nat-pmp v1.0.1
	github.com/jmhodges/levigo v1.0.0
	github.com/karalabe/hid v0.0.0-20180420081245-2b4488a37358
	github.com/kr/logfmt v0.0.0-20140226030751-b84e30acd515
	github.com/magiconair/properties v1.8.0
	github.com/matttproud/golang_protobuf_extensions v1.0.1
	github.com/mitchellh/mapstructure v1.1.2
	github.com/pborman/uuid v0.0.0-20180906182336-adf5a7427709
	github.com/pelletier/go-toml v1.2.0
	github.com/pkg/errors v0.8.1
	github.com/pmezard/go-difflib v1.0.0
	github.com/prometheus/client_golang v0.9.2
	github.com/prometheus/client_model v0.0.0-20190129233127-fd36f4220a90
	github.com/prometheus/common v0.2.0
	github.com/prometheus/procfs v0.0.0-20190227231451-bbced9601137
	github.com/rcrowley/go-metrics v0.0.0-20180503174638-e2704e165165
	github.com/rjeczalik/notify v0.9.2
	github.com/rs/cors v1.6.0
	github.com/spf13/afero v1.2.1
	github.com/spf13/cast v1.3.0
	github.com/spf13/cobra v0.0.5
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.3.2
	github.com/stretchr/testify v1.3.0
	github.com/syndtr/goleveldb v0.0.0-20181105012736-f9080354173f
	github.com/tendermint/btcd v0.1.1
	github.com/tendermint/go-amino v0.15.0
	github.com/tendermint/iavl v0.12.2
	github.com/tendermint/tendermint v0.31.5
	golang.org/x/crypto v0.0.0-20190618222545-ea8f1a30c443 // indirect
	golang.org/x/net v0.0.0-20190404232315-eb5bcb51f2a3
	golang.org/x/sys v0.0.0-20190412213103-97732733099d
	golang.org/x/text v0.3.0
	google.golang.org/genproto v0.0.0-20181107211654-5fc9ac540362
	google.golang.org/grpc v1.19.0
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce
	gopkg.in/yaml.v2 v2.2.2
)

replace (
	github.com/cosmos/cosmos-sdk v0.0.0-20181218000439-ec9c4ea543b5 => ../../cosmos/cosmos-sdk
	github.com/ethereum/go-ethereum v1.8.27 => github.com/alexanderbez/go-ethereum v1.8.17-0.20181024144731-0a57b29f0c8e
	golang.org/x/crypto => github.com/tendermint/crypto v0.0.0-20180820045704-3764759f34a5
)
