// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flag

import (
	"bufio"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type FlagSetExtras struct {
	// prefix to all env variable names /* jnovack/flag */
	envPrefix          string
	readUnderscoreFile bool
	trimFileContent    bool
}

var (
	// EnvironmentPrefix defines a string that will be implicitly prefixed to a
	// flag name before looking it up in the environment variables.
	EnvironmentPrefix = ""

	// DefaultConfigFlagname defines the flag name of the optional config file
	// path. Used to lookup and parse the config file when a default is set and
	// available on disk.
	DefaultConfigFlagname = "config"

	// If ReadUnderscoreFile is set to true, env parsing will look for envkey with _FILE suffix.
	// If ENVKEY_FILE contains a filepath, the corrensponding file will be read and its value be used.
	// This enables seemless support for Docker secrets.
	ReadUnderscoreFile = false

	// If TrimFileContent is set to true, any annoying new lines at the end
	// or preceding spaces will be trimmed. Whitespaces enclosed by other characters are not affected.
	// This is mainly thought for cases, where a new line might change an important key.
	TrimFileContent = false
)

func parseEnvToMap(environ []string) map[string]string {
	env := make(map[string]string)
	for _, s := range environ {
		i := strings.Index(s, "=")
		if i < 1 {
			continue
		}
		env[s[0:i]] = s[i+1:]
	}
	return env
}

func flagNameToEnvKey(name, envPrefix string) string {
	envKey := strings.ToUpper(name)
	if envPrefix != "" {
		envKey = envPrefix + "_" + envKey
	}
	// Replace all dashes (-) with underscores (_)
	envKey = strings.Replace(envKey, "-", "_", -1)
	// Replace all periods (.) with underscores (_)
	envKey = strings.Replace(envKey, ".", "_", -1)
	return envKey
}

// ParseEnv parses flags from environment variables.
// Flags already set will be ignored.
func (f *FlagSet) ParseEnv(environ []string) error {

	// Get the registered flags
	// flags := f.formal

	// Create a map of all environment variables
	env := parseEnvToMap(environ)

	// Iterate over all registered flags
	for _, registeredFlag := range f.formal {
		name := registeredFlag.Name
		// if flag has already been set, skip it
		_, exist := f.actual[name]
		if exist {
			continue
		}

		// seems like we are checking for unkown environment variables and help,
		// but i cant get my head around how it should be possible to have a miss here,
		// when we are iterating around range of f.formal
		// welp, dont touch something that aint broken
		flag, exist := f.formal[name]
		if !exist {
			if name == "help" || name == "h" { // special case for nice help message.
				f.usage()
				return ErrHelp
			}
			return f.failf("environment variable provided but not defined: %s", name)
		}

		envKey := flagNameToEnvKey(flag.Name, f.envPrefix)

		envValue, exist := env[envKey]
		if !exist {
			// parsing of _FILE
			if !f.readUnderscoreFile {
				continue
			}
			envKey = envKey + "_FILE"
			envValue, exist = env[envKey]
			if !exist {
				continue
			}
			if len(envValue) <= 0 {
				return f.failf("provided an _FILE env variable but it was empty")
			}
			fileBytes, err := os.ReadFile(envValue)
			if err != nil {
				return f.failf("could not read file %s provided by %s", envValue, envKey)
			}
			envValue = string(fileBytes)
			if f.trimFileContent {
				envValue = strings.TrimSpace(envValue)
			}
		}

		isEmpty := len(envValue) <= 0

		if fv, ok := flag.Value.(boolFlag); ok && fv.IsBoolFlag() { // special case: doesn't need an arg
			if !isEmpty {
				if err := fv.Set(envValue); err != nil {
					return f.failf("invalid boolean value %q for environment variable %s: %v", envValue, name, err)
				}
			} else {
				// flag without value is regarded a bool
				fv.Set("true")
			}
		} else {
			if err := flag.Value.Set(envValue); err != nil {
				return f.failf("invalid value %q for environment variable %s: %v", envValue, name, err)
			}
		}

		// update f.actual
		if f.actual == nil {
			f.actual = make(map[string]*Flag)
		}
		f.actual[name] = flag

	}
	return nil
}

// NewFlagSetWithExtras returns a new empty flag set with the specified name and error handling,
// as well as an environment variable prefix, and if ENVKEY_FILE should be supported.
func NewFlagSetWithExtras(name string, errorHandling ErrorHandling, envPrefix string, readUnderscoreFile bool, trimFileContent bool) *FlagSet {
	f := NewFlagSet(name, errorHandling)
	f.envPrefix = envPrefix
	f.readUnderscoreFile = readUnderscoreFile
	f.trimFileContent = trimFileContent
	return f
}

// ParseFile parses flags from the file in path.
//
// If the file is a YAML (.yaml, .yaml) file, it will be loaded as actual YAML.
func (f *FlagSet) ParseFile(path string) error {
	if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
		return f.parseFile_YAML(path)
	}

	return f.parseFile_PlainText(path)
}

