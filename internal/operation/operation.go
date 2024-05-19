package operation

import (
	"encoding/json"
	"fmt"
)

const (
	Add     = "add"
	Replace = "replace"
	Delete  = "delete"
)

type Operation struct {
	Path      string `json:"path"`
	Operation string `json:"op"`
}

type OperationsList struct {
	Count      int         `json:"count"`
	Operations []Operation `json:"operations"`
}

func (o *OperationsList) Len() int {
	return len(o.Operations)
}

func (o *OperationsList) Add(path, op string) {
	o.Operations = append(o.Operations, Operation{
		Path: path, Operation: op,
	})
}

func (o *OperationsList) MarshalJSONReport() ([]byte, error) {
	o.Count = len(o.Operations)
	result, err := json.Marshal(o)
	if err != nil {
		return nil, fmt.Errorf("can't marshal report: %v", err)
	}
	return result, nil
}

func NewOperationsList() OperationsList {
	return OperationsList{
		Operations: make([]Operation, 0),
	}
}
