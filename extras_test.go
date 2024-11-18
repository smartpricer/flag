// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flag

import (
	"gopkg.in/yaml.v3"
	"os"
	"syscall"
	"testing"
	"time"
)

// Test parsing a environment variables
func TestParseEnv(t *testing.T) {

	syscall.Setenv("BOOL", "")
	syscall.Setenv("BOOL2", "true")
	syscall.Setenv("INT", "22")
	syscall.Setenv("INT64", "0x23")
	syscall.Setenv("UINT", "24")
	syscall.Setenv("UINT64", "25")
	syscall.Setenv("STRING", "hello")
	syscall.Setenv("FLOAT64", "2718e28")
	syscall.Setenv("DURATION", "2m")

	f := NewFlagSet(os.Args[0], ContinueOnError)

	boolFlag := f.Bool("bool", false, "bool value")
	bool2Flag := f.Bool("bool2", false, "bool2 value")
	intFlag := f.Int("int", 0, "int value")
	int64Flag := f.Int64("int64", 0, "int64 value")
	uintFlag := f.Uint("uint", 0, "uint value")
	uint64Flag := f.Uint64("uint64", 0, "uint64 value")
	stringFlag := f.String("string", "0", "string value")
	float64Flag := f.Float64("float64", 0, "float64 value")
	durationFlag := f.Duration("duration", 5*time.Second, "time.Duration value")

	err := f.ParseEnv(os.Environ())
	if err != nil {
		t.Fatal("expected no error; got ", err)
	}
	if *boolFlag != false {
		t.Error("bool flag should be false, is ", *boolFlag)
	}
	if *bool2Flag != true {
		t.Error("bool2 flag should be true, is ", *bool2Flag)
	}
	if *intFlag != 22 {
		t.Error("int flag should be 22, is ", *intFlag)
	}
	if *int64Flag != 0x23 {
		t.Error("int64 flag should be 0x23, is ", *int64Flag)
	}
	if *uintFlag != 24 {
		t.Error("uint flag should be 24, is ", *uintFlag)
	}
	if *uint64Flag != 25 {
		t.Error("uint64 flag should be 25, is ", *uint64Flag)
	}
	if *stringFlag != "hello" {
		t.Error("string flag should be `hello`, is ", *stringFlag)
	}
	if *float64Flag != 2718e28 {
		t.Error("float64 flag should be 2718e28, is ", *float64Flag)
	}
	if *durationFlag != 2*time.Minute {
		t.Error("duration flag should be 2m, is ", *durationFlag)
	}
}

func TestParseEnvFile(t *testing.T) {
	expectedFileContent := "InHereCouldBeYourPreciosSecretYouWantToKeepSecure"
	syscall.Setenv("YOUR_SECRET_FILE", "./testdata/docker-secret")
	f := NewFlagSetWithExtras(os.Args[0], ContinueOnError, "", true, true)

	secretFlag := f.String("your-secret", "", "secret value")

	err := f.ParseEnv(os.Environ())
	if err != nil {
		t.Fatal("expected no error; got ", err)
	}
	if *secretFlag != expectedFileContent {
		t.Fatal("expected my secret; got ", *secretFlag)
	}
}

// Test parsing a configuration file
func TestParseFile(t *testing.T) {

	f := NewFlagSet(os.Args[0], ContinueOnError)

	boolFlag := f.Bool("bool", false, "bool value")
	bool2Flag := f.Bool("bool2", false, "bool2 value")
	intFlag := f.Int("int", 0, "int value")
	int64Flag := f.Int64("int64", 0, "int64 value")
	uintFlag := f.Uint("uint", 0, "uint value")
	uint64Flag := f.Uint64("uint64", 0, "uint64 value")
	stringFlag := f.String("string", "0", "string value")
	string2Flag := f.String("string2", "0", "string2 value")
	string3Flag := f.String("string3-env-like", "0", "string3 value")
	float64Flag := f.Float64("float64", 0, "float64 value")
	durationFlag := f.Duration("duration", 5*time.Second, "time.Duration value")

	err := f.ParseFile("./testdata/test.conf")
	if err != nil {
		t.Fatal("expected no error; got ", err)
	}
	if *boolFlag != true {
		t.Error("bool flag should be true, is ", *boolFlag)
	}
	if *bool2Flag != true {
		t.Error("bool2 flag should be true, is ", *bool2Flag)
	}
	if *intFlag != 22 {
		t.Error("int flag should be 22, is ", *intFlag)
	}
	if *int64Flag != 0x23 {
		t.Error("int64 flag should be 0x23, is ", *int64Flag)
	}
	if *uintFlag != 24 {
		t.Error("uint flag should be 24, is ", *uintFlag)
	}
	if *uint64Flag != 25 {
		t.Error("uint64 flag should be 25, is ", *uint64Flag)
	}
	if *stringFlag != "hello" {
		t.Error("string flag should be `hello`, is ", *stringFlag)
	}
	if *string2Flag != "world" {
		t.Error("string2 flag should be `world`, is", *string2Flag)
	}
	if *string3Flag != "foo" {
		t.Error("string3-env-like flag should be `foo`, is", *string2Flag)
	}
	if *float64Flag != 2718e28 {
		t.Error("float64 flag should be 2718e28, is ", *float64Flag)
	}
	if *durationFlag != 2*time.Minute {
		t.Error("duration flag should be 2m, is ", *durationFlag)
	}
}

