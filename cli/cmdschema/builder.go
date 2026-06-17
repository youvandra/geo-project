package cmdschema

import (
	"encoding/json"
	"fmt"
	"time"
)

type SchemaType string

const (
	SchemaFAQPage       SchemaType = "FAQPage"
	SchemaArticle       SchemaType = "Article"
	SchemaHowTo         SchemaType = "HowTo"
	SchemaProduct       SchemaType = "Product"
	SchemaLocalBusiness SchemaType = "LocalBusiness"
	SchemaPerson        SchemaType = "Person"
	SchemaOrganization  SchemaType = "Organization"
	SchemaDataset       SchemaType = "Dataset"
)

type SchemaField struct {
	Name  string
	Value string
	Array bool
}

type SchemaBuilder struct {
	Type   SchemaType
	Fields map[string]interface{}
}

func NewSchemaBuilder(schemaType SchemaType, fields []SchemaField) *SchemaBuilder {
	sb := &SchemaBuilder{
		Type:   schemaType,
		Fields: map[string]interface{}{},
	}

	for _, f := range fields {
		if f.Array {
			sb.Fields[f.Name] = []string{f.Value}
		} else {
			sb.Fields[f.Name] = f.Value
		}
	}

	if _, ok := sb.Fields["@context"]; !ok {
		sb.Fields["@context"] = "https://schema.org"
	}
	sb.Fields["@type"] = string(schemaType)

	return sb
}

func (sb *SchemaBuilder) Set(name string, value interface{}) {
	sb.Fields[name] = value
}

func (sb *SchemaBuilder) GetRequiredFields() []string {
	switch sb.Type {
	case SchemaFAQPage:
		return []string{"mainEntity"}
	case SchemaArticle:
		return []string{"headline", "author"}
	case SchemaHowTo:
		return []string{"name", "step"}
	case SchemaProduct:
		return []string{"name"}
	case SchemaLocalBusiness:
		return []string{"name", "address"}
	case SchemaPerson:
		return []string{"name"}
	case SchemaOrganization:
		return []string{"name"}
	case SchemaDataset:
		return []string{"name", "description"}
	default:
		return []string{"name"}
	}
}

func (sb *SchemaBuilder) Validate() error {
	for _, f := range sb.GetRequiredFields() {
		if _, ok := sb.Fields[f]; !ok {
			return fmt.Errorf("missing required field: %s (for %s)", f, sb.Type)
		}
	}
	return nil
}

func (sb *SchemaBuilder) BuildJSON() (string, error) {
	if err := sb.Validate(); err != nil {
		return "", err
	}

	addDateModified(sb)

	data, err := json.MarshalIndent(sb.Fields, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (sb *SchemaBuilder) BuildHTML() (string, error) {
	jsonStr, err := sb.BuildJSON()
	if err != nil {
		return "", err
	}

	html := fmt.Sprintf(
		`<script type="application/ld+json">
%s
</script>`, jsonStr)

	return html, nil
}

func addDateModified(sb *SchemaBuilder) {
	switch sb.Type {
	case SchemaArticle, SchemaDataset:
		if _, ok := sb.Fields["dateModified"]; !ok {
			sb.Fields["dateModified"] = time.Now().Format(time.RFC3339)
		}
	}
}

func SupportedTypes() []SchemaType {
	return []SchemaType{
		SchemaFAQPage,
		SchemaArticle,
		SchemaHowTo,
		SchemaProduct,
		SchemaLocalBusiness,
		SchemaPerson,
		SchemaOrganization,
		SchemaDataset,
	}
}
