package codegen

import (
	"fmt"
	"os"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel"
)

func loadDocumentFromFile(filepath string) (libopenapi.Document, error) {
	contents, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return loadDocumentFromContents(contents)
}

func loadDocumentFromContents(contents []byte) (libopenapi.Document, error) {
	docConfig := &datamodel.DocumentConfiguration{
		SkipCircularReferenceCheck: true,
	}
	doc, err := libopenapi.NewDocumentWithConfiguration(contents, docConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating document: %w", err)
	}
	return doc, nil
}
