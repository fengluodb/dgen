package main

import (
	"dgen/codegen"
	"dgen/config"
	"flag"
	"log"
)

var filename string
var language string
var outputDir string
var encodeType string

func init() {
	flag.StringVar(&filename, "f", "", "filename")
	flag.StringVar(&language, "l", "", "the language to generate")
	flag.StringVar(&outputDir, "o", ".", "the dir of output file")
	flag.StringVar(&encodeType, "e", "", "the type of encoding")
}

func main() {
	flag.Parse()
	if filename == "" {
		log.Println("filename cannot be empty")
		return
	}
	if language == "" {
		log.Println("language cannot be empty")
		return
	}

	config := &config.CodegenConfig{
		Filename:   filename,
		OutputDir:  outputDir,
		EncodeType: encodeType,
	}

	err := codegen.CodegenMap[language](config)
	if err != nil {
		log.Println("failed to generate code, error: ", err)
	}
}
