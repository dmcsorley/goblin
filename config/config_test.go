// import github.com/dmcsorley/goblin/config
package config

import (
	"reflect"
	"testing"
)

func TestLoadHCLBytes(t *testing.T) {
	cs := `build "foo" {
  step git-clone {
    url = "foo"
  }
}
build "bar" {
  step docker-build {
    image = "example"
  }
}`

	c, err := LoadBytes([]byte(cs))
	if err != nil {
		t.Error(err)
	}

	expected := &Record{
		Builds: map[string]*BuildRecord{
			"foo": &BuildRecord{
				Steps: []*StepRecord{
					&StepRecord{
						Type:          "git-clone",
						Url:           "foo",
						DecodedFields: []string{"Url"},
					},
				},
			},
			"bar": &BuildRecord{
				Steps: []*StepRecord{
					&StepRecord{
						Type:          "docker-build",
						Image:         "example",
						DecodedFields: []string{"Image"},
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(c, expected) {
		t.Error("Not equal")
	}
}

func TestLoadJSONBytes(t *testing.T) {
	cs := `{ "build": {
  "foo": { "step": [
    { "git-clone": { "url": "foo" } }
  ] },
  "bar": { "step": [
    { "docker-build": { "image": "example" } }
  ] }
} }`

	c, err := LoadBytes([]byte(cs))
	if err != nil {
		t.Error(err)
	}

	expected := &Record{
		Builds: map[string]*BuildRecord{
			"foo": &BuildRecord{
				Steps: []*StepRecord{
					&StepRecord{
						Type:          "git-clone",
						Url:           "foo",
						DecodedFields: []string{"Url"},
					},
				},
			},
			"bar": &BuildRecord{
				Steps: []*StepRecord{
					&StepRecord{
						Type:          "docker-build",
						Image:         "example",
						DecodedFields: []string{"Image"},
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(c, expected) {
		t.Error("Not equal")
	}
}

func TestLoadValues(t *testing.T) {
	cs := `
value "avalue" {
  literal = "example"
}
value "bvalue" {
  literal = "example2"
}
build "foo" {
  step git-clone {
    url = "${avalue}"
  }
}`

	c, err := LoadBytes([]byte(cs))
	if err != nil {
		t.Error(err)
	}

	expected := &Record{
		Values: []*ValueRecord{
			&ValueRecord{
				Name:          "avalue",
				Literal:       "example",
				DecodedFields: []string{"Literal"},
			},
			&ValueRecord{
				Name:          "bvalue",
				Literal:       "example2",
				DecodedFields: []string{"Literal"},
			},
		},
		Builds: map[string]*BuildRecord{
			"foo": &BuildRecord{
				Steps: []*StepRecord{
					&StepRecord{
						Type:          "git-clone",
						Url:           "${avalue}",
						DecodedFields: []string{"Url"},
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(c, expected) {
		t.Error("Not equal")
	}
}
