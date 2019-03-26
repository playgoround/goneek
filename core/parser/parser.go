package parser

type Parser interface {
	Parse() ([]byte, error)
}
