/*
Copyright © 2025 Johandré van Deventer <johandre.vandeventer@rubiconsa.com>
*/
package main

import (
	"github.com/JohandrevanDeventer/bms-mqtt-client-cli/cmd"
	"github.com/JohandrevanDeventer/bms-mqtt-client-cli/pkg/splash_screen"
)

func main() {
	splash_screen.PrintSplashScreen()
	cmd.Execute()
}
