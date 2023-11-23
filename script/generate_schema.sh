#!/bin/sh

set -u

sed="sed -i \'\'"
if sed --version >/dev/null 2>&1; then # BSD sed exits with 1 on `sed --version`, GNU sed doesn't
  sed="sed -i" # GNU sed
fi

# Delete zz_schema.go if it exists so that there are no compiler errors during generation.
rm -f zz_schema.go

# Replace "var c converter = &converterImpl{}" with "var c converter" to regenerate the schema.
# This is done so that there is no compiler error during generation.
$sed -e "s/var c converter = \&converterImpl{}/var c converter/g" schema.go

# Generate zz_schema.go
go run github.com/jmattheis/goverter/cmd/goverter gen ./...

# After the file is generated we can add the reference back in.
$sed -e "s/var c converter/var c converter = \&converterImpl{}/g" schema.go
