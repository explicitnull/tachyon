package types

const TokenStatusActive = "active"
const TokenStatusInactive = "inactive"

type Token struct {
	ID               string
	Type             string
	Status           string
	Token            string
	CreatedBy        string
	CreatedTimestamp string
}

type Tokens struct {
	Items []*Token

	Total    int
	Active   int
	Inactive int
}
