package main

import (
	"github.com/kelseyhightower/envconfig"
	"log"
)

type Specification struct {
	Redis    string
	ListenOn string
	AuthKey  string
	CryptKey string
}

func ParseSpec() Specification {
	var spec Specification
	err := envconfig.Process("mustache", &spec)
	if err != nil {
		log.Fatalf("Parsing spec, %v", err)
	}
	return spec
}
