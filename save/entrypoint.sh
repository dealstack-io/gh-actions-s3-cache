#!/bin/bash

$GITHUB_ACTION_PATH/build/$(echo "$OS" | tr "[:upper:]" "[:lower:]") save