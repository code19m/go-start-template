package config

import "go-start-template/pkg/ds"

// Application Modes
// To change the application mode, set the "APP_MODE environment variable.
// The default mode is "PRODUCTION".
const (
	LocalMode = "LOCAL"
	TestMode  = "TEST"
	StageMode = "STAGE"
	ProdMode  = "PRODUCTION"
)

var availableModes = ds.NewSet(LocalMode, TestMode, StageMode, ProdMode)

const (
	baseCfgFilename = "base.yaml"
	envFilename     = ".env"
)

var cfgFileMapper = map[string]string{
	LocalMode: "local.yaml",
	StageMode: "stage.yaml",
	ProdMode:  "prod.yaml",
}
