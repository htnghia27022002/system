// Package apidoc builds the OpenAPI 3.1 document from Go types. Each route is
// described once, next to where it is mounted, by an Op that names its request
// and response DTOs; schemas are derived from those structs by reflection
// (json / form / binding tags), so the spec cannot drift from the code. A test
// compares the ops with the live router so an undocumented route fails CI.
package apidoc

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"
)

// Op documents one HTTP route.
type Op struct {
	Method string
	// Path is the gin path, e.g. "/api/admin/users/:id".
	Path    string
	Tag     string
	Summary string
	// Auth: the route needs a Bearer access token.
	Auth bool
	// Permission is the RBAC key enforced (documented as x-permission).
	Permission string
	// Query is a struct whose `form` fields become query parameters.
	Query any
	// Body is the JSON request body type.
	Body any
	// Response is the success body type; nil means no body.
	Response any
	// Status is the success status (default 200, or 204 when Response is nil).
	Status int
	// Codes lists notable error codes (documented as x-error-codes).
	Codes []string
	// Raw marks a non-JSON response (files, redirects).
	Raw string
}

// Info is the document header.
type Info struct {
	Title, Version, Description string
}

// Build returns the OpenAPI document as a JSON-ready map.
func Build(info Info, ops []Op) map[string]any {
	g := newGenerator()
	paths := map[string]map[string]any{}

	sorted := append([]Op(nil), ops...)
	sort.SliceStable(sorted, func(i, j int) bool {
		if sorted[i].Path != sorted[j].Path {
			return sorted[i].Path < sorted[j].Path
		}
		return sorted[i].Method < sorted[j].Method
	})

	for _, op := range sorted {
		path, params := openAPIPath(op.Path)
		operation := map[string]any{
			"operationId": operationID(op.Method, op.Path),
			"summary":     op.Summary,
			"tags":        []string{op.Tag},
		}
		parameters := []any{}
		for _, p := range params {
			parameters = append(parameters, map[string]any{
				"name": p, "in": "path", "required": true, "schema": map[string]any{"type": "string"},
			})
		}
		if op.Query != nil {
			parameters = append(parameters, g.queryParams(reflect.TypeOf(op.Query))...)
		}
		if len(parameters) > 0 {
			operation["parameters"] = parameters
		}
		if op.Body != nil {
			operation["requestBody"] = map[string]any{
				"required": true,
				"content":  map[string]any{"application/json": map[string]any{"schema": g.schema(reflect.TypeOf(op.Body))}},
			}
		}
		if op.Auth {
			operation["security"] = []any{map[string]any{"bearerAuth": []string{}}}
		}
		if op.Permission != "" {
			operation["x-permission"] = op.Permission
		}
		if len(op.Codes) > 0 {
			operation["x-error-codes"] = op.Codes
		}

		status := op.Status
		if status == 0 {
			status = http.StatusOK
			if op.Response == nil && op.Raw == "" {
				status = http.StatusNoContent
			}
		}
		success := map[string]any{"description": http.StatusText(status)}
		switch {
		case op.Raw != "":
			success["content"] = map[string]any{op.Raw: map[string]any{"schema": map[string]any{"type": "string", "format": "binary"}}}
		case op.Response != nil:
			success["content"] = map[string]any{"application/json": map[string]any{"schema": g.schema(reflect.TypeOf(op.Response))}}
		}
		operation["responses"] = map[string]any{
			fmt.Sprint(status): success,
			"default": map[string]any{
				"description": "Error",
				"content":     map[string]any{"application/json": map[string]any{"schema": map[string]any{"$ref": "#/components/schemas/Error"}}},
			},
		}

		if paths[path] == nil {
			paths[path] = map[string]any{}
		}
		method := strings.ToLower(op.Method)
		if method == "any" {
			for _, m := range []string{"get", "post", "put", "patch", "delete"} {
				copyOp := cloneMap(operation)
				copyOp["operationId"] = operationID(m, op.Path)
				paths[path][m] = copyOp
			}
			continue
		}
		paths[path][method] = operation
	}

	g.schemas["Error"] = map[string]any{
		"type":     "object",
		"required": []string{"error", "code"},
		"properties": map[string]any{
			"error":  map[string]any{"type": "string", "description": "Human-readable message (English)."},
			"code":   map[string]any{"type": "string", "description": "Stable machine-readable code, e.g. not_found, validation_failed, email_not_verified."},
			"fields": map[string]any{"type": "object", "additionalProperties": map[string]any{"type": "string"}, "description": "validation_failed only: field → rule."},
		},
	}

	return map[string]any{
		"openapi": "3.1.0",
		"info":    map[string]any{"title": info.Title, "version": info.Version, "description": info.Description},
		"paths":   paths,
		"components": map[string]any{
			"schemas": g.schemas,
			"securitySchemes": map[string]any{
				"bearerAuth": map[string]any{"type": "http", "scheme": "bearer", "bearerFormat": "JWT"},
			},
		},
	}
}

