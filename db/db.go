package db

import (
	"fmt"
	"tac-gateway/options"

	log "github.com/sirupsen/logrus"

	"github.com/BurntSushi/toml"
)

func dbconf() string {
	var c options.Options
	_, err := toml.DecodeFile("/etc/ams.conf", &c)
	if err != nil {
		log.Fatal(err)
	}
	par := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable", c.DbUser, c.DbPassword, c.DbName, c.DbHost)
	return par
}
