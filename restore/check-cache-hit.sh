#!/bin/bash

FILE=$1.zip

if [ -f "$FILE" ]; then
  echo "true"
else
  echo "false"
fi