package types

type AccountingRec struct {
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
	Items []*AccountingRec

	Total int

	MoreItems         bool
	ItemsPerPageLimit int
}
