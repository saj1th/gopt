package gopt

type Marshaler interface {
	MarshalGoption(s string) error
}
