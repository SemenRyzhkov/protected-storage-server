package utils

import (
	"flag"
	"os"
)

func HandleFlag() {
	flag.Func("a", "GRPc server address", func(aFlagValue string) error {
		return os.Setenv("RUN_ADDRESS", aFlagValue)
	})

	flag.Func("d", "Address of db connection", func(dFlagValue string) error {
		return os.Setenv("DATABASE_URI", dFlagValue)
	})
}
