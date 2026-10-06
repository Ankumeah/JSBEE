package netlifyblob

import _ "embed"

//go:embed scripts/blob.mjs
var helperScriptBytes []byte

//go:embed scripts/package.json
var helperPackageJSONBytes []byte
