#!/bin/bash
prj_name="backend"
context_name=$(basename "$PWD")
if [[ $context_name != "$prj_name" ]]; then
    echo "Wrong pwd - it should be $prj_name"
    exit 1
fi
go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/upsert ./ent/schema
