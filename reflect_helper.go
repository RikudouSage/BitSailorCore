package bitwarden

import (
	"errors"
	"reflect"
)

var errValIsNil = errors.New("field value is nil")

func getStringValue(field reflect.StructField, target any) (string, error) {
	var strVal string
	if field.Type.Kind() == reflect.Pointer && field.Type.Elem().Kind() == reflect.String {
		if reflect.ValueOf(target).Elem().FieldByName(field.Name).IsNil() {
			return "", errValIsNil
		}
		strVal = reflect.ValueOf(target).Elem().FieldByName(field.Name).Elem().String()
	} else if field.Type.Kind() == reflect.String {
		strVal = reflect.ValueOf(target).Elem().FieldByName(field.Name).String()
	} else {
		return "", errors.New("the field must be of type string or *string")
	}

	return strVal, nil
}
