module github.com/sergioa/hystrix-go

go 1.22

//replace github.com/cactus/go-statsd-client => github.com/cactus/go-statsd-client/v4 v4.0.0

require (
	github.com/DataDog/datadog-go v4.8.3+incompatible
	github.com/cactus/go-statsd-client/v5 v5.1.0
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475
	github.com/smartystreets/goconvey v1.8.1
)

require (
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/smarty/assertions v1.15.0 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
)
