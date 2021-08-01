package options

import "github.com/BurntSushi/toml"

const defaultOptionsFilename = "ams.conf"

type Options struct {
	AuthenFile         string
	AcctFile           string
	DbHost             string
	DbName             string
	DbUser             string
	DbPassword         string
	FailInterval       int
	FailThold          int
	LogName            string
	TlsFullChain       string
	TlsPrivKey         string
	MinPassLen         int
	Maintenance        string
	EqpTimeBefore      string
	AerospikeNamespace string // new
}

func Load(o *Options) error {
	_, err := toml.DecodeFile(defaultOptionsFilename, &o)
	if err != nil {
		return err
	}

	return nil
}
