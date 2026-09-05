package repo

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"be/pkg/postgres"
)

type columnMapOpts struct {
	PK           string
	SkipEmptyPK  bool
	SkipZeroTime bool
	SkipNilPtr   bool
	SkipPK       bool
	SkipCreated  bool
}

func columnMap(entity any, opts columnMapOpts) (map[string]any, error) {
	value := derefStruct(entity)
	if !value.IsValid() {
		return nil, fmt.Errorf("column map: not a struct")
	}

	out := make(map[string]any)
	typ := value.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !field.IsExported() {
			continue
		}
		col := dbColumn(field)
		if col == "" {
			continue
		}
		if opts.SkipPK && col == opts.PK {
			continue
		}
		if opts.SkipCreated && col == "created_at" {
			continue
		}

		fv := value.Field(i)
		if opts.SkipEmptyPK && col == opts.PK && isEmptyPK(fv) {
			continue
		}
		if opts.SkipZeroTime && isZeroTime(fv) {
			continue
		}

		mapped, skip, err := encodeField(fv, opts.SkipNilPtr)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", col, err)
		}
		if skip {
			continue
		}
		out[postgres.QuoteIdent(col)] = mapped
	}
	return out, nil
}

func updateMap(entity any, pk string) (string, map[string]any, error) {
	value := derefStruct(entity)
	if !value.IsValid() {
		return "", nil, fmt.Errorf("update map: not a struct")
	}

	var pkValue string
	typ := value.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if dbColumn(field) == pk {
			pkValue = fmt.Sprint(value.Field(i).Interface())
			break
		}
	}

	cols, err := columnMap(entity, columnMapOpts{
		PK:           pk,
		SkipPK:       true,
		SkipCreated:  true,
		SkipNilPtr:   false,
		SkipZeroTime: false,
	})
	if err != nil {
		return "", nil, err
	}
	cols[postgres.QuoteIdent("updated_at")] = nowUTC()
	return pkValue, cols, nil
}

func dbColumn(field reflect.StructField) string {
	tag := field.Tag.Get("db")
	if tag == "" || tag == "-" {
		return ""
	}
	return strings.TrimSpace(strings.Split(tag, ",")[0])
}

func derefStruct(entity any) reflect.Value {
	value := reflect.ValueOf(entity)
	for value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return reflect.Value{}
		}
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return reflect.Value{}
	}
	return value
}

func isEmptyPK(fv reflect.Value) bool {
	switch fv.Kind() {
	case reflect.String:
		return strings.TrimSpace(fv.String()) == ""
	default:
		return fv.IsZero()
	}
}

func isZeroTime(fv reflect.Value) bool {
	if !fv.IsValid() {
		return false
	}
	t, ok := unwrapTime(fv)
	return ok && t.IsZero()
}

func unwrapTime(fv reflect.Value) (time.Time, bool) {
	if fv.Kind() == reflect.Pointer {
		if fv.IsNil() {
			return time.Time{}, false
		}
		fv = fv.Elem()
	}
	if fv.Type() == reflect.TypeOf(time.Time{}) {
		return fv.Interface().(time.Time), true
	}
	return time.Time{}, false
}

func encodeField(fv reflect.Value, skipNilPtr bool) (any, bool, error) {
	if !fv.IsValid() {
		return nil, true, nil
	}

	if fv.Kind() == reflect.Pointer {
		if fv.IsNil() {
			if skipNilPtr {
				return nil, true, nil
			}
			return nil, false, nil
		}
		return encodeField(fv.Elem(), skipNilPtr)
	}

	if fv.CanInterface() {
		if valuer, ok := fv.Interface().(driver.Valuer); ok {
			v, err := valuer.Value()
			return v, false, err
		}
	}

	switch fv.Kind() {
	case reflect.Slice, reflect.Array:
		if fv.Type().Elem().Kind() == reflect.Uint8 {
			if fv.Len() == 0 {
				return []byte(nil), false, nil
			}
			return fv.Bytes(), false, nil
		}
		if fv.IsNil() {
			if skipNilPtr {
				return emptyJSON(fv.Type()), false, nil
			}
			return emptyJSON(fv.Type()), false, nil
		}
		raw, err := json.Marshal(fv.Interface())
		if err != nil {
			return nil, false, err
		}
		return string(raw), false, nil
	case reflect.Map:
		if fv.IsNil() {
			return "{}", false, nil
		}
		raw, err := json.Marshal(fv.Interface())
		if err != nil {
			return nil, false, err
		}
		return string(raw), false, nil
	case reflect.String:
		return fv.String(), false, nil
	default:
		return fv.Interface(), false, nil
	}
}

func emptyJSON(t reflect.Type) string {
	if t.Kind() == reflect.Map {
		return "{}"
	}
	return "[]"
}

func nowUTC() time.Time {
	return time.Now().UTC()
}

func columnSelectExprs[T any]() []string {
	typ := reflect.TypeOf((*T)(nil)).Elem()
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil
	}

	out := make([]string, 0, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !field.IsExported() {
			continue
		}
		col := dbColumn(field)
		if col == "" {
			continue
		}
		quoted := postgres.QuoteIdent(col)
		ft := field.Type
		if ft.Kind() == reflect.Pointer {
			out = append(out, quoted)
			continue
		}
		switch {
		case ft.Kind() == reflect.String:
			out = append(out, fmt.Sprintf("COALESCE(%s::text, '') AS %s", quoted, quoted))
		case ft == reflect.TypeOf(time.Time{}):
			out = append(out, fmt.Sprintf("COALESCE(%s, TIMESTAMPTZ '1970-01-01') AS %s", quoted, quoted))
		case ft.Kind() == reflect.Bool:
			out = append(out, fmt.Sprintf("COALESCE(%s, FALSE) AS %s", quoted, quoted))
		case isIntKind(ft.Kind()):
			out = append(out, fmt.Sprintf("COALESCE(%s, 0) AS %s", quoted, quoted))
		default:
			out = append(out, quoted)
		}
	}
	return out
}

func isIntKind(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	default:
		return false
	}
}

func setStringColumn(entity any, column, value string) {
	v := derefStruct(entity)
	if !v.IsValid() || !v.CanSet() {
		// pointer to struct: derefStruct returns elem which is settable if original was pointer
		rv := reflect.ValueOf(entity)
		if rv.Kind() == reflect.Pointer && !rv.IsNil() {
			v = rv.Elem()
		}
	}
	if !v.IsValid() || v.Kind() != reflect.Struct {
		return
	}
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		if dbColumn(t.Field(i)) != column {
			continue
		}
		fv := v.Field(i)
		if fv.Kind() == reflect.String && fv.CanSet() {
			fv.SetString(value)
		}
		return
	}
}
