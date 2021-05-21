package access

type Unmarshaler interface {
	Default(val string) error
	IsNil() bool
}
