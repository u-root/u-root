#!/bin/bash

# Must run in website directory.
mkdocs build -d ../docs/ && exit 0
echo "#### we need mkdocs: use pip install mkdocs"