var pathParam = regexp.MustCompile(`[:*]([A-Za-z0-9_]+)`)

func openAPIPath(ginPath string) (string, []string) {
	var params []string
	path := pathParam.ReplaceAllStringFunc(ginPath, func(m string) string {
		params = append(params, m[1:])
		return "{" + m[1:] + "}"
	})
	return path, params
}

func operationID(method, path string) string {
	parts := []string{strings.ToLower(method)}
	for _, seg := range strings.Split(strings.TrimPrefix(path, "/api/"), "/") {
		seg = strings.TrimLeft(seg, ":*")
		for _, word := range strings.FieldsFunc(seg, func(r rune) bool { return r == '-' || r == '_' }) {
			parts = append(parts, strings.ToUpper(word[:1])+word[1:])
		}
	}
	return strings.Join(parts, "")
}

func cloneMap(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

// ---- schema generation ----

type generator struct {
	schemas map[string]any
	names   map[reflect.Type]string
	taken   map[string]reflect.Type
}

func newGenerator() *generator {
	return &generator{schemas: map[string]any{}, names: map[reflect.Type]string{}, taken: map[string]reflect.Type{}}
}

var (
	timeType = reflect.TypeOf(time.Time{})
	anyType  = reflect.TypeOf((*any)(nil)).Elem()
	rawJSON  = reflect.TypeOf([]byte(nil))
)

func (g *generator) schema(t reflect.Type) map[string]any {
	nullable := false
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
		nullable = true
	}
	s := g.schemaNonNull(t)
	if nullable {
		if ref, ok := s["$ref"]; ok {
			return map[string]any{"oneOf": []any{map[string]any{"$ref": ref}, map[string]any{"type": "null"}}}
		}
		if typ, ok := s["type"].(string); ok {
			out := cloneMap(s)
			out["type"] = []string{typ, "null"}
			return out
		}
	}
	return s
}

func (g *generator) schemaNonNull(t reflect.Type) map[string]any {
	switch {
	case t == timeType:
		return map[string]any{"type": "string", "format": "date-time"}
	case t == anyType || t.Kind() == reflect.Interface:
		return map[string]any{}
	case t.String() == "json.RawMessage":
		return map[string]any{}
	case t == rawJSON:
		return map[string]any{"type": "string", "format": "byte"}
	}
	switch t.Kind() {
	case reflect.String:
		return map[string]any{"type": "string"}
	case reflect.Bool:
		return map[string]any{"type": "boolean"}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return map[string]any{"type": "integer"}
	case reflect.Float32, reflect.Float64:
		return map[string]any{"type": "number"}
	case reflect.Slice, reflect.Array:
		return map[string]any{"type": "array", "items": g.schema(t.Elem())}
	case reflect.Map:
		return map[string]any{"type": "object", "additionalProperties": g.schema(t.Elem())}
	case reflect.Struct:
		return map[string]any{"$ref": "#/components/schemas/" + g.define(t)}
	default:
		return map[string]any{}
	}
}

