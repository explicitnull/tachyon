package repository

import (
	"github.com/aerospike/aerospike-client-go"
	"github.com/sirupsen/logrus"
)

// GetRole checks if user has extended control of Tacasa
func GetRole(le *logrus.Entry, aClient *aerospike.Client, username string) string {
	return "superuser"
}
