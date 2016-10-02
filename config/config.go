package config

import (
	"github.com/hashicorp/hcl"
	"strings"
)

type Record struct {
	Values []*ValueRecord          `hcl:"value"`
	Builds map[string]*BuildRecord `hcl:"build"`
}

type ValueRecord struct {
	Name          string `hcl:",key"`
	Literal       string
	DecodedFields []string `hcl:",decodedFields"`
}

type BuildRecord struct {
	Steps []*StepRecord `hcl:"step,expand"`
}

type StepRecord struct {
	Type          string `hcl:",key"`
	Url           string
	Image         string
	Cmd           string
	Dir           string
	DecodedFields []string `hcl:",decodedFields"`
}

func hasField(fields []string, s string) bool {
	for _, f := range fields {
		if strings.EqualFold(s, f) {
			return true
		}
	}
	return false
}

func (vr *ValueRecord) HasField(s string) bool {
	return hasField(vr.DecodedFields, s)
}

func (sr *StepRecord) HasField(s string) bool {
	return hasField(sr.DecodedFields, s)
}

func LoadBytes(b []byte) (*Record, error) {
	r := &Record{}
	err := hcl.Unmarshal(b, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}