func TestParseFileUnknownFlag(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	if err := f.ParseFile("./testdata/bad_test.conf"); err != nil {
		t.Error("parse did not ignore unknown flag; ", err)
	}
}

func TestDefaultConfigFlagname(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)

	f.Bool("bool", false, "bool value")
	f.Bool("bool2", false, "bool2 value")
	f.Int("int", 0, "int value")
	f.Int64("int64", 0, "int64 value")
	f.Uint("uint", 0, "uint value")
	f.Uint64("uint64", 0, "uint64 value")
	stringFlag := f.String("string", "0", "string value")
	f.String("string2", "0", "string value")
	f.String("string3-env-like", "0", "string value")
	f.Float64("float64", 0, "float64 value")
	f.Duration("duration", 5*time.Second, "time.Duration value")

	f.String(DefaultConfigFlagname, "./testdata/test.conf", "config path")

	if err := os.Unsetenv("STRING"); err != nil {
		t.Error(err)
	}

	if err := f.Parse([]string{}); err != nil {
		t.Error("parse failed; ", err)
	}

	if *stringFlag != "hello" {
		t.Error("string flag should be `hello`, is", *stringFlag)
	}
}

func TestDefaultConfigFlagnameYAML(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)

	stringFlag := f.String("string", "0", "string value")

	f.String(DefaultConfigFlagname, "./testdata/test.yml", "config path")

	if err := os.Unsetenv("STRING"); err != nil {
		t.Error(err)
	}

	if err := f.Parse([]string{}); err != nil {
		t.Error("parse failed; ", err)
	}

	if *stringFlag != "helloYAML" {
		t.Error("string flag should be `hello`, is", *stringFlag)
	}
}

func TestDefaultConfigFlagnameMultiple(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)

	f.Bool("bool", false, "bool value")
	f.Bool("bool2", false, "bool2 value")
	f.Int("int", 0, "int value")
	f.Int64("int64", 0, "int64 value")
	f.Uint("uint", 0, "uint value")
	f.Uint64("uint64", 0, "uint64 value")
	stringFlag := f.String("string", "0", "string value")
	f.String("string2", "0", "string value")
	f.String("string3-env-like", "0", "string value")
	f.Float64("float64", 0, "float64 value")
	f.Duration("duration", 5*time.Second, "time.Duration value")

	f.String(DefaultConfigFlagname, "./testdata/test.yml:./testdata/test.conf", "config path")

	if err := os.Unsetenv("STRING"); err != nil {
		t.Error(err)
	}

	if err := f.Parse([]string{}); err != nil {
		t.Error("parse failed; ", err)
	}

	if *stringFlag != "helloYAML" {
		t.Error("string flag should be `hello`, is", *stringFlag)
	}
}

func TestDefaultConfigFlagnameMissingFile(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	f.String(DefaultConfigFlagname, "./testdata/missing", "config path")

	if err := os.Unsetenv("STRING"); err != nil {
		t.Error(err)
	}
	if err := f.Parse([]string{}); err == nil {
		t.Error("expected error of missing config file, got nil")
	}
}

