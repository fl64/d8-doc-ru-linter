package crd

import (
	"d8-doc-ru-linter/internal/operation"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type CRD struct {
	Spec struct {
		Versions []struct {
			Name   string `yaml:"name"`
			Schema Schema `yaml:"schema"`
		} `yaml:"versions"`
	} `yaml:"spec"`
}

func (c *CRD) Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("can't load crd: %v", err)
	}

	if err = yaml.Unmarshal(data, &c); err != nil {
		return fmt.Errorf("can't unmarshal crd: %v", err)
	}

	for i, scheme := range c.Spec.Versions {
		c.Spec.Versions[i].Schema = *scheme.Schema.Normalize()
	}

	return nil
}

func (c *CRD) Marshal() ([]byte, error) {
	result, err := yaml.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("can't marshal crd: %v", err)
	}
	return result, nil
}

func (c CRD) CompareWith(dst CRD) (CRD, operation.OperationsList) {
	operations := operation.NewOperationsList()

	for originIndex, originVersions := range c.Spec.Versions {
		originName := originVersions.Name
		realDestinationIndex := -1
		for destinationIndex, destinationVersions := range dst.Spec.Versions {
			if destinationVersions.Name == originName {
				realDestinationIndex = destinationIndex
			}
		}

		if realDestinationIndex == -1 {
			continue
		}

		root := fmt.Sprintf("/spec/versions/%s/schema/openAPIV3Schema", originVersions.Name)
		res, ops := originVersions.Schema.CompareWith(dst.Spec.Versions[realDestinationIndex].Schema, root)
		c.Spec.Versions[originIndex].Schema.OpenAPIV3Schema = res.OpenAPIV3Schema
		operations.Operations = append(operations.Operations, ops.Operations...)
	}
	return c, operations
}
