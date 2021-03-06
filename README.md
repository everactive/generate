# generate

Generates Go (golang) Structs and Validation code from JSON schema.

Note: A hash value is added to each generated file to make the id's unqiue when they are put in the resulting  output files to avoid conflicts.  There is a possibility of hashes colliding.

# Requirements

* Go 1.8+

# Usage

Install

```console
$ go get -u github.com/everactive/generate/...
```

or

Build

```console
$ make
```

Run

```console
$ schema-generate exampleschema.json
```

To specify an output file:

```console
$ schema-generate exampleschema.json -o example.go
```

To make primitive values that are optional to be pointers, use e.g.,
```console
schema-generate -pointerPrimitives -o sensor.go sensor_schema.json
```

# Example

This schema

```json
{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "title": "Example",
  "id": "http://example.com/exampleschema.json",
  "type": "object",
  "description": "An example JSON Schema",
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
  },
  "definitions": {
    "address": {
      "id": "address",
      "type": "object",
      "description": "Address",
      "properties": {
        "street": {
          "type": "string",
          "description": "Address 1",
          "maxLength": 40
        },
        "houseNumber": {
          "type": "integer",
          "description": "House Number"
        }
      }
    },
    "status": {
      "type": "object",
      "properties": {
        "favouritecat": {
          "enum": [
            "A",
            "B",
            "C"
          ],
          "type": "string",
          "description": "The favourite cat.",
          "maxLength": 1
        }
      }
    }
  }
}
```

generates

```go
package main

type Address struct {
  HouseNumber int `json:"houseNumber,omitempty"`
  Street string `json:"street,omitempty"`
}

type Example struct {
  Address *Address `json:"address,omitempty"`
  Name string `json:"name,omitempty"`
  Status *Status `json:"status,omitempty"`
}

type Status struct {
  Favouritecat string `json:"favouritecat,omitempty"`
}
```

See the [test/](./test/) directory for more examples.
