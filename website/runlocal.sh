#!/bin/bash

[ x"$U-root" != x ] && cd $U-root/web
mkdocs serve && exit 0 
echo "#### we need mkdocs: use pip install mkdocs"


