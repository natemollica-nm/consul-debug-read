#!/bin/bash

# Function to convert a number from one base to another
convert_base() {
  local input_number="$1"
  local from_base="$2"
  local to_base="$3"

  # Convert the input number from the current base to base 10 (decimal)
  local decimal_number=0
  local power=0

  while [ "$input_number" -gt 0 ]; do
    digit=$((input_number % 10))
    decimal_number=$((decimal_number + digit * (from_base**power)))
    input_number=$((input_number / 10))
    power=$((power + 1))
  done

  # Convert the decimal number to the target base
  local converted_number=""
  while [ "$decimal_number" -gt 0 ]; do
    remainder=$((decimal_number % to_base))
    converted_number="$remainder$converted_number"
    printf "$decimal_number/$to_base = "
    decimal_number=$((decimal_number / to_base))
    printf "$decimal_number, remainder $remainder\n"
  done

  echo "$1 (base $2) is $converted_number (base $3)"
}

# Check for the correct number of arguments
if [ "$#" -ne 3 ]; then
  echo "Usage: $0 <input_number> <from_base> <to_base>"
  exit 1
fi

# Input number, current base, and target base
input_number="$1"
from_base="$2"
to_base="$3"

# Validate the input bases
if ! [[ "$from_base" =~ ^[0-9]+$ ]] || ! [[ "$to_base" =~ ^[0-9]+$ ]] || [ "$from_base" -lt 2 ] || [ "$to_base" -lt 2 ]; then
  echo "Base values must be integers greater than or equal to 2."
  exit 1
fi

# Call the conversion function
convert_base "$input_number" "$from_base" "$to_base"
