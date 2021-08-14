package types

type Authentication struct {
	ID          string // new
	Timestamp   string
	AccountName string
	DeviceIP    string
	DeviceName  string
	EventType   string
	UserIP      string
	UserFQDN    string // new
}

type Authentications struct {
	Items []Authentication

	Total int

	NotFound          bool
	MoreItems         bool
	ItemsPerPageLimit int
	SearchValue       string
}
