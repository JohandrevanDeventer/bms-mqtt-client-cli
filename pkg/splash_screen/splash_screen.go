package splash_screen

import (
	"fmt"

	"github.com/JohandrevanDeventer/bms-mqtt-client-cli/pkg/text_style"
)

// func PrintSplashScreen(appName, appVersion, buildVersion, developer, environment string) {
// 	goVersion := strings.Replace(runtime.Version(), "go", "", 1)

// 	fmt.Print(colorText(Blue, `
// ██████  ██    ██ ██████  ██  ██████  ██████  ███    ██     ██████  ███    ███ ███████
// ██   ██ ██    ██ ██   ██ ██ ██      ██    ██ ████   ██     ██   ██ ████  ████ ██
// ██████  ██    ██ ██████  ██ ██      ██    ██ ██ ██  ██     ██████  ██ ████ ██ ███████
// ██   ██ ██    ██ ██   ██ ██ ██      ██    ██ ██  ██ ██     ██   ██ ██  ██  ██      ██
// ██   ██  ██████  ██████  ██  ██████  ██████  ██   ████     ██████  ██      ██ ███████

// `))

// 	if appName == "" {
// 		appName = defaultAppName
// 	}

// 	if appVersion == "" {
// 		appVersion = defaultVersion
// 	}

// 	if buildVersion == "" {
// 		buildVersion = defaultBuildVersion
// 	}

// 	fullVersion := fmt.Sprintf("v%s-%s", appVersion, buildVersion)

// 	if developer == "" {
// 		developer = defaultDeveloper
// 	}

// 	switch environment {
// 	case "d":
// 		environment = "Development"
// 	case "t":
// 		environment = "Testing"
// 	case "p":
// 		environment = "Production"
// 	default:
// 		environment = defaultEnvironment
// 	}

// 	fmt.Printf("Welcome to %s!\n", colorText(Green, (boldText(appName))))
// 	fmt.Printf("Built with Go %s\n", colorText(Yellow, (boldText(goVersion))))
// 	fmt.Printf("Running version %s\n", colorText(Magenta, (boldText(fullVersion))))
// 	fmt.Printf("Developed by %s\n", colorText(Cyan, (boldText(developer))))
// }

func PrintSplashScreen() {

	fmt.Print(text_style.ColorText(text_style.Blue, `
██████  ██    ██ ██████  ██  ██████  ██████  ███    ██     ██████  ███    ███ ███████ 
██   ██ ██    ██ ██   ██ ██ ██      ██    ██ ████   ██     ██   ██ ████  ████ ██      
██████  ██    ██ ██████  ██ ██      ██    ██ ██ ██  ██     ██████  ██ ████ ██ ███████ 
██   ██ ██    ██ ██   ██ ██ ██      ██    ██ ██  ██ ██     ██   ██ ██  ██  ██      ██ 
██   ██  ██████  ██████  ██  ██████  ██████  ██   ████     ██████  ██      ██ ███████

`))
}
