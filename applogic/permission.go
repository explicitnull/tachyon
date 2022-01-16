package applogic

import (
	"tachyon/repository"
	"tachyon/types"

	"github.com/aerospike/aerospike-client-go"
	"github.com/sirupsen/logrus"
)

func EditPermissionAction(le *logrus.Entry, aclient *aerospike.Client, fperm types.Permission) error {
	// getting account data from DB
	dbperm, err := repository.GetPermissionByName(le, aclient, fperm.Name)
	if err != nil {
		return err
	}

	// applying changes
	if fperm.Description != "" && fperm.Description != dbperm.Description {
		err = repository.SetPermissionDescription(le, aclient, fperm.Name, fperm.Description)
		if err != nil {
			return err
		}
	}

	if fperm.Status != dbperm.Status {
		err = repository.SetPermissionStatus(le, fperm.Name, fperm.Status)
		if err != nil {
			return err
		}
	}

	return nil
}
