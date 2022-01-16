// insert into tacacs.options(PK, name, value) values ('min_pass_len', 'min_pass_len', 9)

package repository

import (
	"tachyon/types"

	"github.com/aerospike/aerospike-client-go"
	"github.com/sirupsen/logrus"
)

const optionsSet = "options"

func GetOptions(le *logrus.Entry, aclient *aerospike.Client) ([]*types.Option, error) {
	recs, err := getAllRecords(aclient, optionsSet)
	if err != nil {
		return nil, err
	}

	result := make([]*types.Option, 0)

	for _, v := range recs {
		opt, err := extractOption(v.Bins)
		if err != nil {
			le.WithError(err).Error("extracting bins failed")
			return nil, err
		}

		result = append(result, opt)
	}

	return result, nil
}

func SetOptionMinimumPasswordLength(le *logrus.Entry, name string, length int) error {
	return nil
}

func extractOption(bins aerospike.BinMap) (*types.Option, error) {
	opt := &types.Option{}

	var err error

	opt.Name, err = extractString(bins, "name")
	if err != nil {
		return nil, err
	}

	opt.Value, err = extractString(bins, "value")
	if err != nil {
		return nil, err
	}

	return opt, nil
}
