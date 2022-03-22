module gitlab.wikimedia.org/repos/generated-data-platform/datasets/image-suggestions

go 1.15

replace gitlab.wikimedia.org/eevans/cassandra-http-gateway => /home/eevans/dev/src/git/cassandra/http-gateway

require (
	gerrit.wikimedia.org/r/mediawiki/services/servicelib-golang v0.0.0-20220322011350-df509f780b5c
	github.com/gocql/gocql v1.0.0
	github.com/golang/snappy v0.0.4 // indirect
	github.com/julienschmidt/httprouter v1.3.0
	github.com/prometheus/client_golang v1.12.1
	gitlab.wikimedia.org/eevans/cassandra-http-gateway v0.0.0-20220322001004-14949cbf268d
	golang.org/x/sys v0.0.0-20220319134239-a9b59b0215f8 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/yaml.v2 v2.4.0
)
