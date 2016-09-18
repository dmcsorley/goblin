package main

import (
	"testing"
)

func expectLoadConfigBytesFailure(t *testing.T, jsonString string, errorMessage string) {
	_, err := loadConfigBytes([]byte(jsonString))
	if err == nil {
		t.Error(errorMessage)
	}
}

func TestEmptyServerConfig(t *testing.T) {
	expectLoadConfigBytesFailure(
		t,
		"",
		"should have failed on empty server config",
	)
}

func TestServerConfigRequiresBuild(t *testing.T) {
	expectLoadConfigBytesFailure(
		t,
		"{}",
		"should have failed when server config has no builds",
	)
}

func TestServerConfigNonObjectBuilds(t *testing.T) {
	expectLoadConfigBytesFailure(
		t,
		`{"builds":[1, 2, 3]}`,
		"should have failed when builds are not JSON objects",
	)
}

func TestBuildConfigRequiresName(t *testing.T) {
	expectLoadConfigBytesFailure(
		t,
		`{"builds":[{}]}`,
		"should have failed when build has no name",
	)
}

func TestBuildConfigRequiresNameIsString(t *testing.T) {
	expectLoadConfigBytesFailure(
		t,
		`{"builds":[{"name":true}]}`,
		"should have failed when build name is not string",
	)
}

func TestBuildConfigRequiresStep(t *testing.T) {
	expectLoadConfigBytesFailure(
		t,
		`{"builds":[{"name":"aname"}]}`,
		"should have failed when build has no steps",
	)
}
