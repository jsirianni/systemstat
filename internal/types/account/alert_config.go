package account

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/pkg/errors"
)

type alertConfig map[string]interface{}

func (a alertConfig) Value() (driver.Value, error) {
	j, err := json.Marshal(a)
	return j, err
}

func (a *alertConfig) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed.")
	}

	var i interface{}
	err := json.Unmarshal(source, &i)
	if err != nil {
		return err
	}

	*a, ok = i.(map[string]interface{})
	if !ok {
		return errors.New("Type assertion .(map[string]interface{}) failed.")
	}

	return nil
}

func (a alertConfig) JSON() ([]byte, error) {
	return json.Marshal(a)
}

func (a alertConfig) String() (string, error) {
	x, err := a.JSON()
	return string(x), err
}
