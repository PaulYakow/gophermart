package entity

import (
	"database/sql/driver"
	"fmt"
)

type NullString string
type NullFloat32 float32

func (s *NullString) Scan(value interface{}) error {
	if value == nil {
		*s = ""
		return nil
	}

	val, ok := value.(string)
	if !ok {
		return fmt.Errorf("not a string")
	}

	*s = NullString(val)
	return nil
}

func (s NullString) Value() (driver.Value, error) {
	if len(s) == 0 {
		return nil, nil
	}

	return string(s), nil
}

func (f *NullFloat32) Scan(value interface{}) error {
	if value == nil {
		*f = 0.0
		return nil
	}

	val, ok := value.(float32)
	if !ok {
		return fmt.Errorf("not a string")
	}

	*f = NullFloat32(val)
	return nil
}

func (f NullFloat32) Value() (driver.Value, error) {
	if f == 0.0 {
		return nil, nil
	}

	return float32(f), nil
}
