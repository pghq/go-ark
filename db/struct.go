package db

import (
	"reflect"
	"strings"

	"github.com/pghq/go-tea"
)

// Map Convert a struct (w. optional tags) to a map using reflection
// variation of: https://play.golang.org/p/2Qi3thFf--
// meant to be used for data persistence.
func Map(v interface{}, transient ...interface{}) (map[string]interface{}, error) {
	if m, ok := v.(map[string]interface{}); ok {
		return m, nil
	}

	if m, ok := v.(*map[string]interface{}); ok {
		return *m, nil
	}

	rv := reflect.Indirect(reflect.ValueOf(v))
	for {
		if rv.Kind() != reflect.Ptr {
			break
		}
		rv = reflect.Indirect(rv)
	}

	if rv.Kind() != reflect.Struct {
		return nil, tea.NewErrorf("item of type %T is not a struct", v)
	}

	item := make(map[string]interface{})
	t := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		sf := t.Field(i)
		key := sf.Tag.Get("db")
		if key == "" {
			key = sf.Name
		}

		if key == "-" || len(transient) == 0 && strings.HasSuffix(key, ",transient") {
			continue
		}

		item[strings.Split(key, ",")[0]] = rv.Field(i).Interface()
	}

	return item, nil
}

// Copy Copies src value to destination
func Copy(src, dst interface{}) error {
	dv := reflect.Indirect(reflect.ValueOf(dst))
	if !dv.CanSet() {
		return tea.NewError("bad destination")
	}

	sv := reflect.Indirect(reflect.ValueOf(src))
	if dv.Type() != sv.Type() {
		return tea.NewError("type mismatch")
	}

	dv.Set(sv)
	return nil
}

// KeyName is a helper to determine a good default for struct persistence keys
func KeyName(v interface{}) string {
	if v == nil {
		return ""
	}

	rv := reflect.Indirect(reflect.ValueOf(v))
	for {
		if rv.Kind() != reflect.Ptr {
			break
		}
		rv = reflect.Indirect(rv)
	}

	switch rv.Kind() {
	case reflect.String:
		return ""
	case reflect.Struct:
		return ToSnakeCase(rv.Type().Name()) + "_id"
	default:
		return "id"
	}
}
