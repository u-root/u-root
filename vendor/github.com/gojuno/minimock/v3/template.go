package minimock

import _ "embed"

// HeaderTemplate is used to generate package clause and go:generate instruction
//
//go:embed internal/templates/header.tmpl
var HeaderTemplate string

// BodyTemplate is used to generate mock body
//
//go:embed internal/templates/body.tmpl
var BodyTemplate string
