package types

type AccountingRecord struct {
	ID          string // new
	Timestamp   string
	DeviceIP    string
	DeviceName  string
	AccountName string
	UserIP      string
	UserFQDN    string // new
	Command     string
}

type AccountingRecords struct {
	Items []AccountingRecord

	Total int

	NotFound          bool
	MoreItems         bool
	ItemsPerPageLimit int

	SearchValue string
}
