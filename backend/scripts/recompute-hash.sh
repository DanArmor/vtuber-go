#!/bin/bash
prj_name="vtuber-go"
context_name=$(basename "$PWD")
if [[ $context_name != "$prj_name" ]]; then
    echo "Wrong pwd - it should be $prj_name"
    exit 1
fi
cd ent/migrate &&
atlas migrate hash