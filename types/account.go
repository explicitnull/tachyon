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
	UILevel string // new

	SubdivisionID    int // new
	PermissionID     int // new
	SubdivisionsList []string
	PermissionsList  []string
}

type Accounts struct {
	Items []*Account

	Total              int
	Active             int
	PasswordNotChanged int
	Suspended          int
	NeverSignedIn      int

	MoreItems         bool
	ItemsPerPageLimit int
}
