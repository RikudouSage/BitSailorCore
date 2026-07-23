package types

import (
	"errors"
	"strings"
)

type CSVSlice []string

func (receiver *CSVSlice) UnmarshalJSON(bytes []byte) error {
	str := string(bytes)
	if str == "null" {
		return nil
	}

	if str[0] != '"' && str[len(str)-1] != '"' {
		return errors.New("the value must be a string")
	}

	inner := str[1 : len(str)-1]
	*receiver = strings.Split(inner, ",")
	return nil
}

func (receiver CSVSlice) MarshalJSON() ([]byte, error) {
	return []byte("\"" + receiver.String() + "\""), nil
}

func (receiver CSVSlice) String() string {
	return strings.Join(receiver, ",")
}
