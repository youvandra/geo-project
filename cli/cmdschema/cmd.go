package cmdschema

import (
	"fmt"
	"os"
	"strings"
)

func Run(args []string) {
	if len(args) < 2 {
		printUsage()
		os.Exit(1)
	}

	schemaType := SchemaType(args[0])
	valid := false
	for _, t := range SupportedTypes() {
		if t == schemaType {
			valid = true
			break
		}
	}
	if !valid {
		fmt.Fprintf(os.Stderr, "unsupported schema type: %s\n", schemaType)
		fmt.Fprintf(os.Stderr, "supported types: %s\n", strings.Join(func() []string {
			types := []string{}
			for _, t := range SupportedTypes() {
				types = append(types, string(t))
			}
			return types
		}(), ", "))
		os.Exit(1)
	}

	outputFormat := "json"
	fieldArgs := []string{}

	for _, a := range args[1:] {
		switch a {
		case "--html":
			outputFormat = "html"
		default:
			fieldArgs = append(fieldArgs, a)
		}
	}

	fields := parseFields(fieldArgs)
	if len(fields) == 0 {
		fmt.Fprintln(os.Stderr, "error: no fields provided")
		fmt.Fprintln(os.Stderr, "usage: geo schema <type> <key=value>... [--html]")
		os.Exit(1)
	}

	sb := NewSchemaBuilder(schemaType, fields)

	var output string
	var err error

	switch outputFormat {
	case "html":
		output, err = sb.BuildHTML()
	default:
		output, err = sb.BuildJSON()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(output)
}

func parseFields(args []string) []SchemaField {
	fields := []SchemaField{}
	for _, a := range args {
		if !strings.Contains(a, "=") {
			continue
		}
		parts := strings.SplitN(a, "=", 2)
		name := parts[0]
		value := parts[1]

		isArray := false
		if strings.HasSuffix(name, "[]") {
			isArray = true
			name = strings.TrimSuffix(name, "[]")
		}

		fields = append(fields, SchemaField{
			Name:  name,
			Value: value,
			Array: isArray,
		})
	}
	return fields
}

func printUsage() {
	fmt.Println("Usage: geo schema <type> <key=value>... [--html]")
	fmt.Println()
	fmt.Println("Types:")
	fmt.Println("  FAQPage       FAQ schema with questions and answers")
	fmt.Println("  Article       Blog/article schema")
	fmt.Println("  HowTo         Tutorial/steps schema")
	fmt.Println("  Product       Product schema")
	fmt.Println("  LocalBusiness Local business schema")
	fmt.Println("  Person        Person/author schema")
	fmt.Println("  Organization  Organization/brand schema")
	fmt.Println("  Dataset       Dataset schema")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  geo schema Article headline=\"My Article\" author=\"John Doe\"")
	fmt.Println("  geo schema FAQPage mainEntity=\"[{\\\"@type\\\":\\\"Question\\\",...}]\"")
	fmt.Println("  geo schema Person name=\"John Doe\" jobTitle=\"Writer\" --html")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --html  Output as <script> tag (default: raw JSON)")
}
