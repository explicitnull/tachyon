package types

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
}
