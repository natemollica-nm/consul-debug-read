#!/bin/bash

set -e

# Function to convert a number from one base to another
convert_base() {
  local input_number="$1"
  local from_base="$2"
  local to_base="$3"
  local integar_part fractional_part
  # Split the input into integer and fractional parts
  integer_part=$(echo "$input_number" | awk -F'.' '{print $1}')
  fractional_part=$(echo "$input_number" | awk -F'.' '{print $2}')

  # Convert the integer part from the current base to decimal
  local decimal_integer=0
  local power=0
  while [ "$integer_part" -gt 0 ]; do
    digit=$((integer_part % 10))
    decimal_integer=$((decimal_integer + digit * (from_base**power)))
    integer_part=$((integer_part / 10))
    power=$((power + 1))
  done

#  # Convert the decimal part from the current base to decimal
#  local decimal_fraction=0
#  local fractional_length="${#fractional_part}"
#  for ((i = 0; i < fractional_length; i++)); do
#    digit="${fractional_part:$i:1}"
#    decimal_fraction=$( echo $(( decimal_fraction + $((digit / $(( from_base**((i + 1)) )) )) )) | bc)
#  done
#  echo "$fractional_part converted to base 10: $decimal_fraction"

  # Convert the integer and fractional parts to the target base
  local converted_integer=""
  while [ "$decimal_integer" -gt 0 ]; do
    remainder=$((decimal_integer % to_base))
    converted_integer="$remainder$converted_integer"
    decimal_integer=$((decimal_integer / to_base))
  done

  local converted_fraction=""
  local fractional_length="$(echo -n "$fractional_part" | wc -c)"
  for ((i = 0; i < 6; i++)); do
    decimal_fraction=$(( ((decimal_fraction)) * to_base ))
    digit="$(echo $decimal_fraction | awk -F'.' '{print $1}')"
    converted_fraction="$converted_fraction$digit"
    decimal_fraction="${decimal_fraction#*.}"
  done
   echo "$fractional_part converted to base 10: $converted_fraction"
  # Combine the integer and fractional parts
  if [ -z "$converted_integer" ]; then
    converted_integer=0
  fi

  result="$converted_integer"
  if [ -n "$converted_fraction" ]; then
    result="$result.$converted_fraction"
  fi

  echo "$1 (base $2) is $result (base $3)"
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
