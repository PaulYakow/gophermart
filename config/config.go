package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/spf13/pflag"
)

type Cfg struct {
	Address        string `env:"RUN_ADDRESS" env-default:"localhost:8080"`
	Dsn            string `env:"DATABASE_DSN" env-default:""`
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS" env-default:""`
}

var address = struct {
	name         string
	shorthand    string
	value        *string
	defaultValue string
}{
	"address",
	"a",
	new(string),
	"localhost:8080",
}

var dsn = struct {
	name         string
	shorthand    string
	value        *string
	defaultValue string
}{
	"dsn",
	"d",
	new(string),
	"",
}

var accrualAddress = struct {
	name         string
	shorthand    string
	value        *string
	defaultValue string
}{
	"accrualAddress",
	"r",
	new(string),
	"",
}

func (cfg *Cfg) updateCfgFromFlags() {
	address.value = pflag.StringP(address.name, address.shorthand, address.defaultValue, "address of server in host:port format")
	dsn.value = pflag.StringP(dsn.name, dsn.shorthand, dsn.defaultValue, "DSN for database connect")
	accrualAddress.value = pflag.StringP(accrualAddress.name, accrualAddress.shorthand, accrualAddress.defaultValue, "address of accrual system")

	pflag.Parse()

	cfg.Address = *address.value
	cfg.Dsn = *dsn.value
	cfg.AccrualAddress = *accrualAddress.value
}

func NewConfig() (*Cfg, error) {
	cfg := &Cfg{}

	cfg.updateCfgFromFlags()

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
