# jstransform CHANGELOG

# 0.1.0
Tim Kuhlman - Add support for creating time.Time fields based on the format type in the JSON schema.

# 0.0.9
Kent Lee - Expand input/output paths and add go generate examples

# 0.0.8
Tim Kuhlman - Generate a parent struct to embed in others for JSON schemas with allOf refs

# 0.0.7
John Lin - handle JSON files with no oneOfTypes

# 0.0.6
Tim Kuhlman - Converting a nil to a string should result in a nil return with no error

# 0.0.5
Tim Kuhlman - Implement go struct generation from JSON schema

# 0.0.4
Dennis Nguyen - default nil booleans

# 0.0.3
Tim Kuhlman - Expand the walk function to walk by instance or raw JSON

# 0.0.2
Tim Kuhlman - split out transform code from jsonschema

# 0.0.1
Tim Kuhlman - initial creation of jstransform repo
