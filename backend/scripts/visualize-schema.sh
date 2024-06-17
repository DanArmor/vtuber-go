#!/bin/bash
sudo atlas schema inspect \
  -u "ent://ent/schema" \
  --dev-url "docker://postgres/15/test?search_path=public" \
  -w