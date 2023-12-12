#!/bin/sh

set -u

replace() {
  if sed --version >/dev/null 2>&1; then
    # GNU
    sed -i -e "$1" schema.go
  else
    # BSD ('' needed or a backup file will be created)
    sed -i '' -e "$1" schema.go
  fi
}

# Delete zz_schema.go if it exists so that there are no compiler errors during generation.
rm -f zz_schema.go

# Replace "var c converter = &converterImpl{}" with "var c converter" to regenerate the schema.
# This is done so that there is no compiler error during generation.
replace "s/var c converter = \&converterImpl{}/var c converter/g"

# Generate zz_schema.go
go run github.com/jmattheis/goverter/cmd/goverter gen ./...

# After the file is generated we can add the reference back in.
replace "s/var c converter/var c converter = \&converterImpl{}/g"
