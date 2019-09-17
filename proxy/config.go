package proxy

import (
	"github.com/hyperledger/burrow/keys"
)

type ProxyConfig struct {
	InternalProxyEnabled    bool
	ListenHost              string
	ListenPort              string
	AllowBadFilePermissions bool
	KeysDirectory           string
}

func DefaultProxyConfig() *ProxyConfig {
	return &ProxyConfig{
		// Default Monax keys port
		InternalProxyEnabled:    false,
		AllowBadFilePermissions: false,
		ListenHost:              "0.0.0.0",
		ListenPort:              "10998",
		KeysDirectory:           keys.DefaultKeysDir,
	}
}