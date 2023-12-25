module open_im_sdk

go 1.20

// go get -u github.com/OpenIMSDK/Open-IM-Server@main
require (
	github.com/antonfisher/nested-logrus-formatter v1.3.1
	github.com/golang/protobuf v1.5.3
	github.com/gorilla/websocket v1.5.0
	github.com/imCloud v0.0.0
	github.com/jinzhu/copier v0.3.5
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/mitchellh/mapstructure v1.5.0
	github.com/pkg/errors v0.9.1
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5
	github.com/shamsher31/goimgtype v1.0.0
	github.com/sirupsen/logrus v1.9.3
	github.com/tencentyun/qcloud-cos-sts-sdk v0.0.0-20220106031843-2efeb10ca2f6
	google.golang.org/protobuf v1.31.0 // indirect
	gorm.io/gorm v1.25.2-0.20230530020048-26663ab9bf55
	nhooyr.io/websocket v1.8.7
)

replace github.com/imCloud => ../imCloud

require (
	github.com/brian-god/xy-apis v0.0.20
	golang.org/x/net v0.19.0
	golang.org/x/text v0.14.0
	gorm.io/driver/sqlite v1.5.0
)

require (
	github.com/bwmarrin/snowflake v0.3.0 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.0.2 // indirect
	github.com/go-kratos/kratos/v2 v2.7.1 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/lestrrat-go/strftime v1.0.6 // indirect
	github.com/mattn/go-sqlite3 v1.14.17 // indirect
	github.com/montanaflynn/stats v0.0.0-20171201202039-1bf9dbcd8cbe // indirect
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646 // indirect
	github.com/shamsher31/goimgext v1.0.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	go.mongodb.org/mongo-driver v1.12.1 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	go.uber.org/zap v1.24.0 // indirect
	golang.org/x/crypto v0.16.0 // indirect
	golang.org/x/exp v0.0.0-20231219180239-dc181d75b848 // indirect
	golang.org/x/image v0.14.0 // indirect
	golang.org/x/mobile v0.0.0-20231127183840-76ac6878050a // indirect
	golang.org/x/mod v0.14.0 // indirect
	golang.org/x/sync v0.5.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/tools v0.16.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230822172742-b8732ec3820d // indirect
	google.golang.org/grpc v1.59.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