func TestFlagSetParseErrors(t *testing.T) {
	fs := NewFlagSet("test", ContinueOnError)
	fs.Int("int", 0, "int value")

	args := []string{"-int", "bad"}
	expected := `invalid value "bad" for flag -int: parse error`
	if err := fs.Parse(args); err == nil || err.Error() != expected {
		t.Errorf("expected error %q parsing from args, got: %v", expected, err)
	}

	if err := os.Setenv("INT", "bad"); err != nil {
		t.Fatalf("error setting env: %s", err.Error())
	}
	expected = `invalid value "bad" for environment variable int: parse error`
	if err := fs.Parse([]string{}); err == nil || err.Error() != expected {
		t.Errorf("expected error %q parsing from env, got: %v", expected, err)
	}
	if err := os.Unsetenv("INT"); err != nil {
		t.Fatalf("error unsetting env: %s", err.Error())
	}

	fs.String("config", "", "config filename")
	args = []string{"-config", "testdata/bad_test.conf"}
	expected = `invalid value "bad" for configuration variable int: parse error`
	if err := fs.Parse(args); err == nil || err.Error() != expected {
		t.Errorf("expected error %q parsing from config, got: %v", expected, err)
	}
}

func TestFlagSetBooleanVariants(t *testing.T) {
	expectations := []struct {
		Input  string
		Output bool
	}{
		{"true", true},
		{"false", false},
		{"True", true},
		{"False", false},
		{"1", true},
		{"0", false},
		{"yes", true},
		{"no", false},
		{"y", true},
		{"n", false},
		{"on", true},
		{"off", false},
		{"ON", true},
		{"OFF", false},
		{"", false},
	}

	for _, expectation := range expectations {
		f := NewFlagSet("test", ContinueOnError)
		v := f.Bool("my-boolean", false, "test boolean")

		if err := os.Setenv("MY_BOOLEAN", expectation.Input); err != nil {
			t.Error(err)
		}
		if err := f.Parse([]string{}); err != nil {
			t.Error(err)
		}
		if *v != expectation.Output {
			t.Errorf("expected '%v' to be %v", expectation.Input, expectation.Output)
		}
		if err := os.Unsetenv("MY_BOOLEAN"); err != nil {
			t.Error(err)
		}
	}
}

func TestTestingPackageFlags(t *testing.T) {
	f := NewFlagSet("test", ContinueOnError)
	if err := f.Parse([]string{"-test.v", "-test.count", "1"}); err == nil {
		t.Error(err)
	}
}

func TestYAMLValues_Empty(t *testing.T) {
	var v map[string]yamlValue

	err := yaml.Unmarshal([]byte(""), &v)
	if err != nil {
		t.Fatal(err)
	}

	if len(v) > 0 {
		t.Error("expected empty")
	}
}

func TestYAMLValues_Simple(t *testing.T) {
	var v map[string]yamlValue

	err := yaml.Unmarshal([]byte(`
str: hello_world
str2: "hello_world"
str3: 'hello_world'
ml: |
  line 1
  line 2
i: 12345
f: 12345.6789
b: true
obj:
  ignore: me
ar:
  - ignore
  - me
`), &v)
	if err != nil {
		t.Fatal(err)
	}

	if vv, ok := v["str"]; !ok || vv.Value != "hello_world" {
		t.Error("unexpected value for str")
	}
	if vv, ok := v["str2"]; !ok || vv.Value != "hello_world" {
		t.Error("unexpected value for str2")
	}
	if vv, ok := v["str3"]; !ok || vv.Value != "hello_world" {
		t.Error("unexpected value for str3")
	}
	if vv, ok := v["ml"]; !ok || vv.Value != "line 1\nline 2\n" {
		t.Error("unexpected value for ml")
	}
	if vv, ok := v["i"]; !ok || vv.Value != "12345" {
		t.Error("unexpected value for i")
	}
	if vv, ok := v["f"]; !ok || vv.Value != "12345.6789" {
		t.Error("unexpected value for f")
	}
	if vv, ok := v["b"]; !ok || vv.Value != "true" {
		t.Error("unexpected value for b")
	}
	if vv, ok := v["obj"]; !ok || vv.Value != "" || vv.Error == nil {
		t.Error("unexpected value for obj")
	}
	if vv, ok := v["ar"]; !ok || vv.Value != "" || vv.Error == nil {
		t.Error("unexpected value for ar")
	}
}
