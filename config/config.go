package config

import (
	"github.com/hashicorp/hcl"
	"strings"
)

type Record struct {
	Builds map[string]*BuildRecord `hcl:"build"`
}

type BuildRecord struct {
	Steps []*StepRecord `hcl:"step,expand"`
}

type StepRecord struct {
	Type string  `hcl:",key"`
	Url string
	Image string
	DecodedFields []string `hcl:",decodedFields"`
}

func (sr *StepRecord) HasField(s string) bool {
	for _, f := range sr.DecodedFields {
		if strings.EqualFold(s, f) {
			return true
		}
	}
	return false
}

func LoadBytes(b []byte) (*Record, error) {
	r := &Record{}
	err := hcl.Unmarshal(b, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}
