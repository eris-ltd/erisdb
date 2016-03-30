package logger

import (
	cfg "github.com/eris-ltd/eris-db/tendermint/tendermint/config"
)

var config cfg.Config = nil

func init() {
	cfg.OnConfig(func(newConfig cfg.Config) {
		config = newConfig
		Reset() // reset log root upon config change.
	})
}
