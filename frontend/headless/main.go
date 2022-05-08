package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ivanizag/izapple2"
	"github.com/ivanizag/izapple2/screen"
)

func main() {
	a := izapple2.MainApple()
	fe := &headLessFrontend{}
	fe.keyChannel = make(chan uint8, 200)
	a.SetKeyboardProvider(fe)
	go a.Run()

	inReader := bufio.NewReader(os.Stdin)
	running := true
	for running {
		fmt.Print("* ")
		text, err := inReader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		text = strings.TrimSpace(text)
		parts := strings.Split(text, " ")
		command := strings.ToLower(parts[0])
		switch command {
		case "exit":
			a.SendCommand(izapple2.CommandKill)
			running = false

		case "pts":
			fallthrough
		case "printtextscreen":
			fmt.Print(izapple2.DumpTextModeAnsi(a))

		case "ss":
			fallthrough
		case "savescreen":
			err := screen.SaveSnapshot(a, screen.ScreenModeNTSC, "snapshot.png")
			if err != nil {
				fmt.Printf("Error saving screen: %v.\n.", err)
			} else {
				fmt.Println("Saving screen 'snapshot.png'")
			}

		case "ssm":
			fallthrough
		case "savescreenmono":
			err := screen.SaveSnapshot(a, screen.ScreenModePlain, "snapshot.png")
			if err != nil {
				fmt.Printf("Error saving screen: %v.\n.", err)
			} else {
				fmt.Println("Saving screen 'snapshot.png'")
			}

		case "k":
			fallthrough
		case "key":
			if len(parts) < 2 {
				fmt.Println("No key specified.")
			} else {
				key := uint8(parts[1][0])
				fe.keyChannel <- key
			}

		case "ks":
			fallthrough
		case "keys":
			text := strings.Join(parts[1:], " ")
			for _, char := range text {
				fe.keyChannel <- uint8(char)
			}

		case "kr":
			text := strings.Join(parts[1:], " ")
			for _, char := range text {
				fe.keyChannel <- uint8(char)
			}
			fe.keyChannel <- 13

		case "r":
			fallthrough
		case "return":
			fe.keyChannel <- 13

		case "help":
			fmt.Print(`
Available commands:
	Exit: Stops the emulator and quits
	PrintTextScreen or pts: Prints the text mode screen
	PrintTextScreen, pts: Prints the text mode screen
	SaveScreen or ss: Saves the screen with NTSC colors to "snapshot.png"
	SaveScreenMono or ssm: Saves the monochromatic screen to "snapshot.png"
	Key or k: Sends a key to the emulator
	Keys or ks: Sends a string to the emulator
	Return or r: Sends a return to the emulator
	Help: Prints this help
`)
		default:
			fmt.Println("Unknown command.")
		}
	}
}

/*
Uses the console to send commands and queries to an emulated machine.
*/
type headLessFrontend struct {
	keyChannel chan uint8
}

func (fe *headLessFrontend) GetKey(strobed bool) (key uint8, ok bool) {
	if !strobed {
		// We must use the strobe to control the flow from stdin
		ok = false
		return
	}

	select {
	case key = <-fe.keyChannel:
		ok = true
	default:
		ok = false
	}
	return
}
