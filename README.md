d8-doc-ru-linter
================

Build:

```bash
go build ./cmd/d8-doc-ru-linter
```

Usage:

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

Example:

```bash
./d8-doc-ru-linter -s /go/src/github.com/deckhouse/candi/openapi/node_group.yaml -d /go/src/github.com/deckhouse/candi/openapi/doc-ru-node_group.yaml -n /dev/null --fail | jq .

# Output:
# {
#   "count": 7,
#   "operations": [
#     {
#       "path": "/spec/versions/v1alpha1/schema/openAPIV3Schema/properties/status",
#       "op": "add"
#     },
#     {
#       "path": "/spec/versions/v1alpha1/schema/openAPIV3Schema/properties/spec/properties/cri/properties/type",
#       "op": "add"
#     },
#     {
#       "path": "/spec/versions/v1alpha1/schema/openAPIV3Schema/properties/spec/properties/disruptions/properties/rollingUpdate",
#       "op": "add"
#     },
#     {
#       "path": "/spec/versions/v1alpha2/schema/openAPIV3Schema/properties/status",
#       "op": "add"
#     },
#     {
#       "path": "/spec/versions/v1alpha2/schema/openAPIV3Schema/properties/spec/properties/disruptions/properties/rollingUpdate",
#       "op": "add"
#     },
#     {
#       "path": "/spec/versions/v1/schema/openAPIV3Schema/properties/spec/properties/staticInstances/properties/labelSelector/properties/matchExpressions/items/description",
#       "op": "add"
#     },
#     {
#       "path": "/spec/versions/v1/schema/openAPIV3Schema/properties/status",
#       "op": "add"
#     }
#   ]
# }
```
