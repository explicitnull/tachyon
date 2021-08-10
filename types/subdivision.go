package types

const SubdivisionStatusActive = "active"
const SubdivisionStatusInactive = "inactive"

type Subdivision struct {
	Id               int
	Name             string
	Description      string
	Status           string
	CreatedBy        string
	CreatedTimestamp string
}

type Subdivisions struct {
	Items []*Subdivision

	Total    int
	Active   int
	Inactive int
}
