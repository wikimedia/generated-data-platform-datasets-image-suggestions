# Image Suggestions Data Gateway Service

Provides read-only Image Suggestions access (via the [cassandra-http-gateway] framework) for the
[Data Gateway Service][Data Gateway].


## Usage

Building requires:

- Go (>= 1.15)
- Make

```sh-session
$ make
...
$ ./image-suggestions -h
Usage of ./image-suggestions:
  -config string
    	Path to the configuration file (default "./config.yaml")
```

Running the service requires access to a [Cassandra] database.  The Image Suggestions schema can be
recreated using `cassandra_schema.cql` (ala `cqlsh -f cassandra_schema.cql`), but (for the time being)
test data is left as an exercise for the reader.


## Configuration

See the well-commented example, `config.yaml`.


## Deployment

See: https://wikitech.wikimedia.org/wiki/Kubernetes/Deployments


[Cassandra]:              http://cassandra.apache.org
[cassandra-http-gateway]: https://gitlab.wikimedia.org/repos/generated-data-platform/cassandra-http-gateway
[Data Gateway]:           https://www.mediawiki.org/wiki/Platform_Engineering_Team/Data_Value_Stream/Data_Gateway
[T293807]:                https://phabricator.wikimedia.org/T293807 
