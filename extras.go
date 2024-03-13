// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flag

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var (
// EnvironmentPrefix defines a string that will be implicitly prefixed to a
// flag name before looking it up in the environment variables.
	EnvironmentPrefix = ""

	// DefaultConfigFlagname defines the flag name of the optional config file
	// path. Used to lookup and parse the config file when a default is set and
	// available on disk.
	DefaultConfigFlagname = "config"

	// If ReadValueFromUnderscoreFile is set to true, env parsing will look for envkey with _FILE suffix.
	// If ENVKEY_FILE contains a filepath, the corrensponding file will be read and its value be used.
	// This enables seemless support for Docker secrets.
	ReadValueFromUnderscoreFile = false

	// If WhitespaceTrimUnderscoreFileContent is set to true, any annoying new lines at the end
	// or preceding spaces will be trimmed. Whitespaces enclosed by other characters are not affected.
	// This is mainly thought for cases, where a new line might change an important key.
	WhitespaceTrimUnderscoreFileContent = false
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

		// TODO: can it even be possible to hit a non-existing flag due to the range func loop?
		// TODO: who on earth would call help from an environment variable?
		// seems like we are checking for unkown

		//EXPLAIN
		flag, exist := f.formal[name]
		if !exist {
			fmt.Printf("FERRIS: found a flag that does not exist: %s", name)
			if name == "help" || name == "h" { // special case for nice help message.
				f.usage()
				return ErrHelp
			}
			return f.failf("environment variable provided but not defined: %s", name)
		}

		envKey := flagNameToEnvKey(flag.Name, f.envPrefix)

		envValue, exist := env[envKey]
		if !exist {
			continue
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

// NewFlagSetWithEnvPrefix returns a new empty flag set with the specified name,
// environment variable prefix, and error handling property.
func NewFlagSetWithEnvPrefix(name string, prefix string, errorHandling ErrorHandling) *FlagSet {
	f := NewFlagSet(name, errorHandling)
	f.envPrefix = prefix
	return f
}

// ParseFile parses flags from the file in path.
// Same format as commandline argumens, newlines and lines beginning with a
// "#" charater are ignored. Flags already set will be ignored.
func (f *FlagSet) ParseFile(path string) error {

	// Extract arguments from file
	fp, err := os.Open(path)
	if err != nil {
		return err
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
		if line[:1] == "#" {
			continue
		}

		// Match `key=value` and `key value`
		var name, value string
		hasValue := false
		for i, v := range line {
			if v == '=' || v == ' ' {
				hasValue = true
				name, value = line[:i], line[i+1:]
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
