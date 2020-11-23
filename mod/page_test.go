package mod

import (
	"encoding/json"
	"log"
	"testing"
)

func TestPageInfo_Value(t *testing.T) {
	p := PageInfo{}
	data, err := json.MarshalIndent(&p, "", "\t")
	if err != nil {
		panic(err)
	}
	log.Printf("\n%v", string(data))
}
