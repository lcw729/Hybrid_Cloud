#!/bin/bash

for dir in */; do mkdir -- "$dir/git"; done

for dir in */; do mv "$dir/.git" "$dir/git"; done

git rm -r --cached