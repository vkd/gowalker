package gowalker

import (
	"fmt"
)

// Config fills config structure.
func Config(cfg interface{}, setters ...Setter) error {
	for _, s := range setters {
		if i, ok := s.(Initer); ok {
			err := i.Init(cfg)
			if err != nil {
				return fmt.Errorf("init %T setter: %w", i, err)
			}
		}
	}

	fs := make(Fields, 0, 6)

	err := Walk(cfg, fs, MultiSetterOR(setters...))
	if err != nil {
		return fmt.Errorf("walk through config structure: %w", err)
	}

	return nil
}

type Initer interface {
	Init(cfg interface{}) error
}
