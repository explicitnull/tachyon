package types

type Authentications struct {
	Items []Authentication
	Total int
}
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
