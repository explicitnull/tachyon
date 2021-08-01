package repository

import (
	"tacasa-web/types"

	"github.com/aerospike/aerospike-client-go"
	"github.com/sirupsen/logrus"
)

const subdivisionsSet = "subdivisions"

func GetSubdivisions(le *logrus.Entry, aclient *aerospike.Client) ([]*types.Subdivision, error) {
	recs, err := getAllRecords(aclient, subdivisionsSet)
	if err != nil {
		return nil, err
	}

	result := make([]*types.Subdivision, 0)

	for _, v := range recs {
		subdiv, err := extractSubdivision(v.Bins)
		if err != nil {
			le.WithError(err).Error("extracting bins failed")
			return nil, err
		}

		result = append(result, subdiv)
	}

	return result, nil
}

// 	GetSubdivisionID returns subdivision ID for DB normalization
func GetSubdivisionID(le *logrus.Entry, aClient *aerospike.Client, div string) (int, error) {
	return 1, nil
}

func GetSubdivisionsList(le *logrus.Entry, aclient *aerospike.Client) []string {
	return []string{"europe", "asia"}
}

func extractSubdivision(bins aerospike.BinMap) (*types.Subdivision, error) {
	subdiv := &types.Subdivision{}

	var err error

	subdiv.Name, err = extractString(bins, "name")
	if err != nil {
		return nil, err
	}

	subdiv.Description, err = extractString(bins, "description")
	if err != nil {
		return nil, err
	}

	subdiv.Status, err = extractString(bins, "status")
	if err != nil {
		return nil, err
	}

	subdiv.CreatedBy, err = extractString(bins, "created_by")
	if err != nil {
		return nil, err
	}

	subdiv.CreatedTimestamp, err = extractString(bins, "created_ts")
	if err != nil {
		return nil, err
	}

	return subdiv, nil
}
