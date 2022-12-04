package postcard

type Address struct {
	Name  string `anonymize:"true"`
	Line1 string
	Line2 string
	Line3 string
}
