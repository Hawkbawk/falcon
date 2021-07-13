package lib

import (
	"log"

	"github.com/spf13/viper"
)


var DebugMode bool = viper.GetBool("PROX_DEBUG_MODE")


func LogDebugOnly(statement ...interface{}) {
	if DebugMode {
		log.Print("PROX DEBUG STATEMENT:", statement)
	}
}