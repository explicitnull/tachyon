package applogic

import (
	"tacacs-webconsole/repository"

	"github.com/aerospike/aerospike-client-go"
	"github.com/sirupsen/logrus"
)

func LoginAction(le *logrus.Entry, username, formPassword string, aClient *aerospike.Client) bool {
	dbhash, err := repository.GetPasswordHash(le, aClient, username)
	if err != nil {
		le.WithError(err).Errorf("GetPasswordHash() failed")
		return false
	}

	if dbhash == "" {
		le.Warning("user not found")
		return false
	}

	// hashParts := strings.Split(dbhash, "$")
	// if len(hashParts) != 3 {
	// 	le.Error("wrong database hash format")
	// 	return false
	// }

	// salt := hashParts[2]
	formHash := applogic.MakeHash(le, formPassword)

	if formHash != dbhash {
		le.Warning("wrong password")
		return false
	}

	le.Info("logged in")

	return true
}
