# d8-doc-ru-linter

## Build

```bash
go build ./cmd/d8-doc-ru-linter
```

## Usage

```
Usage: d8-doc-ru-linter [--debug] [--fail] --source SOURCE [--destination DESTINATION] [--new NEW]

Options:
  --debug                turn on debug mode
  --fail                 fail if there are diffs
  --source SOURCE, -s SOURCE
                         origin CRD
  --destination DESTINATION, -d DESTINATION
                         destination CRD
  --new NEW, -n NEW      file to save merged CRD
  --help, -h             display this help and exit

```

## Examples

Generate new normalized CRD

```bash
./d8-doc-ru-linter -s /go/src/github.com/deckhouse/candi/openapi/node_group.yaml -n /go/src/github.com/deckhouse/candi/openapi/doc-ru-node_group.yaml
```

Merge two CRDs

- Keys that are present in source and missing in destination will be added
- Keys that are not in source and are in destination will be deleted
- The values of all other keys in destination will remain unchanged

```bash
./d8-doc-ru-linter -s /go/src/github.com/deckhouse/candi/openapi/node_group.yaml -d /go/src/github.com/deckhouse/candi/openapi/doc-ru-node_group.yaml -n /go/src/github.com/deckhouse/candi/openapi/doc-ru-node_group.yaml
```

Just check and fail if some diffs exist

```bash
./d8-doc-ru-linter -s /go/src/github.com/deckhouse/candi/openapi/node_group.yaml -d /go/src/github.com/deckhouse/candi/openapi/doc-ru-node_group.yaml -n /dev/null --fail | jq .
```

Output:

```json
{
  "count": 1,
  "operations": [
    {
      "path": "/spec/versions/v1alpha1/schema/openAPIV3Schema/properties/status",
      "op": "add"
    }
  ]
}
```
