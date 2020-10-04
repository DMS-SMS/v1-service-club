module club

go 1.13

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/VividCortex/mysqlerr v0.0.0-20200629151747-c28746d985dd
	github.com/go-playground/validator/v10 v10.3.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/hashicorp/consul/api v1.7.0
	github.com/jinzhu/gorm v1.9.16
	github.com/micro/go-micro/v2 v2.9.1
	github.com/stretchr/testify v1.6.1
	gorm.io/driver/mysql v1.0.1
	gorm.io/gorm v1.20.1
)
