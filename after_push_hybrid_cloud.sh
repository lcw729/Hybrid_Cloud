#!/bin/bash


for dir in */; do mv "./${dir}git/.git" "$dir"; done