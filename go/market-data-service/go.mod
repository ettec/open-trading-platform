module github.com/ettech/open-trading-platform/go/market-data-service

go 1.13

require (
	github.com/BurntSushi/xgb v0.0.0-20200324125942-20f126ea2843 // indirect
	github.com/bradfitz/gomemcache v0.0.0-20190913173617-a41fca850d0b // indirect
	github.com/chzyer/logex v1.1.10 // indirect
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e // indirect
	github.com/chzyer/test v0.0.0-20180213035817-a1ea475d72b1 // indirect
	github.com/emicklei/go-restful-swagger12 v0.0.0-20170926063155-7524189396c6 // indirect
	github.com/ettec/open-trading-platform/go/common v0.0.0

	github.com/ettec/open-trading-platform/go/model v0.0.0
	github.com/garyburd/redigo v1.6.0 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/golang/protobuf v1.3.2
	github.com/google/uuid v1.1.1
	github.com/gorilla/schema v1.1.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.14.4 // indirect
	github.com/iancoleman/strcase v0.0.0-20191112232945-16388991a334 // indirect
	github.com/ianlancetaylor/demangle v0.0.0-20200414190113-039b1ae3a340 // indirect
	github.com/inconshreveable/go-vhost v0.0.0-20160627193104-06d84117953b // indirect
	github.com/jessevdk/go-flags v1.4.0 // indirect
	github.com/lyft/protoc-gen-star v0.4.14 // indirect
	github.com/mjibson/appstats v0.0.0-20151004071057-0542d5f0e87e // indirect
	github.com/nbutton23/zxcvbn-go v0.0.0-20180912185939-ae427f1e4c1d // indirect
	github.com/pkg/sftp v1.11.0 // indirect
	github.com/rogpeppe/go-charset v0.0.0-20190617161244-0dc95cdf6f31 // indirect
	github.com/syndtr/goleveldb v1.0.0 // indirect
	golang.org/x/image v0.0.0-20200119044424-58c23975cae1 // indirect
	golang.org/x/mobile v0.0.0-20200329125638-4c31acba0007 // indirect
	google.golang.org/grpc v1.25.1
	gopkg.in/vmihailenco/msgpack.v2 v2.9.1 // indirect
	k8s.io/api v0.17.4
	k8s.io/apimachinery v0.17.4
	k8s.io/client-go v0.17.4
)

replace github.com/ettec/open-trading-platform/go/common v0.0.0 => ../common

replace github.com/ettec/open-trading-platform/go/model v0.0.0 => ../model