// define registers a named struct schema and returns its component name.
func (g *generator) define(t reflect.Type) string {
	if name, ok := g.names[t]; ok {
		return name
	}
	name := schemaName(t)
	if named, ok := reflect.Zero(t).Interface().(interface{ OpenAPIName() string }); ok {
		name = named.OpenAPIName()
	}
	if other, clash := g.taken[name]; clash && other != t {
		pkg := t.PkgPath()
		name = exportName(pkg[strings.LastIndex(pkg, "/")+1:]) + name
	}
	g.names[t] = name
	g.taken[name] = t
	g.schemas[name] = map[string]any{} // placeholder for recursive types

	props := map[string]any{}
	var required []string
	hasBinding := structHasBinding(t)
	g.collectFields(t, props, &required, hasBinding)
	s := map[string]any{"type": "object", "properties": props}
	if len(required) > 0 {
		sort.Strings(required)
		s["required"] = required
	}
	g.schemas[name] = s
	return name
}

func (g *generator) collectFields(t reflect.Type, props map[string]any, required *[]string, hasBinding bool) {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		jsonTag := f.Tag.Get("json")
		name, opts, _ := strings.Cut(jsonTag, ",")
		if f.Anonymous && name == "" {
			ft := f.Type
			if ft.Kind() == reflect.Pointer {
				ft = ft.Elem()
			}
			if ft.Kind() == reflect.Struct {
				g.collectFields(ft, props, required, hasBinding)
				continue
			}
		}
		if name == "-" {
			continue
		}
		if name == "" {
			name = f.Name
		}
		s := g.schema(f.Type)
		applyBinding(s, f.Tag.Get("binding"))
		props[name] = s

		binding := f.Tag.Get("binding")
		isRequired := strings.Contains(","+binding+",", ",required,")
		if !hasBinding && !strings.Contains(opts, "omitempty") && f.Type.Kind() != reflect.Pointer {
			// Response DTOs: every non-optional field is always present.
			isRequired = true
		}
		if isRequired {
			*required = append(*required, name)
		}
	}
}

func (g *generator) queryParams(t reflect.Type) []any {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	var out []any
	var walk func(reflect.Type)
	walk = func(t reflect.Type) {
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if f.Anonymous && f.Type.Kind() == reflect.Struct {
				walk(f.Type)
				continue
			}
			name := strings.SplitN(f.Tag.Get("form"), ",", 2)[0]
			if name == "" || name == "-" {
				continue
			}
			s := g.schema(f.Type)
			binding := f.Tag.Get("binding")
			applyBinding(s, binding)
			out = append(out, map[string]any{
				"name": name, "in": "query", "required": strings.Contains(","+binding+",", ",required,"), "schema": s,
			})
		}
	}
	walk(t)
	return out
}

// applyBinding maps go-playground/validator rules onto JSON Schema keywords.
func applyBinding(s map[string]any, binding string) {
	if binding == "" {
		return
	}
	typ, _ := s["type"].(string)
	for _, rule := range strings.Split(binding, ",") {
		key, val, _ := strings.Cut(rule, "=")
		switch key {
		case "email":
			s["format"] = "email"
		case "uuid":
			s["format"] = "uuid"
		case "url":
			s["format"] = "uri"
		case "oneof":
			s["enum"] = strings.Fields(val)
		case "min", "max":
			n := 0
			fmt.Sscan(val, &n)
			switch typ {
			case "string":
				s[map[string]string{"min": "minLength", "max": "maxLength"}[key]] = n
			case "integer", "number":
				s[map[string]string{"min": "minimum", "max": "maximum"}[key]] = n
			case "array":
				s[map[string]string{"min": "minItems", "max": "maxItems"}[key]] = n
			}
		}
	}
}

func structHasBinding(t reflect.Type) bool {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Tag.Get("binding") != "" {
			return true
		}
		if f.Anonymous && f.Type.Kind() == reflect.Struct && structHasBinding(f.Type) {
			return true
		}
	}
	return false
}

var genericArgs = regexp.MustCompile(`\[(.*)\]$`)

// schemaName turns "Page[be/internal/dto/user.UserResponse]" into "PageUserResponse".
func schemaName(t reflect.Type) string {
	name := t.Name()
	if m := genericArgs.FindStringSubmatch(name); m != nil {
		base := strings.TrimSuffix(name, m[0])
		var args []string
		for _, a := range strings.Split(m[1], ",") {
			a = a[strings.LastIndex(a, ".")+1:]
			args = append(args, exportName(a))
		}
		return base + strings.Join(args, "")
	}
	if name == "" {
		return "Anonymous"
	}
	return exportName(name)
}

func exportName(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
