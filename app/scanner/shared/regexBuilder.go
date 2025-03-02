package scannerShared

import (
	"fmt"
	"reflect"
	"strings"
)

const stringRegex string = "[\"'`]"

func GenerateRegex(s interface{}) string {
	// assumes all fields are optional

	val := reflect.ValueOf(s)
	t := reflect.TypeOf(s)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return val.String()
	}

	parts := []string{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		nameRe := field.Tag.Get("regex")
		if nameRe == "" {
			nameRe = field.Name
		}

		valueRe := ""
		if field.Type.Kind() == reflect.String || (field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.String) {
			valueRe = fmt.Sprintf(`%s.*%s`, stringRegex, stringRegex)
		} else if field.Type.Kind() == reflect.Int || (field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Int) {
			valueRe = `\d+`
		} else if field.Type.Kind() == reflect.Bool || (field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Bool) {
			valueRe = `(true|false)`
		} else if field.Type.Kind() == reflect.Slice || (field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Slice) {
			valueRe = fmt.Sprintf(`\[%s?.*%s?\]`, stringRegex, stringRegex)
		} else if field.Type.Kind() == reflect.Interface || (field.Type.Kind() == reflect.Pointer && field.Type.Elem().Kind() == reflect.Struct) {
			valueRe = GenerateRegex(reflect.New(field.Type).Elem().Interface())
		}

		parts = append(parts, "("+strings.Join([]string{fmt.Sprintf(`%s?%s%s?`, stringRegex, nameRe, stringRegex), ":", valueRe}, `\s*`)+")")
	}

	return fmt.Sprintf(`{\s*((%s)*\s*,?)*\s*}`, strings.Join(parts, `|`))
}
