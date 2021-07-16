package generate

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
)

func getOrderedFieldNames(m map[string]Field) []string {
	keys := make([]string, len(m))
	idx := 0
	for k := range m {
		keys[idx] = k
		idx++
	}
	sort.Strings(keys)
	return keys
}

func getOrderedStructNames(m map[string]Struct) []string {
	keys := make([]string, len(m))
	idx := 0
	for k := range m {
		keys[idx] = k
		idx++
	}
	sort.Strings(keys)
	return keys
}

// Output generates code and writes to w.
func Output(w io.Writer, g *Generator, pkg string, pointerPrimitives bool) {
	structs := g.Structs
	aliases := g.Aliases

	fmt.Fprintln(w, "// Code generated by schema-generate. DO NOT EDIT.")
	fmt.Fprintln(w)

	packageName := cleanPackageName(pkg)

	fmt.Fprintf(w, "// Package %v contains the structs and types as defined by this schema.\n", packageName)
	fmt.Fprintf(w, "package %v\n", packageName)

	// write all the code into a buffer, compiler functions will return list of imports
	// write list of imports into main output stream, followed by the code
	codeBuf := new(bytes.Buffer)
	imports := make(map[string]bool)

	if len(imports) > 0 {
		fmt.Fprintf(w, "\nimport (\n")
		for k := range imports {
			fmt.Fprintf(w, "    \"%s\"\n", k)
		}
		fmt.Fprintf(w, ")\n")
	}

	if len(g.schemas) > 0 {
		fmt.Fprintln(w, "\nconst (")
		for index, schema := range g.schemas {
			schemaIdParts := strings.Split(schema.ID(), "/")
			schemaIdPartLength := len(schemaIdParts)
			if schemaIdPartLength > 1 {
				fmt.Fprintf(w, "    SchemaID%d", index)
				fmt.Fprintf(w, " = \"%s\"\n", schemaIdParts[schemaIdPartLength-1])
			}
		}
		fmt.Fprintf(w, "\n)\n")
	}

	for _, k := range getOrderedFieldNames(aliases) {
		a := aliases[k]

		fmt.Fprintln(w, "")
		fmt.Fprintf(w, "// %s\n", a.Name)
		fmt.Fprintf(w, "type %s %s\n", a.Name, a.Type)
	}

	for _, k := range getOrderedStructNames(structs) {
		s := structs[k]

		fmt.Fprintln(w, "")
		outputNameAndDescriptionComment(s.Name, s.Description, w)
		fmt.Fprintf(w, "type %s struct {\n", s.Name)

		for _, fieldKey := range getOrderedFieldNames(s.Fields) {
			f := s.Fields[fieldKey]

			// Only apply omitempty if the field is not required.
			omitempty := ",omitempty"
			typePrefix := ""
			if f.Required {
				omitempty = ""
			}
			if !f.Required && pointerPrimitives {
				if f.Type == "bool" || f.Type == "int" || f.Type == "string" || f.Type == "float64" {
					typePrefix = "*"
				}
			}

			if f.Description != "" {
				outputFieldDescriptionComment(f.Description, w)
			}

			fmt.Fprintf(w, "  %s %s%s `json:\"%s%s\"`\n", f.Name, typePrefix, f.Type, f.JSONName, omitempty)
		}

		fmt.Fprintln(w, "}")
	}

	// write code after structs for clarity
	w.Write(codeBuf.Bytes())
}

func outputNameAndDescriptionComment(name, description string, w io.Writer) {
	if strings.Index(description, "\n") == -1 {
		fmt.Fprintf(w, "// %s %s\n", name, description)
		return
	}

	dl := strings.Split(description, "\n")
	fmt.Fprintf(w, "// %s %s\n", name, strings.Join(dl, "\n// "))
}

func outputFieldDescriptionComment(description string, w io.Writer) {
	if strings.Index(description, "\n") == -1 {
		fmt.Fprintf(w, "\n  // %s\n", description)
		return
	}

	dl := strings.Split(description, "\n")
	fmt.Fprintf(w, "\n  // %s\n", strings.Join(dl, "\n  // "))
}

func cleanPackageName(pkg string) string {
	pkg = strings.Replace(pkg, ".", "", -1)
	pkg = strings.Replace(pkg, "_", "", -1)
	pkg = strings.Replace(pkg, "-", "", -1)
	return pkg
}
