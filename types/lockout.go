package types

type Lockout struct {
	IP                    string
	FQDN                  string
	Attempts              int // new
	FirstAttemptTimestamp string
	LastAttemptTimestamp  string
	LastDevice            string
	LastAccountName       string
}

type Lockouts struct {
	Items []*Lockout

	Total int
}
