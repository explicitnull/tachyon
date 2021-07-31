package types

type Account struct {
	Name        string
	Password    string
	Subdivision string
	Permission  string
	Mail        string

	SubdivID int
	PermisID int

	Status string

	CreatedTimestamp string
	CreatedBy        string

	PasswordSetTimestamp string
}

type AccountTemplateData struct {
	Id        int
	Name      string
	Hash      string // Password has
	Cleartext string
	Subdiv    string
	Prm       string
	Mail      string

	Active     bool
	ActiveSt   string
	ActiveBox  string // Is "checked" or "" for HTML form
	CreaTime   string // Full time form
	CreaTimeS  string // Short time form
	CreaBy     string
	PassChd    bool
	SubdivList []string
	PrmList    []string
}

type TemplateAccountsSummary struct {
	Total  int
	Active int
}
