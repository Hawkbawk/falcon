// The logger package contains helpful functions for quickly and beautifully logging things to
// stdout.
package logger

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

var DebugMode bool = viper.GetBool("PROX_DEBUG_MODE")

// Logs the given statement if debug mode is enabled.
func LogDebugOnly(statement ...interface{}) {
	if DebugMode {
		log.Print("PROX DEBUG STATEMENT:", statement)
	}
}

// Logs the given statement(s) in white text to stdout, after formatting it using fmt.Sprintf
func LogInfo(format string, substitutions ...interface{}) {
	formatted := fmt.Sprintf(format, substitutions...)
	color.White(formatted)
}

// Logs the given statement in red text to stdout, after formatting it using fmt.Sprintf
// and then ends the program with an exit code of 1.
func LogError(format string, substitutions ...interface{}) {
	formatted := fmt.Sprintf(format, substitutions...)
	color.Red(formatted)
	panic(1)
}
