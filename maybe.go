package gowalker

type Maybe[T any] struct {
	IsSet bool
	Value T
}

func (m Maybe[T]) Get() (zero T, _ bool) {
	if m.IsSet {
		return m.Value, true
	}
	return zero, false
}

func (m *Maybe[T]) Set(v T) {
	m.IsSet = true
	m.Value = v
}

func (Maybe[T]) GoWalkerMaybe() {}
func (Maybe[T]) GoWalkerField() {}
