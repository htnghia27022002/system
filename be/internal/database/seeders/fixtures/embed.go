package fixtures

import _ "embed"

// VietnamDivisionsJSON is the P1 VN country + 34 provinces + commune subset.
//
//go:embed vn_divisions.json
var VietnamDivisionsJSON []byte
