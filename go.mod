module github.com/brian-god/imcloud_sdk

go 1.18

require (
	github.com/golang/protobuf v1.5.4
	github.com/gorilla/websocket v1.4.2
	github.com/jinzhu/copier v0.3.5
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible // indirect
	github.com/pkg/errors v0.9.1
	google.golang.org/protobuf v1.33.0 // indirect
	nhooyr.io/websocket v1.8.10
)

require golang.org/x/net v0.22.0

require (
	github.com/OpenIMSDK/protocol v0.0.45
	github.com/OpenIMSDK/tools v0.0.24
	github.com/glebarez/sqlite v1.11.0
	github.com/google/go-cmp v0.6.0
	github.com/spark-comm/spark-api v0.0.0
	github.com/patrickmn/go-cache v2.1.0+incompatible
	golang.org/x/image v0.14.0
	golang.org/x/text v0.14.0
	gorm.io/gorm v1.25.7
)

replace github.com/spark-comm/spark-api => ../miliao_apis

require (
	github.com/bwmarrin/snowflake v0.3.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.0.4 // indirect
	github.com/glebarez/go-sqlite v1.21.2 // indirect
	github.com/go-kratos/kratos/v2 v2.7.3 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/lestrrat-go/strftime v1.0.6 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.24.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240318140521-94a12d6c2237 // indirect
	google.golang.org/grpc v1.64.0 // indirect
	modernc.org/libc v1.22.5 // indirect
	modernc.org/mathutil v1.5.0 // indirect
	modernc.org/memory v1.5.0 // indirect
	modernc.org/sqlite v1.23.1 // indirect
)
