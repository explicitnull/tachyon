package repository

import (
	"errors"
	"fmt"

	"github.com/aerospike/aerospike-client-go"
)

const namespace = "tacacs"

// type Metrics struct {
// 	count int
// 	total int
// }

// var setMap = make(map[string]Metrics)

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func printOK(format string, a ...interface{}) {
	fmt.Printf("ok: "+format+"\n", a...)
}

func printError(format string, a ...interface{}) {
	fmt.Printf("error: "+format+"\n", a...)
}

func extractString(bins aerospike.BinMap, bin string) (string, error) {
	valueI, ok := bins[bin]
	if !ok {
		return "", errors.New("failed to get value from BinMap for key " + bin)
	}

	value, ok := valueI.(string)
	if !ok {
		return "", errors.New("BinMap value is not string")
	}

	return value, nil
}

func extractInt(bins aerospike.BinMap, bin string) (int, error) {
	valueI, ok := bins[bin]
	if !ok {
		return 0, errors.New("failed to get value from BinMap for key " + bin)
	}

	value, ok := valueI.(int)
	if !ok {
		return 0, errors.New("BinMap value is not integer")
	}

	return value, nil
}

func getAllRecords(aclient *aerospike.Client, setName string) ([]*aerospike.Record, error) {
	policy := aerospike.NewScanPolicy()
	policy.RecordsPerSecond = 10000

	nodeList := aclient.GetNodes()

	records := make([]*aerospike.Record, 0)

	for _, node := range nodeList {
		recordset, err := aclient.ScanNode(policy, node, namespace, setName)
		if err != nil {
			return nil, err
		}

	L:
		for {
			select {
			case record := <-recordset.Records:
				if record == nil {
					break L
				}

				// metrics, exists := setMap[record.Key.SetName()]

				// if !exists {
				// 	metrics = Metrics{}
				// }
				// metrics.count++
				// metrics.total++
				// setMap[record.Key.SetName()] = metrics

				records = append(records, record)

			case <-recordset.Errors:
				return nil, err
			}
		}

		// for k, v := range setMap {
		// 	log.Println("Node ", node, " set ", k, " count: ", v.count)
		// 	v.count = 0
		// }
	}

	return records, nil
}
