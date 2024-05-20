# d8-doc-ru-linter

## Build

```bash
task local-build
```

## Test

```bash
task test
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
./d8-doc-ru-linter -s ./testdata/src.yaml -n new.yaml
```

Merge two CRDs

- Keys that are present in source and missing in destination will be added
- Keys that are not in source and are in destination will be deleted
- The values of all other keys in destination will remain unchanged

```bash
./d8-doc-ru-linter -s ./testdata/src.yaml -d ./testdata/dst.yaml -n new.yaml
```

Just check if some diffs exist

```bash
./d8-doc-ru-linter -s ./testdata/src.yaml -d ./testdata/dst.yaml -n /dev/null | jq .
```

Output:

```json
{
  "count": 4,
  "operations": [
    {
      "path": "/spec/versions/v2",
      "op": "add"
    },
    {
      "path": "/spec/versions/v1/schema/openAPIV3Schema/properties/p1-object/properties/p1-1-object/properties/only-in-src",
      "op": "add"
    },
    {
      "path": "/spec/versions/v1/schema/openAPIV3Schema/properties/p1-object/properties/p1-1-object/properties/only-in-dst",
      "op": "delete"
    },
    {
      "path": "/spec/versions/v0",
      "op": "delete"
    }
  ]
}
```

Fail with exit code 33 if some diffs exist

```bash
./d8-doc-ru-linter -s ./testdata/src.yaml -d ./testdata/dst.yaml -n /dev/null --fail
```

Using docker

```bash
docker run -it -v ${PWD}:/tmp index.docker.io/fl64/d8-doc-ru-linter:v0.0.1-dev0 /d8-doc-ru-linter -s /tmp/testdata/src.yaml -d /tmp/testdata/dst.yaml -n /dev/null

# {"count":4,"operations":[{"path":"/spec/versions/v2","op":"add"},{"path":"/spec/versions/v1/schema/openAPIV3Schema/properties/p1-object/properties/p1-1-object/properties/only-in-src","op":"add"},{"path":"/spec/versions/v1/schema/openAPIV3Schema/properties/p1-object/properties/p1-1-object/properties/only-in-dst","op":"delete"},{"path":"/spec/versions/v0","op":"delete"}]}
```
