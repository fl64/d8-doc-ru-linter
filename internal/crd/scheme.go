package crd

import (
	"reflect"

	log "github.com/sirupsen/logrus"

	"d8-doc-ru-linter/internal/operation"
)

var keysToKeep = map[string]struct{}{
	"description": {},
}

func normalize(data interface{}, keysToKeep map[string]struct{}, parentPath string) interface{} {
	if data == nil {
		return nil
	}

	switch reflect.TypeOf(data).Kind() {
	case reflect.Map:
		m := make(map[interface{}]interface{})
		for k, v := range data.(map[interface{}]interface{}) {
			keyName := k.(string)
			parentPath := parentPath + "/" + keyName

			if _, ok := keysToKeep[keyName]; ok {
				if v == nil || reflect.TypeOf(v).Kind() == reflect.String {
					m[k] = v
					continue
				}
			}

			if v == nil {
				continue
			}

			res := normalize(v, keysToKeep, parentPath)
			if res != nil {
				m[k] = res
			} else {
				log.Debugf("%s\n", parentPath)
			}
		}

		if len(m) > 0 {
			return m
		}

	case reflect.Slice:
		s := make([]interface{}, 0)
		for index, v := range data.([]interface{}) {
			res := normalize(v, keysToKeep, parentPath)
			if res != nil {
				s = append(s, res)
			} else {
				log.Debugf("%s/%d\n", parentPath, index)
			}
		}
		if len(s) > 0 {
			return s
		}
	}

	return nil
}

type Schema struct {
	OpenAPIV3Schema interface{} `yaml:"openAPIV3Schema"`
}

func (s *Schema) Normalize() *Schema {
	return &Schema{
		OpenAPIV3Schema: normalize(s.OpenAPIV3Schema, keysToKeep, ""),
	}
}

func (s *Schema) CompareWith(dstScheme Schema, root string) (*Schema, operation.OperationsList) {
	operations := operation.NewOperationsList()

	var Compare func(interface{}, interface{}, string) interface{}
	Compare = func(origin, destination interface{}, parentPath string) interface{} {
		typeOfOrigin := reflect.TypeOf(origin).Kind()
		if destination == nil || reflect.TypeOf(destination).Kind() != typeOfOrigin {
			operations.Add(parentPath, operation.Replace)
			return origin
		}

		switch typeOfOrigin {
		case reflect.Map:
			{
				m := make(map[interface{}]interface{})
				originMap := origin.(map[interface{}]interface{})
				destinationMap := destination.(map[interface{}]interface{})
				for originKey, originValue := range originMap {
					parentPath := parentPath + "/" + originKey.(string)

					if _, ok := destinationMap[originKey.(string)]; !ok {
						operations.Add(parentPath, operation.Add)
						m[originKey.(string)] = originValue
						continue
					}
					m[originKey.(string)] = Compare(originValue, destinationMap[originKey.(string)], parentPath)
				}

				for destinationKey := range destinationMap {
					if _, ok := originMap[destinationKey.(string)]; !ok {
						operations.Add(parentPath+"/"+destinationKey.(string), operation.Delete)
					}
				}

				return m
			}
		default:
			return destination
		}
	}

	return &Schema{
		OpenAPIV3Schema: Compare(s.OpenAPIV3Schema, dstScheme.OpenAPIV3Schema, root),
	}, operations

}
