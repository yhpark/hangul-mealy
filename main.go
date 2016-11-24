package hangulmealy

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
)

var verboseFlag, choPriorityFlag = flag.Bool("v", false, "verbose"), flag.Bool("c", false, "cho priority (초성 우선)")

func readyStty() {
	restoreOnSignal := func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		go func() {
			for sig := range c {
				fmt.Println(sig)
				restoreStty()
				os.Exit(1)
			}
		}()
	}

	switch runtime.GOOS {
	case "linux":
		exec.Command("/bin/stty", "-F", "/dev/tty", "-icanon", "min", "1").Run()
		exec.Command("/bin/stty", "-F", "/dev/tty", "-echo").Run()
		restoreOnSignal()
	case "darwin":
		exec.Command("/bin/stty", "-f", "/dev/tty", "-icanon", "min", "1").Run()
		exec.Command("/bin/stty", "-f", "/dev/tty", "-echo").Run()
		restoreOnSignal()
	default:
		fmt.Println("Unbuffered typing is not supported in this operationg system: %v", runtime.GOOS)
	}
}

func restoreStty() {
	switch runtime.GOOS {
	case "linux":
		exec.Command("stty", "-F", "/dev/tty", "echo").Run()
	case "darwin":
		exec.Command("stty", "-f", "/dev/tty", "echo").Run()
	}
}

func main() {
	flag.Parse()

	var mealy, e = MakeHangulMealy(*choPriorityFlag)
	if e != nil {
		fmt.Println(e.Error())
		return
	}

	readyStty()
	defer restoreStty() // defer for panic

	b := make([]byte, 1)
	for {
		os.Stdin.Read(b)

		if !(32 <= b[0] && b[0] <= 127) {
			fmt.Println("Invalid Character " + string(b[0]) + " (char must be ascii 32~127)")
			continue
		}

		e := mealy.RunByRune(engByteToKorRune(b[0]))
		if e != nil {
			fmt.Println("Error - ", e)
		}

		fmt.Println(mealy.HangulString())
	}
}
