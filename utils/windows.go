package utils

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"
)

func CheckForEdge() bool {
	if runtime.GOOS != "windows" {
		// for now we are concerned with windows only
		return true
	}
	c := exec.Command("powershell.exe", `(Get-AppxPackage Microsoft.MicrosoftEdge).Version`)
	if out, err := c.Output(); err == nil {
		if strings.TrimSpace(string(out)) != "" {
			return true
		}
		c = exec.Command("powershell.exe", `(Get-AppxPackage Microsoft.MicrosoftEdge.Stable).Version`)
		if out, err = c.Output(); err == nil {
			return strings.TrimSpace(string(out)) != ""
		}
	}
	return false
}

func OpenBrowser(url string) {
	// a copypaste
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}
