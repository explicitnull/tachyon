package types

const PermissionStatusActive = "active"
const PermissionStatusInactive = "inactive"

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

	Total    int
	Active   int
	Inactive int
}
