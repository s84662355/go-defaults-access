package access

type Unmarshaler interface {
	Default(val string) error
	Empty() bool
}
