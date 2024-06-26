module github.com/openimsdk/openim-sdk-core/v3

go 1.18

require (
	github.com/golang/protobuf v1.5.4
	github.com/gorilla/websocket v1.4.2
	github.com/jinzhu/copier v0.3.5
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible // indirect
	github.com/pkg/errors v0.9.1
	google.golang.org/protobuf v1.33.0 // indirect
	gorm.io/driver/sqlite v1.3.6
	nhooyr.io/websocket v1.8.10
)

require golang.org/x/net v0.20.0

require (
	github.com/OpenIMSDK/protocol v0.0.45
	github.com/OpenIMSDK/tools v0.0.24
	github.com/google/go-cmp v0.5.9
	github.com/miliao_apis v0.0.0
	github.com/patrickmn/go-cache v2.1.0+incompatible
	golang.org/x/image v0.14.0
	gorm.io/gorm v1.23.8
)

replace github.com/miliao_apis => ../miliao_apis

require (
	github.com/bwmarrin/snowflake v0.3.0 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.0.4 // indirect
	github.com/go-kratos/kratos/v2 v2.7.3 // indirect
	github.com/go-playground/form/v4 v4.2.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/iancoleman/strcase v0.3.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/lestrrat-go/strftime v1.0.6 // indirect
	github.com/lyft/protoc-gen-star/v2 v2.0.3 // indirect
	github.com/mattn/go-sqlite3 v1.14.12 // indirect
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646 // indirect
	github.com/spf13/afero v1.10.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.24.0 // indirect
	golang.org/x/mod v0.12.0 // indirect
	golang.org/x/sync v0.3.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/tools v0.13.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230717213848-3f92550aa753 // indirect
	google.golang.org/grpc v1.56.3 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
