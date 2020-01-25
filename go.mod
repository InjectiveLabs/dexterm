module github.com/InjectiveLabs/dexterm

go 1.13

require (
	github.com/InjectiveLabs/injective-core/api/gen/http/relayer/client v0.0.0-00010101000000-000000000000
	github.com/InjectiveLabs/injective-core/api/gen/relayer v0.0.0-00010101000000-000000000000
	github.com/InjectiveLabs/zeroex-go v0.0.0-20200125063848-29c3866c47f5
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5 // indirect
	github.com/apex/log v1.1.1
	github.com/c-bata/go-prompt v0.2.3
	github.com/ethereum/go-ethereum v1.9.10
	github.com/fatih/color v1.9.0
	github.com/go-playground/locales v0.13.0
	github.com/gogo/protobuf v1.2.0
	github.com/google/uuid v1.1.1
	github.com/graph-gophers/graphql-go v0.0.0-20191115155744-f33e81362277
	github.com/jawher/mow.cli v1.1.0
	github.com/magiconair/properties v1.8.1
	github.com/mailru/easyjson v0.0.0-20190626092158-b2ccc519800e
	github.com/mattn/go-tty v0.0.3 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/multiformats/go-multiaddr-dns v0.2.0 // indirect
	github.com/oklog/ulid v1.3.1
	github.com/olekukonko/tablewriter v0.0.4
	github.com/openzipkin/zipkin-go v0.2.2 // indirect
	github.com/pelletier/go-toml v1.6.0
	github.com/pkg/errors v0.9.1
	github.com/pkg/term v0.0.0-20190109203006-aa71e9d9e942 // indirect
	github.com/serialx/hashring v0.0.0-20190515033939-7706f26af194
	github.com/shopspring/decimal v0.0.0-20200105231215-408a2507e114
	github.com/sirupsen/logrus v1.4.2
	github.com/sony/gobreaker v0.4.1 // indirect
	github.com/streadway/handy v0.0.0-20190108123426-d5acb3125c2a // indirect
	github.com/stretchr/testify v1.4.0
	github.com/tj/go-spin v1.1.0
	github.com/xlab/closer v0.0.0-20190328110542-03326addb7c2
	github.com/xlab/structwalk v1.1.1
	goa.design/goa/v3 v3.0.9
	golang.org/x/crypto v0.0.0-20200117160349-530e935923ad
	google.golang.org/api v0.15.0 // indirect
	honnef.co/go/tools v0.0.0-20190523083050-ea95bdfd59fc
)

replace github.com/InjectiveLabs/injective-core/api/gen/http/relayer/client => ./gen/http/relayer/client

replace github.com/InjectiveLabs/injective-core/api/gen/relayer => ./gen/relayer
