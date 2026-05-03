#!/bin/sh

set -e
# Move to the directory where the script is located
cd "$(dirname "$0")"

./iknow-tools load "$1"

mkdir -p images

last_day_of_last_year=$(date -d "$(date +%Y-01-01) -1 day" +%Y-%m-%d)
last_day_of_this_year=$(date -d "$(date +%Y-01-01) +1 year -1 day" +%Y-%m-%d)
last_day_of_last_month=$(date -d "$(date +%Y-%m-01) -1 day" +%Y-%m-%d)
last_day_of_this_month=$(date -d "$(date +%Y-%m-01) +1 month -1 day" +%Y-%m-%d)
this_year=$(date +%Y)
this_month=$(date +%Y-%m)

mkdir -p "images/${this_year}"

./iknow-tools plot -predict-completed "$last_day_of_last_year" "$last_day_of_this_year"
mv plot.png "images/${this_year}/${this_year}.png"

./iknow-tools plot -predict-completed "$last_day_of_last_month" "$last_day_of_this_month"
mv plot.png "images/${this_year}/${this_month}.png"

./iknow-tools diff --markdown "$last_day_of_last_year" "$last_day_of_this_year" > "images/${this_year}/${this_year}.md"
./iknow-tools diff --markdown "$last_day_of_last_month" "$last_day_of_this_month" > "images/${this_year}/${this_month}.md"
