package applogic

import (
	"tachyon/repository"
	"tachyon/types"
	"time"

	"github.com/aerospike/aerospike-client-go"
	"github.com/sirupsen/logrus"
)

const authOffset = 60

func ShowAuthentication(le *logrus.Entry, aclient *aerospike.Client) ([]types.Authentication, error) {
	now := time.Now()
	begin := now.Add(-authOffset * time.Minute)
	end := now

	le.Debugf("applogic begin: %s, end: %s", begin, end)

	items, err := repository.GetAuthenticationWithTimeFilter(le, aclient, begin, end)
	if err != nil {
		le.WithError(err).Error("getting authentication failed")
		return nil, err
	}

	return items, nil
}

func SearchAuthentication(le *logrus.Entry, field, value string, begin, end time.Time, aclient *aerospike.Client) []types.Authentication {
	items := make([]types.Authentication, 0)

	var err error

	if value != "" {
		items, err = repository.GetAuthenticationWithEqualFilter(le, aclient, field, value)
		if err != nil {
			le.WithError(err).Error("searching accounting failed")
			return nil
		}
	} else {
		items, err = repository.GetAuthenticationWithTimeFilter(le, aclient, begin, end)
		if err != nil {
			le.WithError(err).Error("searching accounting failed")
			return nil
		}
	}

	return items
}
