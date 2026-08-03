#!/bin/sh

set -e

script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)

user_id=${1:-${IKNOW_USER_ID:-}}
if [ -z "$user_id" ]; then
    echo "IKNOW_USER_ID is required" >&2
    exit 1
fi

data_dir=${IKNOW_DATA_DIR:-${script_dir}}
db_path=${IKNOW_DB_PATH:-${data_dir}/iknow.sqlite3}
output_dir=${data_dir}/images

mkdir -p "$output_dir"
cd "$data_dir"

"${script_dir}/iknow-tools" --db "$db_path" load "$user_id"


last_day_of_last_year=$(date -d "$(date +%Y-01-01) -1 day" +%Y-%m-%d)
last_day_of_this_year=$(date -d "$(date +%Y-01-01) +1 year -1 day" +%Y-%m-%d)
last_day_of_last_month=$(date -d "$(date +%Y-%m-01) -1 day" +%Y-%m-%d)
last_day_of_this_month=$(date -d "$(date +%Y-%m-01) +1 month -1 day" +%Y-%m-%d)
this_year=$(date +%Y)
this_month=$(date +%Y-%m)

mkdir -p "${output_dir}/${this_year}"

"${script_dir}/iknow-tools" --db "$db_path" plot -predict-completed "$last_day_of_last_year" "$last_day_of_this_year"
mv plot.png "${output_dir}/${this_year}/${this_year}.png"

"${script_dir}/iknow-tools" --db "$db_path" plot -predict-completed "$last_day_of_last_month" "$last_day_of_this_month"
mv plot.png "${output_dir}/${this_year}/${this_month}.png"

"${script_dir}/iknow-tools" --db "$db_path" diff --markdown "$last_day_of_last_year" "$last_day_of_this_year" > "${output_dir}/${this_year}/${this_year}.txt"
"${script_dir}/iknow-tools" --db "$db_path" diff --markdown "$last_day_of_last_month" "$last_day_of_this_month" > "${output_dir}/${this_year}/${this_month}.txt"
