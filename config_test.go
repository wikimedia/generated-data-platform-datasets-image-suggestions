package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	log "gerrit.wikimedia.org/r/mediawiki/services/servicelib-golang/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaults(t *testing.T) {
	var config *Config
	var err error

	config, err = NewConfig([]byte{})

	require.NoError(t, err)
	require.NotNil(t, config)

	assert.Equal(t, "service-scaffold-golang", config.ServiceName)
	assert.Equal(t, "/", config.BaseURI)
	assert.Equal(t, "localhost", config.Address)
	assert.Equal(t, 8080, config.Port)
	assert.Equal(t, "info", strings.ToLower(config.LogLevel))
	assert.Equal(t, 9042, config.Cassandra.Port)
	assert.Equal(t, "quorum", strings.ToLower(config.Cassandra.Consistency))
	assert.Len(t, config.Cassandra.Hosts, 1)
	assert.Equal(t, "localhost", config.Cassandra.Hosts[0])
}

func TestFullConfig(t *testing.T) {
	var err error
	var obj *Config
	var conf string = `
service_name: test-service
base_uri: /v1
listen_address: 127.0.0.5
listen_port: 8081
log_level: debug
cassandra:
    port: 9043
    consistency: localQuorum
    hosts:
        - 127.0.0.6
        - 127.0.0.7
    local_dc: datacenter1
    authentication:
        username: eevans
        password: qwerty
    tls:
        ca: /tmp/ca/rootCa.pem
        cert: /tmp/ca/cert.pem
        key: /tmp/ca/key.pem
`
	obj, err = NewConfig([]byte(conf))

	require.NoError(t, err)
	require.NotNil(t, obj)

	assert.Equal(t, "test-service", obj.ServiceName)
	assert.Equal(t, "/v1", obj.BaseURI)
	assert.Equal(t, "127.0.0.5", obj.Address)
	assert.Equal(t, 8081, obj.Port)
	assert.Equal(t, "debug", strings.ToLower(obj.LogLevel))
	assert.Equal(t, 9043, obj.Cassandra.Port)
	assert.Equal(t, "localquorum", strings.ToLower(obj.Cassandra.Consistency))
	assert.Len(t, obj.Cassandra.Hosts, 2)
	assert.Contains(t, obj.Cassandra.Hosts, "127.0.0.6")
	assert.Contains(t, obj.Cassandra.Hosts, "127.0.0.7")
	assert.Equal(t, "datacenter1", obj.Cassandra.LocalDC)
	assert.Equal(t, "eevans", obj.Cassandra.Authentication.Username)
	assert.Equal(t, "qwerty", obj.Cassandra.Authentication.Password)
	assert.Equal(t, "/tmp/ca/rootCa.pem", obj.Cassandra.TLS.CaPath)
	assert.Equal(t, "/tmp/ca/cert.pem", obj.Cassandra.TLS.CertPath)
	assert.Equal(t, "/tmp/ca/key.pem", obj.Cassandra.TLS.KeyPath)
}

func TestValidConsistencies(t *testing.T) {
	var conf = `
cassandra:
    consistency: %s
`
	var consistencyLevels = []string{
		"any",
		"one",
		"two",
		"three",
		"quorum",
		"all",
		"eachquorum",
		"localquorum",
		"localone",
		"QuOruM",
		"localONE",
	}

	for _, consistency := range consistencyLevels {
		t.Run(consistency, func(t *testing.T) {
			_, err := NewConfig([]byte(fmt.Sprintf(conf, consistency)))
			require.NoError(t, err)
		})
	}
}

func TestBogusConsistency(t *testing.T) {
	var conf string = `
cassandra:
    consistency: unreal
`

	_, err := NewConfig([]byte(conf))
	require.Error(t, err)
}

func TestValidLogLevels(t *testing.T) {
	for _, level := range []string{"debug", "info", "warning", "error", "fatal", "FaTaL", "INFO"} {
		t.Run(level, func(t *testing.T) {
			_, err := NewConfig([]byte(fmt.Sprintf("log_level: %s", level)))
			require.NoError(t, err)
		})
	}
}

func TestBogusLogLevel(t *testing.T) {
	_, err := NewConfig([]byte("log_level: unreal"))
	require.Error(t, err)
}

func TestMergeEnvironment(t *testing.T) {
	config, err := NewConfig([]byte("cassandra: {authentication: {username: eevans, password: qwerty}}"))
	require.NoError(t, err)

	logger, err := log.NewLogger(os.Stdout, config.ServiceName, config.LogLevel)
	require.NoError(t, err)

	// Verify starting values...
	require.Equal(t, "eevans", config.Cassandra.Authentication.Username)
	require.Equal(t, "qwerty", config.Cassandra.Authentication.Password)

	// Username w/o password should default to the config values
	os.Setenv("CASSANDRA_USERNAME", "newuser")
	MergeEnvironment(config, logger)

	require.Equal(t, "eevans", config.Cassandra.Authentication.Username)
	require.Equal(t, "qwerty", config.Cassandra.Authentication.Password)

	// Password w/o username should default to the config values
	os.Unsetenv("CASSANDRA_USERNAME")
	os.Setenv("CASSANDRA_PASSWORD", "newpass")
	MergeEnvironment(config, logger)

	require.Equal(t, "eevans", config.Cassandra.Authentication.Username)
	require.Equal(t, "qwerty", config.Cassandra.Authentication.Password)

	// Username & password in environment should override config values
	os.Setenv("CASSANDRA_USERNAME", "newuser")
	MergeEnvironment(config, logger)

	assert.Equal(t, "newuser", config.Cassandra.Authentication.Username)
	assert.Equal(t, "newpass", config.Cassandra.Authentication.Password)
}
