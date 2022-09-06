package image

import (
	"strings"
	"testing"
)

func TestBuildName(t *testing.T) {
	assert_eq := func(input, validation string) {
		res := BuildName(input)
		if res != validation {
			t.Errorf("%q != %q", res, validation)
		}
	}

	input := "docker/Dockerfile.app"
	validation := "app"
	assert_eq(input, validation)

	input = "docker/app/Dockerfile"
	validation = "docker-app"
	assert_eq(input, validation)

	input = "docker/app/Dockerfile.image-name"
	validation = "image-name"
	assert_eq(input, validation)

}

func TestParseImageSource(t *testing.T) {
	// FROM [--platform=<platform>] <image> [AS <name>]
	// FROM [--platform=<platform>] <image>[:<tag>] [AS <name>]
	// FROM [--platform=<platform>] <image>[@<digest>] [AS <name>]
	//
	// FROM {1} [--platform=<platform>] {2} <image>[@<digest>] {3} [AS <name>]
	assert_eq := func(input, validation string) {
		res := parseImageSource(input)
		if res != validation {
			t.Errorf("%q != %q", res, validation)
		}
	}

	input := "FROM debian:stable-slim"
	validation := "debian:stable-slim"
	assert_eq(input, validation)

	input = "FROM debian:stable-slim AS builder"
	validation = "debian:stable-slim"
	assert_eq(input, validation)

	input = "FROM --platform=linux/arm/v7 debian:stable-slim"
	validation = "debian:stable-slim"
	assert_eq(input, validation)

	input = "FROM --platform=linux/arm/v7 debian:stable-slim AS builder"
	validation = "debian:stable-slim"
	assert_eq(input, validation)

	input = "FROM --platform=linux/arm/v7 debian:stable-slim"
	validation = "debian:stable-slim"
	assert_eq(input, validation)

}

func TestParseImageStr(t *testing.T) {
	assert_eq := func(input string, validation Image) {
		res := ParseImageStr(input)
		if res != validation {
			t.Errorf("%#v != %#v", res, validation)
		}
	}

	input := "debian"
	validation := Image{Name: "debian"}
	assert_eq(input, validation)

	input = "debian:stable-slim"
	validation = Image{Name: "debian", Tag: "stable-slim"}
	assert_eq(input, validation)

	input = "debian@1234"
	validation = Image{Name: "debian", Hash: "1234"}
	assert_eq(input, validation)

}

func TestImageSource(t *testing.T) {
	assert_eq := func(input string, validation Image) {
		r := strings.NewReader(input)
		res, err := ImageSource(r)
		if err != nil {
			t.Errorf("%s", err)
		}
		if res != validation {
			t.Errorf("%#v != %#v", res, validation)
		}
	}

	input := "FROM debian"
	validation := Image{Name: "debian"}
	assert_eq(input, validation)

	input = `
# some comments
FROM debian
	`
	validation = Image{Name: "debian"}
	assert_eq(input, validation)

	// input = "FROM debian:stable-slim AS builder"
	// validation = Image{Name: "debian", Tag: "stable-slim"}
	// assert_eq(input, validation)

	// input = "FROM debian:stable-slim AS builder"
	// validation = Image{Name: "debian"}
	// assert_eq(input, validation)

}
