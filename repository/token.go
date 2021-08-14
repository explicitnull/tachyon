// insert into tacacs.api_tokens(PK, id, type, status, token) values ('sidecar_m7', 'sidecar_m7', 'rw', 'active', 'vKUmF6j5gy3eobQAQ18m')

package repository

import (
	"tacacs-webconsole/types"

	"github.com/aerospike/aerospike-client-go"
	"github.com/sirupsen/logrus"
)

const tokensSet = "api_tokens"

func GetTokens(le *logrus.Entry, aclient *aerospike.Client) ([]*types.Token, error) {
	recs, err := getAllRecords(aclient, tokensSet)
	if err != nil {
		return nil, err
	}

	result := make([]*types.Token, 0)

	for _, v := range recs {
		perm, err := extractToken(v.Bins)
		if err != nil {
			le.WithError(err).Error("extracting bins failed")
			return nil, err
		}

		result = append(result, perm)
	}

	return result, nil
}

func extractToken(bins aerospike.BinMap) (*types.Token, error) {
	tk := &types.Token{}

	var err error

	tk.ID, err = extractString(bins, "id")
	if err != nil {
		return nil, err
	}

	tk.Type, err = extractString(bins, "type")
	if err != nil {
		return nil, err
	}

	tk.Status, err = extractString(bins, "status")
	if err != nil {
		return nil, err
	}

	tk.Token, err = extractString(bins, "token")
	if err != nil {
		return nil, err
	}

	// tk.CreatedBy, err = extractString(bins, "created_by")
	// if err != nil {
	// 	return nil, err
	// }

	// tk.CreatedTimestamp, err = extractString(bins, "created_ts")
	// if err != nil {
	// 	return nil, err
	// }

	return tk, nil
}
