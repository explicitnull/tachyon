package types

type Option struct {
	Name  string
	Value string
}

type Options struct {
	OptionItems []*Option
	TokenItems  []*Token
}
