package gowalker

import (
	"reflect"
	"strings"
)

// TagString - parsed struct field tag
type TagString struct {
	Value string

	IsDefaultValue bool
	DefaultValue   string
}

// TagStringParse - parse struct field tag
func TagStringParse(field reflect.StructField, tagKey string) (out TagString) {
	tag := field.Tag.Get(tagKey)
	tag, opts := head(tag, ",")
	var opt string
	for len(opts) > 0 {
		opt, opts = head(opts, ",")

		k, v := head(opt, "=")
		switch k {
		case "default":
			out.IsDefaultValue = true
			out.DefaultValue = v
		}
	}

	if tag == "" {
		tag = field.Name
	}

	out.Value = tag
	return
}

func head(s string, sep string) (head string, tail string) {
	idx := strings.Index(s, sep)
	if idx < 0 {
		return s, ""
	}
	return s[:idx], s[idx+len(sep):]
}
