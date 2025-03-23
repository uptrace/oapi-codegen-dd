package codegen

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// RequestBodyDefinition describes a request body.
// Name is the name of the body.
// Required is whether the body is required.
// GoSchema is the GoSchema object describing the body.
// NameTag is the tag used to generate the type name,
// i.e. JSON, in which case we will produce "JSONBody".
// ContentType is the content type of the body.
// Default is whether this is the default body type.
// Encoding is the encoding options for formdata.
type RequestBodyDefinition struct {
	Name        string
	Required    bool
	Schema      GoSchema
	NameTag     string
	ContentType string
	Default     bool
	Encoding    map[string]RequestBodyEncoding
}

// TypeDef returns the Go type definition for a request body
func (r RequestBodyDefinition) TypeDef(opID string) TypeDefinition {
	return TypeDefinition{
		Name:   fmt.Sprintf("%s%sRequestBody", opID, r.NameTag),
		Schema: r.Schema,
	}
}

// CustomType returns whether the body is a custom inline type, or pre-defined. This is
// poorly named, but it's here for compatibility reasons post-refactoring
// TODO: clean up the templates code, it can be simpler.
func (r RequestBodyDefinition) CustomType() bool {
	return r.Schema.RefType == ""
}

// Suffix is needed when we're generating multiple functions which relate to request bodies,
// this generates the suffix. Such as Operation DoFoo would be suffixed with
// DoFooWithXMLBody.
func (r RequestBodyDefinition) Suffix() string {
	// The default response is never suffixed.
	if r.Default {
		return ""
	}
	return "With" + r.NameTag + "Body"
}

// IsSupportedByClient returns true if we support this content type for client. Otherwise only generic method will ge generated
func (r RequestBodyDefinition) IsSupportedByClient() bool {
	return r.IsJSON() || r.NameTag == "Formdata" || r.NameTag == "Text"
}

// IsJSON returns whether this is a JSON media type, for instance:
// - application/json
// - application/vnd.api+json
// - application/*+json
func (r RequestBodyDefinition) IsJSON() bool {
	return isMediaTypeJson(r.ContentType)
}

// IsSupported returns true if we support this content type for server. Otherwise io.Reader will be generated
func (r RequestBodyDefinition) IsSupported() bool {
	return r.NameTag != ""
}

// IsFixedContentType returns true if content type has fixed content type, i.e. contains no "*" symbol
func (r RequestBodyDefinition) IsFixedContentType() bool {
	return !strings.Contains(r.ContentType, "*")
}

func (r RequestBodyDefinition) IsOptional() bool {
	return !r.Schema.Constraints.Required
}

type RequestBodyEncoding struct {
	ContentType string
	Style       string
	Explode     *bool
}

// createBodyDefinition turns the OpenAPI body definitions into a list of our body definitions
// which will be used for code generation.
func createBodyDefinition(operationID string, bodyOrRef *openapi3.RequestBodyRef) (*RequestBodyDefinition, *TypeDefinition, error) {
	if bodyOrRef == nil {
		return nil, nil, nil
	}

	td := TypeDefinition{}

	body := bodyOrRef.Value
	var targetContentType string
	for _, contentType := range sortedMapKeys(body.Content) {
		if contentType == "application/json" {
			targetContentType = contentType
			break
		}
		targetContentType = contentType
	}

	content := body.Content[targetContentType]
	var tag string
	var defaultBody bool
	required := body.Required

	switch {
	case targetContentType == "application/json":
		tag = "JSON"
		defaultBody = true
	case isMediaTypeJson(targetContentType):
		tag = mediaTypeToCamelCase(targetContentType)
	case strings.HasPrefix(targetContentType, "multipart/"):
		tag = "Multipart"
	case targetContentType == "application/x-www-form-urlencoded":
		tag = "Formdata"
	case targetContentType == "text/plain":
		tag = "Text"
	default:
		return nil, nil, nil
	}

	bodyTypeName := operationID + "Body"
	bodySchema, err := GenerateGoSchema(content.Schema, []string{bodyTypeName})
	if err != nil {
		return nil, nil, fmt.Errorf("error generating request body definition: %w", err)
	}

	// If the request has a body, but it's not a user defined
	// type under #/components, we'll define a type for it, so
	// that we have an easy-to-use type for marshaling.
	if bodySchema.RefType == "" {
		if targetContentType == "application/x-www-form-urlencoded" {
			// Apply the appropriate structure tag if the request
			// schema was defined under the operations' section.
			for i := range bodySchema.Properties {
				bodySchema.Properties[i].NeedsFormTag = true
			}

			// Regenerate the Golang struct adding the new form tag.
			fields := genFieldsFromProperties(bodySchema.Properties)
			bodySchema.GoType = bodySchema.createGoStruct(fields)
		}

		td = TypeDefinition{
			Name:         bodyTypeName,
			Schema:       bodySchema,
			SpecLocation: SpecLocationBody,
		}
		// The body schema now is a reference to a type
		bodySchema.RefType = bodyTypeName
	}

	bodySchema.Constraints.Required = required

	bd := &RequestBodyDefinition{
		Name:        bodyTypeName,
		Required:    body.Required,
		Schema:      bodySchema,
		NameTag:     tag,
		ContentType: targetContentType,
		Default:     defaultBody,
	}

	if len(content.Encoding) != 0 {
		bd.Encoding = make(map[string]RequestBodyEncoding)
		for k, v := range content.Encoding {
			encoding := RequestBodyEncoding{ContentType: v.ContentType, Style: v.Style, Explode: v.Explode}
			bd.Encoding[k] = encoding
		}
	}

	return bd, &td, nil
}
