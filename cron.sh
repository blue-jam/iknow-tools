#!/bin/sh

set -e
# Move to the directory where the script is located
cd "$(dirname "$0")"

./iknow-tools load "$1"
