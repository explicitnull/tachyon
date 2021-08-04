package types

const AccountStatusActive = "active"
const AccountStatusSuspended = "suspended"
const AccountStatusPasswordNotChanged = "tmp_pass_not_chd"

type Account struct {
	Name                     string
	Cleartext                string
	Hash                     string
	Subdivision              string
	Permission               string
	Mail                     string
	Status                   string
	CreatedTimestamp         string
	CreatedBy                string
	PasswordChangedTimestamp string // new
	// LastSignedInTimestamp    string // new
	UIRole string // new

	SubdivisionID    int // new
	PermissionID     int // new
	SubdivisionsList []string
	PermissionsList  []string
}

type Accounts struct {
	Items              []*Account
	Total              int
	Active             int
	PasswordNotChanged int
	Suspended          int
	NeverSignedIn      int

	TextStatusActive             string
	TextStatusPasswordNotChanged string
	TextStatusSuspended          string

	MoreRecords              bool
	TextAccountsPerPageLimit int
}
