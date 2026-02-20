// Copyright 2026 Carleton University Library
// All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// Package setfromenv is a library which sets unset flags from environment variables.
package setfromenv

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

// SetFlagsInFlagSet sets unset flags using environment variables.
// It finds unset flags in fs, then sets those flags using the value of an
// environment variable.
// To help group environment variables used by an application, an optional
// prefix can be provided. Prefixes that do not end in a '_' character have
// a '_' character appended.
// The name of the environment variable to lookup is the name of the unset
// flag, with the optional prefix prepended, converted to uppercase, and
// with any '-' characters replaced with '_' characters.
func SetFlagsInFlagSet(fs *flag.FlagSet, prefix string) error {
	// Create a closure which can be used to generate environment variable
	// names.
	makeEnvName := EnvVarNameFromPrefix(prefix)

	// The set of unset flag names.
	unset := make(map[string]struct{})

	// Visit calls a function on "only those flags that have been set."
	// VisitAll calls a function on "all flags, even those not set."
	// No way to ask for "only unset flags". So, we add all, then
	// delete the set flags.

	// First, visit all the flags, and add them to the set.
	fs.VisitAll(func(f *flag.Flag) { unset[f.Name] = struct{}{} })

	// Then delete the set flags.
	fs.Visit(func(f *flag.Flag) { delete(unset, f.Name) })

	for flagName := range unset {
		envName := makeEnvName(flagName)
		// Look for the environment variable.
		// If found, set the flag to that variable's value.
		// If there's a problem setting the value, return an error.
		value, found := os.LookupEnv(envName)
		if found {
			err := fs.Set(flagName, value)
			if err != nil {
				return fmt.Errorf("unable to set flag %q from environment variable %q: %w",
					flagName, envName, err)
			}
		}
	}
	return nil

}

// SetFlags sets unset flags in the flag package's CommandLine FlagSet
// from corresponding environment variables.
// This function will return an error if flag.Parse() has not been called first.
func SetFlags(prefix string) error {
	if !flag.Parsed() {
		return errors.New("command-line arguments not yet parsed, flag.Parse() should be called first")
	}
	return SetFlagsInFlagSet(flag.CommandLine, prefix)
}

// EnvVarNameFromPrefix returns a function which generates environment variable
// names from flag names, using an optional prefix.
func EnvVarNameFromPrefix(prefix string) func(string) string {
	// Add a trailing '_' character to the prefix if it is needed.
	if prefix != "" && !strings.HasSuffix(prefix, "_") {
		prefix += "_"
	}
	return func(flagName string) string {
		envName := prefix + flagName
		envName = strings.ReplaceAll(envName, "-", "_")
		envName = strings.ToUpper(envName)
		return envName
	}
}
