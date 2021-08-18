package repository

import (
	"errors"

	"github.com/aerospike/aerospike-client-go"
	"github.com/sirupsen/logrus"
)

// GetAccessLevel checks if user can see or change tacplus configuration
func GetAccessLevel(le *logrus.Entry, aclient *aerospike.Client, username string) (string, error) {
	key, err := aerospike.NewKey(namespace, accountsSet, username)
	if err != nil {
		return "", err
	}

	policy := aerospike.NewPolicy()

	rec, err := aclient.Get(policy, key, "ui_role")
	if err != nil {
		le.WithError(err).Error("aerospike query failed")
		return "", err
	}

	if rec == nil {
		return "", errors.New("aerospike record not found")
	}

	level, err := extractString(rec.Bins, "ui_role")
	if err != nil {
		return "", err
	}

	le.Debugf("access level - %s", level)

	return level, nil
}
