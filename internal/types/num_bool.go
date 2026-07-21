package types

import (
	"errors"
	"fmt"
)

type NumBool bool

func (receiver NumBool) MarshalJSON() ([]byte, error) {
	var i int
	if receiver {
		i = 1
	}

	return []byte(fmt.Sprintf("%d", i)), nil
}

func (receiver *NumBool) UnmarshalJSON(data []byte) error {
	asString := string(data)
	if asString == "1" {
		*receiver = true
	} else if asString == "0" {
		*receiver = false
	} else {
		return errors.New(fmt.Sprintf("Boolean unmarshal error: invalid input %s", asString))
	}
	return nil
}
