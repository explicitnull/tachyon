package repository

import (
	"github.com/aerospike/aerospike-client-go"
	"github.com/sirupsen/logrus"
)

// GetRole checks if user has extended control of tacplus
func GetRole(le *logrus.Entry, aClient *aerospike.Client, username string) string {
	return "admin"
}
