package types

type Account struct {
	Name                 string
	Cleartext            string
	Hash                 string
	Subdivision          string
	Permission           string
	Mail                 string
	Status               string
	CreatedTimestamp     string
	CreatedBy            string
	PasswordSetTimestamp string // new

	SubdivisionID    int // new
	PermissionID     int // new
	SubdivisionsList []string
	PermissionsList  []string
}

type Accounts struct {
	Items  []*Account
	Total  int
	Active int
}
