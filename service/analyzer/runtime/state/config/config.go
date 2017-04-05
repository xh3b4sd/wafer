package config

import (
	"github.com/xh3b4sd/wafer/service/analyzer/runtime/state/config/history"
)

type Config struct {
	History []history.History `json:"history"`
}