// parseFile_PlainText parses flags from the file in path.
// Same format as commandline argumens, newlines and lines beginning with a
// "#" charater are ignored. Flags already set will be ignored.
func (f *FlagSet) parseFile_PlainText(path string) error {

	// Extract arguments from file
	fp, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file '%s': %v", path, err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text()

		// Ignore empty lines
		if len(line) == 0 {
			continue
		}

		// Ignore comments
		if line[:1] == "#" || line == "---" {
			continue
		}

		// Match `key=value` and `key value`
		var name, value string
		hasValue := false
		for i, v := range line {
			if v == '=' || v == ' ' || v == ':' {
				hasValue = true
				name, value = strings.TrimSpace(line[:i]), strings.TrimSpace(line[i+1:])

				// check if the name is an env name
				for srcName := range f.formal {
					if flagNameToEnvKey(srcName, f.envPrefix) == name {
						name = srcName
						break
					}
				}
				break
			}
		}

		if !hasValue {
			name = line
		}

		// Ignore flag when already set; arguments have precedence over file
		if f.actual[name] != nil {
			continue
		}

		m := f.formal
		flag, alreadythere := m[name]
		if !alreadythere {
			if name == "help" || name == "h" { // special case for nice help message.
				f.usage()
				return ErrHelp
			}
			return f.failf("configuration variable provided but not defined: %s", name)
		}

		if fv, ok := flag.Value.(boolFlag); ok && fv.IsBoolFlag() { // special case: doesn't need an arg
			if hasValue {
				if err := fv.Set(value); err != nil {
					return f.failf("invalid boolean value %q for configuration variable %s: %v", value, name, err)
				}
			} else {
				// flag without value is regarded a bool
				fv.Set("true")
			}
		} else {
			if err := flag.Value.Set(value); err != nil {
				return f.failf("invalid value %q for configuration variable %s: %v", value, name, err)
			}
		}

		// update f.actual
		if f.actual == nil {
			f.actual = make(map[string]*Flag)
		}
		f.actual[name] = flag
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

type yamlValue struct {
	Node  *yaml.Node
	Value string
	Error error
}

var (
	_ yaml.Unmarshaler = &yamlValue{} // ensure type compatibility
)

func (y *yamlValue) UnmarshalYAML(value *yaml.Node) error {
	y.Node = value

	if value.Kind != yaml.ScalarNode {
		y.Error = errors.New("only scalar/single values are supported")
		return nil
	}

	y.Value = value.Value

	return nil
}

func (f *FlagSet) parseFile_YAML(path string) error {
	// open the yaml file
	fp, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file '%s': %v", path, err)
	}
	defer fp.Close()

	// read the root object
	var values map[string]yamlValue

	err = yaml.NewDecoder(fp).Decode(&values)
	if err != nil {
		return fmt.Errorf("failed to parse file '%s': %v", path, err)
	}

	// parse the fields
	for name, value := range values {

		// check if the name is an env name
		for srcName := range f.formal {
			if flagNameToEnvKey(srcName, f.envPrefix) == name {
				name = srcName
				break
			}
		}

		// Ignore flag when already set; arguments have precedence over file
		if f.actual[name] != nil {
			continue
		}

		m := f.formal
		flag, alreadythere := m[name]
		if !alreadythere {
			if name == "help" || name == "h" { // special case for nice help message.
				f.usage()
				return ErrHelp
			}
			return f.failf("configuration variable provided but not defined: %s", name)
		}

		// forward error
		if value.Error != nil {
			return f.failf("invalid value %q for configuration variable %s at line %v: %v", value.Value, name, value.Node.Line, value.Error)
		}

		// set the flag value
		if err := flag.Value.Set(value.Value); err != nil {
			return f.failf("invalid value %q for configuration variable %s: %v", value.Value, name, err)
		}

		// update f.actual
		if f.actual == nil {
			f.actual = make(map[string]*Flag)
		}
		f.actual[name] = flag
	}

	return nil
}
