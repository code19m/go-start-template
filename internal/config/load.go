package config

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
)

func Load() (*Config, error) {
	cfg := new(Config)

	workDir, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get working dir")
	}

	configDir := fmt.Sprintf("%s/configs/", workDir)

	// Read .env file from working directory
	err = cleanenv.ReadConfig(workDir+"/"+envFilename, cfg)
	if err != nil && !os.IsNotExist(err) {
		return nil, errors.Wrapf(err, "failed to read from %s:", envFilename)
	}

	// Read base config
	err = cleanenv.ReadConfig(configDir+baseCfgFilename, cfg)
	if err != nil && !isEOFerr(err) {
		return nil, errors.Wrapf(err, "failed to read from %s:", baseCfgFilename)
	}

	if !availableModes.Contains(cfg.AppMode) {
		cfg.AppMode = ProdMode
	}

	// Read application mode specific config file
	modeFilename, ok := cfgFileMapper[cfg.AppMode]
	if ok {
		err = cleanenv.ReadConfig(configDir+modeFilename, cfg)
		if err != nil && !isEOFerr(err) {
			return nil, errors.Wrapf(err, "failed to read from %s:", modeFilename)
		}
	}

	err = cfg.validate()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func isEOFerr(err error) bool {
	return strings.HasSuffix(err.Error(), io.EOF.Error())
}
