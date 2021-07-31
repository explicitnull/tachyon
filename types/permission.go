package types

type Permission struct {
	Id          int
	Name        string
	Description string
	CreatedBy   string
}

type TemplatePermission struct {
	Id        int
	Name      string
	Descr     string
	Active    bool
	ActiveSt  string
	CreaTime  string // Full time form
	CreaTimeS string // Short time form
	CreaBy    string
}
