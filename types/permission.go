package types

type Permission struct {
	Id               int
	Name             string
	Description      string
	Status           string
	CreatedBy        string
	CreatedTimestamp string
}

type Permissions struct {
	Items []*Permission
}
