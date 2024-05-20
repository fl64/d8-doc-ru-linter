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

func (c *CRD) normalize() {
	for i, scheme := range c.Spec.Versions {
		c.Spec.Versions[i].Schema = *scheme.Schema.Normalize()
	}
}

func (c *CRD) Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("can't load crd: %v", err)
	}

	if err = yaml.Unmarshal(data, &c); err != nil {
		return fmt.Errorf("can't unmarshal crd: %v", err)
	}
	c.normalize()
	return nil
}

func (c *CRD) Marshal() ([]byte, error) {
	result, err := yaml.Marshal(c)
	if err != nil {
		return nil, fmt.Errorf("can't marshal crd: %v", err)
	}
	return result, nil
}

func (c *CRD) VersionIndex(version string) int {
	for index, scheme := range c.Spec.Versions {
		if scheme.Name == version {
			return index
		}
	}
	return -1
}

func (c CRD) CompareWith(dst CRD) (CRD, operation.OperationsList) {
	operations := operation.NewOperationsList()

	for originIndex, originVersions := range c.Spec.Versions {
		originName := originVersions.Name
		destinationIndex := dst.VersionIndex(originName)

		if destinationIndex == -1 {
			operations.Add("/spec/versions/"+originVersions.Name, operation.Add)
			continue
		}

		root := fmt.Sprintf("/spec/versions/%s/schema/openAPIV3Schema", originVersions.Name)
		res, ops := originVersions.Schema.CompareWith(dst.Spec.Versions[destinationIndex].Schema, root)
		c.Spec.Versions[originIndex].Schema.OpenAPIV3Schema = res.OpenAPIV3Schema
		operations.Operations = append(operations.Operations, ops.Operations...)
	}

	// add removed versions
	for _, destinationVersions := range dst.Spec.Versions {
		if c.VersionIndex(destinationVersions.Name) == -1 {
			operations.Add("/spec/versions/"+destinationVersions.Name, operation.Delete)
		}
	}

	return c, operations
}
