package config

import (
	"log"
	"testing"
)

func TestInit(t *testing.T) {
	Init("../bin/config.yaml")
	log.Printf("Conifg :: %+v", Instance)
}
