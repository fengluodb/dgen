package codegen

import (
	"dgen/codegen/gogen"
	"dgen/config"
)

var CodegenMap = map[string]func(config *config.CodegenConfig) error{
	"go": gogen.Gen,
}
