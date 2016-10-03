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

func TestValueValidationAndReplacement(t *testing.T) {
	ve := NewValueEngine()
	ve.AddValue(&ValueRecord{
		Name:          "qwerty",
		Literal:       "avalue",
		DecodedFields: []string{"Literal"},
	})

	testcases := []struct {
		input  string
		valid  bool
		output string
	}{
		{"", true, ""},
		{"$", false, "BAD"},
		{"ax$", false, "BAD"},
		{"$bx", false, "BAD"},
		{"cx$dx", false, "BAD"},
		{"$$", true, "$"},
		{"$$ex", true, "$ex"},
		{"fx$$", true, "fx$"},
		{"gx$$hx", true, "gx$hx"},
		{"{", true, "{"},
		{"ix{", true, "ix{"},
		{"{jx", true, "{jx"},
		{"lx{mx", true, "lx{mx"},
		{"nx${", false, "BAD"},
		{"ox${px", false, "BAD"},
		{"${qx", false, "BAD"},
		{"}", true, "}"},
		{"rx}", true, "rx}"},
		{"}sx", true, "}sx"},
		{"tx}ux", true, "tx}ux"},
		{"${foo}", false, "BAD"},
		{"${qwerty}", true, "avalue"},
		{"$${qwerty}", true, "${qwerty}"},
		{"zxcv-${qwerty}-lkjhg", true, "zxcv-avalue-lkjhg"},
	}

	for _, tc := range testcases {
		err := ve.Validate(tc.input)
		if tc.valid {
			if err != nil {
				t.Errorf("'%s' should have passed but got: %v", tc.input, err)
			} else {
				output, err := ve.Replace(tc.input)
				if err != nil {
					t.Error(err)
				} else if output != tc.output {
					t.Errorf(
						"'%s' should have produced '%s' but got '%s'",
						tc.input,
						tc.output,
						output,
					)
				}
			}
		} else if err == nil {
			t.Errorf("'%s' should have failed", tc.input)
		}
	}
}
