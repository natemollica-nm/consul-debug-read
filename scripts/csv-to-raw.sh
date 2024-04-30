#!/bin/bash

# Path to your CSV file
csv_file="$1"

## Ensure non-null csv file param value is passed.
if [ -z "$csv_file" ]; then
  echo "Usage: $(basename "$0") <csv_file>"
fi

# Check if the file exists
if [[ ! -f "$csv_file" ]]; then
    echo "File $csv_file does not exist."
    exit 1
fi

# Read the CSV file line by line, skipping the header
# Process the CSV file
printf '%s\n' "Date     Host    Service     Thread     Module     Level     Message"
awk -F, 'NR > 1 { # Skip the header line
    # Clean up fields: remove leading/trailing quotes and extra quotes inside fields
    for(i = 1; i <= NF; i++) {
        gsub(/^"|"$/, "", $i);
        gsub(/""/, "\"", $i);
    }

    # Format level to uppercase
    level = toupper($6);

    # Concatenate message parts if the message was split due to internal commas
    message = $7;
    for(i = 8; i <= NF; i++) {
        message = message ", " $i;
    }

    # Print the reformatted log entry
    # Assuming fields are: Date, Host, Service, Thread, Module, Level, Message
    print $1 " " $2":"$3 " [" level "] " $5 " : " message
}' "$csv_file"