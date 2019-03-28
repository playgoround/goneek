package parser

import "github.com/json-iterator/go"

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Parser interface {
	Parse() ([]byte, error)
}
