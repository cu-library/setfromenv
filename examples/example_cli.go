package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/cu-library/setfromenv"
)

// A version flag, which should be overwritten when building using ldflags.
var version = "devel"

const (
	EnvPrefix = "SCANNER"
)

func main() {
	v := flag.Int("power-level", 0, "power level")
	//nolint:errcheck
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Scanner: Vegeta, what does the scanner say about his power level?\n")
		fmt.Fprintf(os.Stderr, "Version %v\n", version)
		fmt.Fprintf(flag.CommandLine.Output(), "Compiled with %v\n", runtime.Version())
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nEnvironment variables read when flag is unset:")
		makeEnvName := setfromenv.EnvVarNameFromPrefix(EnvPrefix)
		flag.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(os.Stderr, "%v: %v\n", f.Name, makeEnvName(f.Name))
		})
	}
	flag.Parse()
	err := setfromenv.SetFlags(EnvPrefix)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Power level: %v\n", *v)
	if *v > 9000 {
		fmt.Println("It's over nine thousand!")
	}
}
