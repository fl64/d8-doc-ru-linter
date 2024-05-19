#! /usr/bin/env bash

go build -o d8-doc-ru-linter ../cmd/d8-doc-ru-linter

exclude="(values|config-values|cloud_discovery_data|cluster_configuration|openapi-case-tests)"

yamls=$(find /home/user/go/src -type f \( -name "*.yaml" -o -name "*.yml" \)  \( -path "*/crds/*" -o -path "*/openapi/*" \))

for yaml in $yamls; do
    # get file name from path
    yaml_file=$(basename $yaml)
    # get path wihout file name
    yaml_path=$(dirname $yaml)

    # get file name without extension
    yaml_file_no_ext="${yaml_file%.*}"

    # if file name is in exclude list, skip it
    if [[ $yaml_file_no_ext =~ $exclude ]]; then continue; fi

    doc_ru_yaml="$yaml_path/doc-ru-$yaml_file"

    if [[ ! -f "$doc_ru_yaml" ]]; then
      continue
    fi

    escaped_doc_ru_yaml=$(echo "${doc_ru_yaml#*/}" | tr '/' '-')

    changes=$(./d8-doc-ru-linter -s $yaml -d $doc_ru_yaml -n $escaped_doc_ru_yaml 2>&1)
    changes_count=$(echo $changes | jq -r .count)
    if [[ "$changes_count" -gt 0 ]]; then
      echo "======================================================================================="
      echo "src: $yaml"
      echo "dst: $doc_ru_yaml"
      echo "$changes" | jq . -c
      #dyff  bw -i -b $doc_ru_yaml $escaped_doc_ru_yaml
    fi
done
rm *.yaml
