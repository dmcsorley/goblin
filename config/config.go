package config

import (
	"bufio"
	"errors"
	"github.com/hashicorp/hcl"
	"os"
	"strings"
)

type Record struct {
	Values []*ValueRecord          `hcl:"value"`
	Builds map[string]*BuildRecord `hcl:"build"`
}

type ValueRecord struct {
	Name          string `hcl:",key"`
	Literal       string
	Env           string
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

func (vr *ValueRecord) value() string {
	if hasField(vr.DecodedFields, "literal") {
		return vr.Literal
	}

	if hasField(vr.DecodedFields, "env") {
		return os.Getenv(vr.Env)
	}

	return ""
}

func (sr *StepRecord) HasParameter(s string) bool {
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

func tokenize(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if len(data) == 0 {
		return 0, nil, nil
	}

	switch data[0] {
	case '$':
		return 1, []byte("$"), nil
	case '{':
		return 1, []byte("{"), nil
	case '}':
		return 1, []byte("}"), nil
	default:
		token = nil
	INNER:
		for _, b := range data {
			switch b {
			case '$', '{', '}':
				break INNER
			default:
				token = append(token, b)
			}
		}

		if len(token) < len(data) || atEOF {
			return len(token), token, nil
		} else {
			return 0, nil, nil
		}
	}
}

type ValueEngine struct {
	values map[string]*ValueRecord
}

func NewValueEngine() *ValueEngine {
	return &ValueEngine{values: map[string]*ValueRecord{}}
}

func (ve *ValueEngine) AddValue(vr *ValueRecord) {
	ve.values[vr.Name] = vr
}

func (ve *ValueEngine) HasValue(name string) bool {
	return ve.values[name] != nil
}

func (ve *ValueEngine) EnvVars() []string {
	names := []string{}
	for _, vr := range ve.values {
		if vr.HasField("env") {
			names = append(names, vr.Env)
		}
	}
	return names
}

type parseState int

const (
	initial parseState = iota
	haveDollar
	haveLeftCurly
	haveValueName
)

func (ve *ValueEngine) ValidateValue(astring string) error {
	s := bufio.NewScanner(strings.NewReader(astring))
	s.Split(tokenize)

	var valueName string
	state := initial

	for s.Scan() {
		t := s.Text()
		switch state {
		case initial:
			if t == "$" {
				state = haveDollar
			}
		case haveDollar:
			switch t {
			case "$":
				state = initial
			case "{":
				state = haveLeftCurly
			default:
				return errors.New("Unexpected '" + t + "' after $")
			}
		case haveLeftCurly:
			switch t {
			case "$", "{", "}":
				return errors.New("Unexpected '" + t + "' after {")
			default:
				state = haveValueName
				valueName = t
			}
		case haveValueName:
			switch t {
			case "}":
				if !ve.HasValue(valueName) {
					return errors.New("Undefined value '" + valueName + "'")
				}
				state = initial
			default:
				return errors.New("Unexpected '" + t + "' after '" + valueName + "'")
			}
		}
	}

	if state != initial {
		return errors.New("Unexpected end of stuff '" + astring + "'")
	}

	return nil
}

func (ve *ValueEngine) ResolveValues(astring string) (string, error) {
	s := bufio.NewScanner(strings.NewReader(astring))
	s.Split(tokenize)

	var valueName string
	var result []byte
	state := initial

	for s.Scan() {
		t := s.Text()
		switch state {
		case initial:
			switch t {
			case "$":
				state = haveDollar
			default:
				result = append(result, []byte(t)...)
			}
		case haveDollar:
			switch t {
			case "$":
				result = append(result, '$')
				state = initial
			case "{":
				state = haveLeftCurly
			default:
				return "", errors.New("Unexpected token '" + t + "'")
			}
		case haveLeftCurly:
			switch t {
			case "}":
				return "", errors.New("Unexpected token '" + t + "'")
			default:
				valueName = t
				state = haveValueName
			}
		case haveValueName:
			switch t {
			case "}":
				result = append(result, []byte(ve.values[valueName].value())...)
				state = initial
			default:
				return "", errors.New("Unexpected token '" + t + "'")
			}
		}
	}

	return string(result), nil
}
