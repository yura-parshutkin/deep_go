package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

const delimiter = '='

func Serialize(input any) string {
	rt := reflect.TypeOf(input)
	if rt.Kind() != reflect.Struct {
		panic("expected struct type, actual: " + rt.Kind().String())
	}
	rv := reflect.ValueOf(input)
	bf := strings.Builder{}
	fieldsCount := rt.NumField()

	for i := 0; i < fieldsCount; i++ {
		prop, ok := rt.Field(i).Tag.Lookup("properties")
		if !ok {
			continue
		}
		key, omitempty := parseProp(prop)
		val := fmt.Sprintf("%v", rv.Field(i))
		// skip empty values if omitempty is set
		if omitempty && val == "" {
			continue
		}
		bf.WriteString(key)
		bf.WriteByte(delimiter)
		bf.WriteString(val)
		// do not use \n for last line
		if i != fieldsCount-1 {
			bf.WriteByte('\n')
		}
	}
	return bf.String()
}

func parseProp(prop string) (val string, omitempty bool) {
	parts := strings.Split(prop, ",")
	if len(parts) < 2 {
		return prop, false
	}
	if parts[0] == "omitempty" {
		return parts[1], true
	}
	if parts[1] == "omitempty" {
		return parts[0], true
	}
	panic(fmt.Sprintf("unknown property %q", prop))
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}
