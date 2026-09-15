package content

import _ "embed"

//go:embed catalog.json
var catalogJSON []byte

//go:embed lessons.json
var lessonsJSON []byte

//go:embed boards.json
var boardsJSON []byte

//go:embed materials.json
var materialsJSON []byte
