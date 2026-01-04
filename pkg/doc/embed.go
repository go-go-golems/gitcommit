package doc

import "embed"

// FS embeds gitcommit help/documentation markdown sections.
//
// Files are loaded into the Glazed help system during CLI initialization.
// See cmd/gitcommit/main.go.
//
//go:embed topics/*.md
var FS embed.FS
