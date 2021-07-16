package logger

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

var DebugMode bool = viper.GetBool("PROX_DEBUG_MODE")

func LogDebugOnly(statement ...interface{}) {
	if DebugMode {
		log.Print("PROX DEBUG STATEMENT:", statement)
	}
}

func LogInfo(statement ...interface{}) {
	formatted := fmt.Sprint(statement...)
	color.White(formatted)
}

func LogError(statement ...interface{}) {
	formatted := fmt.Sprint(statement...)
	color.Red(formatted)
	panic(1)
}