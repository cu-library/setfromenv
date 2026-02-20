# setfromenv
[![License MIT](https://img.shields.io/badge/license-MIT-lightgrey.svg?style=flat)](LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/cu-library/setfromenv.svg)](https://pkg.go.dev/github.com/cu-library/setfromenv)

A small Go library for setting unset flags from environment variables.

If a standard Go flag has not been set on the command line, this library can be used to easily pull the value for that flag
from an environment variable. The environment variable name is the flag name, uppercased, with hyphens replaced with underscores. 
An optional prefix can also be provided. 

| Prefix     | Flag Name           | Environment Variable      |
|------------|---------------------|---------------------------|
|            | `host`              | `HOST`                    |
|            | `log-level`         | `LOG_LEVEL`               |
|            | `http-port`         | `HTTP_PORT`               |
| `app`      | `host`              | `APP_HOST`                |
| `app_`     | `host`              | `APP_HOST`                |
| `svc`      | `http-port`         | `SVC_HTTP_PORT`           |
| `SVC`      | `enable-feature-x`  | `SVC_ENABLE_FEATURE_X`    |

The `EnvVarNameFromPrefix` function can also be used to help generate a helpful usage message.

## Example 

Here's an example of a small command line tool called 'scanner' with a flag which can be set
on the command line or from the environment. Set flags are not overwritten.

```go
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
```

Then, from the command line, assuming the resulting `scanner` program is in the local directory and executable:

```bash
$ ./scanner
Power level: 0
$ ./scanner -h
Scanner: Vegeta, what does the scanner say about his power level?
Version devel
Compiled with go1.25.7
  -power-level int
        power level

Environment variables read when flag is unset:
power-level: SCANNER_POWER_LEVEL
$ export SCANNER_POWER_LEVEL=1000
$ ./scanner
Power level: 1000
$ ./scanner -power-level 9001
Power level: 9001
It\'s over nine thousand!
$ export SCANNER_POWER_LEVEL="One hundred puppies."
$ ./scanner
unable to set flag "power-level" from environment variable "SCANNER_POWER_LEVEL": parse error
```
