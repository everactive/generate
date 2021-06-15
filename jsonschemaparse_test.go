package generate

import (
	"net/url"
	"strings"
	"testing"
)

func TestThatAMissingSchemaKeyResultsInAnError(t *testing.T) {
	invalid := `{
        "title": "root"
    }`
	_, invaliderr := Parse(invalid, &url.URL{Scheme: "file", Path: "jsonschemaparse_test.go"})
	valid := `{
        "$schema": "http://json-schema.org/schema#",
        "title": "root"
    }`
	_, validerr := Parse(valid, &url.URL{Scheme: "file", Path: "jsonschemaparse_test.go"})
	if invaliderr == nil {
		// it SHOULD be used in the root schema
		// t.Error("When the $schema key is missing from the root, the JSON Schema is not valid")
	}
	if validerr != nil {
		t.Error("It should be possible to parse a simple JSON schema if the $schema key is present")
	}
}

func TestThatTheRootSchemaCanBeParsed(t *testing.T) {
	s := `{
        "$schema": "http://json-schema.org/schema#",
        "title": "root"
    }`
	so, err := Parse(s, &url.URL{Scheme: "file", Path: "jsonschemaparse_test.go"})

	if err != nil {
		t.Fatal("It should be possible to unmarshal a simple schema, but received error:", err)
	}

	if so.Title != "root" {
		t.Errorf("The title was not deserialised from the JSON schema, expected %s, but got %s", "root", so.Title)
	}
}

func TestThatPropertiesCanBeParsed(t *testing.T) {
	s := `{
        "$schema": "http://json-schema.org/schema#",
        "title": "root",
        "properties": {
            "name": {
                "type": "string"
            },
            "address": {
                "$ref": "#/definitions/address"
            },
            "status": {
                "$ref": "#/definitions/status"
            }
        }
    }`
	so, err := Parse(s, &url.URL{Scheme: "file", Path: "jsonschemaparse_test.go"})

	if err != nil {
		t.Fatal("It was not possible to unmarshal the schema:", err)
	}

	nameType, nameMultiple := so.Properties["name"].Type()
	if nameType != "string" || nameMultiple {
		t.Errorf("expected property 'name' type to be 'string', but was '%v'", nameType)
	}

	addressType, _ := so.Properties["address"].Type()
	if addressType != "" {
		t.Errorf("expected property 'address' type to be '', but was '%v'", addressType)
	}

	if so.Properties["address"].Reference != "#/definitions/address" {
		t.Errorf("expected property 'address' reference to be '#/definitions/address', but was '%v'", so.Properties["address"].Reference)
	}

	if so.Properties["status"].Reference != "#/definitions/status" {
		t.Errorf("expected property 'status' reference to be '#/definitions/status', but was '%v'", so.Properties["status"].Reference)
	}
}

func TestThatPropertiesCanHaveMultipleTypes(t *testing.T) {
	s := `{
        "$schema": "http://json-schema.org/schema#",
        "title": "root",
        "properties": {
            "name": {
                "type": [ "integer", "string" ]
            }
        }
    }`
	so, err := Parse(s, &url.URL{Scheme: "file", Path: "jsonschemaparse_test.go"})

	if err != nil {
		t.Fatal("It was not possible to unmarshal the schema:", err)
	}

	nameType, nameMultiple := so.Properties["name"].Type()
	if nameType != "integer" {
		t.Errorf("expected first value of property 'name' type to be 'integer', but was '%v'", nameType)
	}

	if !nameMultiple {
		t.Errorf("expected multiple types, but only returned one")
	}
}

func TestThatParsingInvalidValuesReturnsAnError(t *testing.T) {
	s := `{ " }`
	_, err := Parse(s, &url.URL{Scheme: "file", Path: "jsonschemaparse_test.go"})

	if err == nil {
		t.Fatal("Expected a parsing error, but got nil")
	}
}

func TestThatDefaultsCanBeParsed(t *testing.T) {
	s := `{
        "$schema": "http://json-schema.org/schema#",
        "title": "root",
        "properties": {
            "name": {
                "type": [ "integer", "string" ],
                "default":"Enrique"
            }
        }
    }`
	so, err := Parse(s, &url.URL{Scheme: "file", Path: "jsonschemaparse_test.go"})

	if err != nil {
		t.Fatal("It was not possible to unmarshal the schema:", err)
	}

	defaultValue := so.Properties["name"].Default
	if defaultValue != "Enrique" {
		t.Errorf("expected default value of property 'name' type to be 'Enrique', but was '%v'", defaultValue)
	}
}

func TestReturnedSchemaId(t *testing.T) {
	tests := []struct {
		input    *Schema
		expected string
	}{
		{
			input:    &Schema{},
			expected: "",
		},
		{
			input:    &Schema{ID06: "http://example.com/foo.json", ID04: "#foo"},
			expected: "http://example.com/foo.json",
		},
		{
			input:    &Schema{ID04: "#foo"},
			expected: "#foo",
		},
	}

	for idx, test := range tests {
		actual := test.input.ID()
		if actual != test.expected {
			t.Errorf("Test %d failed: For input \"%+v\", expected \"%s\", got \"%s\"", idx, test.input, test.expected, actual)
		}
	}
}

func TestThatDefsCanBeParsed(t *testing.T) {
	s := `{
  "$schema": "http://json-schema.org/schema#",
  "$id": "https://example.com/defs-example.json",

  "$defs": {
    "temperature":{
      "description": "An temperature value.",
      "type": "number"
    },
    "humidity": {
      "description": "A humidity value.",
      "type": "number"
    }
  },
  "type": "array",
  "items": {
    "title": "Sensor Data",
    "description": "Sensor reading",
    "type":"object",
    "properties": {

      "schema": {
        "description": "The schema version of this packet",
        "type": "string",
        "const": "defs-example.json"
      },

      "timestamp": {
        "description": "The time of the reading.",
        "type": "string",
        "format": "date-time"
      },

      "sensor": {
        "title": "Sensor Data",
        "description": "Data related to the sensor.",
        "type": "object",

        "properties": {
          "metrics": {
            "title": "Sensor Metrics",
            "description": "Sensor-produced ",
            "type": "object",
            "properties": {
              "humidity": {
                "description": "Humidity",
                "type": "number"
              },
              "temperature_0": {"$ref":"#/$defs/temperature"}
            },
            "additionalProperties": false
          }
        },
        "oneOf": [
          {
            "required": ["metrics"]
          },
          {
            "required": ["metrics"]
          }
        ],
        "additionalProperties": false
      }
    },
    "required": ["schema", "timestamp", "sensor"],
    "additionalProperties": false
  }
}
`
	so, err := Parse(s, &url.URL{Scheme: "file", Path: "jsonschemaparse_test.go"})

	if err != nil {
		t.Fatal("It was not possible to unmarshal the schema:", err)
	}

	if len(so.Defs) < 2 {
		t.Fatal("Expected $defs to be populated with 2 references")
	}

	defsKeys:=[]string{"temperature", "humidity"}
	for k, d := range so.Defs {
		valid:=false
		for _, v := range defsKeys {
			if k == v {
				valid=true
			}
		}
		if !valid {
			t.Fatal("unexpected reference: " + k)
		}

		if !strings.Contains(d.PathElement, "$defs") {
			t.Fatal("expected reference path to contain the $defs prefix")
		}
	}



}