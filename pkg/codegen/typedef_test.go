package codegen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeDefinition_GetErrorResponse(t *testing.T) {
	t.Run("single property without pointer", func(t *testing.T) {
		typ := TypeDefinition{
			Name: "ResError",
			Schema: GoSchema{
				Properties: []Property{
					{
						GoName:        "ID",
						JsonFieldName: "id",
					},
					{
						GoName:        "Details",
						JsonFieldName: "details",
						Schema: GoSchema{
							GoType: "string",
						},
					},
				},
			},
		}
		res := typ.GetErrorResponse(map[string]string{"ResError": "details"}, "e")
		expected := `res0 := e.Details
return res0`
		assert.Equal(t, expected, res)
	})

	t.Run("single property with pointer", func(t *testing.T) {
		typ := TypeDefinition{
			Name: "ResError",
			Schema: GoSchema{
				Properties: []Property{
					{
						GoName:        "ID",
						JsonFieldName: "id",
					},
					{
						GoName:        "Details",
						JsonFieldName: "details",
						Schema: GoSchema{
							GoType: "string",
						},
						Constraints: Constraints{
							Nullable: true,
						},
					},
				},
			},
		}
		res := typ.GetErrorResponse(map[string]string{"ResError": "details"}, "e")
		expected := `res0 := e.Details
if res0 == nil { return "unknown error" }
res1 := *res0
return res1`
		assert.Equal(t, expected, res)
	})

	t.Run("property with name error", func(t *testing.T) {
		typ := TypeDefinition{
			Name: "ResError",
			Schema: GoSchema{
				Properties: []Property{
					{
						GoName:        "ErrorData",
						JsonFieldName: "error",
						Schema: GoSchema{
							Properties: []Property{
								{
									GoName:        "Message",
									JsonFieldName: "message",
									Schema: GoSchema{
										GoType: "string",
									},
								},
							},
						},
					},
				},
			},
		}
		res := typ.GetErrorResponse(map[string]string{"ResError": "error.message"}, "e")
		expected := `res0 := e.ErrorData
res1 := res0.Message
return res1`
		assert.Equal(t, expected, res)
	})

	t.Run("nested property without pointer", func(t *testing.T) {
		typ := TypeDefinition{
			Name: "ResError",
			Schema: GoSchema{
				Properties: []Property{
					{
						GoName:        "Data",
						JsonFieldName: "data",
						Schema: GoSchema{
							Properties: []Property{
								{
									GoName:        "Details",
									JsonFieldName: "details",
									Schema: GoSchema{
										Properties: []Property{
											{
												GoName:        "Message",
												JsonFieldName: "message",
												Schema:        GoSchema{GoType: "string"},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}
		res := typ.GetErrorResponse(map[string]string{"ResError": "data.details.message"}, "e")
		expected := `res0 := e.Data
res1 := res0.Details
res2 := res1.Message
return res2`
		assert.Equal(t, expected, res)
	})

	t.Run("nested property with pointer", func(t *testing.T) {
		typ := TypeDefinition{
			Name: "ResError",
			Schema: GoSchema{
				Properties: []Property{
					{
						GoName:        "Data",
						JsonFieldName: "data",
						Constraints:   Constraints{Nullable: true},
						Schema: GoSchema{
							Properties: []Property{
								{
									GoName:        "Details",
									JsonFieldName: "details",
									Constraints:   Constraints{Nullable: true},
									Schema: GoSchema{
										Properties: []Property{
											{
												GoName:        "Message",
												JsonFieldName: "message",
												Constraints:   Constraints{Nullable: true},
												Schema: GoSchema{
													GoType: "string",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}
		res := typ.GetErrorResponse(map[string]string{"ResError": "data.details.message"}, "e")
		expected := `res0 := e.Data
if res0 == nil { return "unknown error" }
res1 := *res0
res2 := res1.Details
if res2 == nil { return "unknown error" }
res3 := *res2
res4 := res3.Message
if res4 == nil { return "unknown error" }
res5 := *res4
return res5`
		assert.Equal(t, expected, res)
	})
}
