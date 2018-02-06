package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

var (
	notes = map[string]int{
		"B0":  31,
		"C1":  33,
		"CS1": 35,
		"D1":  37,
		"DS1": 39,
		"E1":  41,
		"F1":  44,
		"FS1": 46,
		"G1":  49,
		"GS1": 52,
		"A1":  55,
		"AS1": 58,
		"B1":  62,
		"C2":  65,
		"CS2": 69,
		"D2":  73,
		"DS2": 78,
		"E2":  82,
		"F2":  87,
		"FS2": 93,
		"G2":  98,
		"GS2": 104,
		"A2":  110,
		"AS2": 117,
		"B2":  123,
		"C3":  131,
		"CS3": 139,
		"D3":  147,
		"DS3": 156,
		"E3":  165,
		"F3":  175,
		"FS3": 185,
		"G3":  196,
		"GS3": 208,
		"A3":  220,
		"AS3": 233,
		"B3":  247,
		"C4":  262,
		"CS4": 277,
		"D4":  294,
		"DS4": 311,
		"E4":  330,
		"F4":  349,
		"FS4": 370,
		"G4":  392,
		"GS4": 415,
		"A4":  440,
		"AS4": 466,
		"B4":  494,
		"C5":  523,
		"CS5": 554,
		"D5":  587,
		"DS5": 622,
		"E5":  659,
		"F5":  698,
		"FS5": 740,
		"G5":  784,
		"GS5": 831,
		"A5":  880,
		"AS5": 932,
		"B5":  988,
		"C6":  1047,
		"CS6": 1109,
		"D6":  1175,
		"DS6": 1245,
		"E6":  1319,
		"F6":  1397,
		"FS6": 1480,
		"G6":  1568,
		"GS6": 1661,
		"A6":  1760,
		"AS6": 1865,
		"B6":  1976,
		"C7":  2093,
		"CS7": 2217,
		"D7":  2349,
		"DS7": 2489,
		"E7":  2637,
		"F7":  2794,
		"FS7": 2960,
		"G7":  3136,
		"GS7": 3322,
		"A7":  3520,
		"AS7": 3729,
		"B7":  3951,
		"C8":  4186,
		"CS8": 4435,
		"D8":  4699,
		"DS8": 4978,
	}
	music_sheet []byte
)

func init() {
	var filename *string
	var err error

	filename = flag.String("music_sheet", "", "the music sheet to load and play")
	flag.Parse()
	music_sheet, err = ioutil.ReadFile(*filename)
	if err != nil {
		panic(err)
	}
}

func main() {
	var err error
	fmt.Println("beep-jam 0.1")

	beeper, err := NewBeeper()
	if err != nil {
		fmt.Println("Could not create beeper")
		panic(err)
	}

	// Create the handler for SIGINT
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("\n\nGot Ctrl-C'd")
		beeper.Beep(0.0, 1)
		os.Exit(1)
	}()

	s := string(music_sheet[:])
	lines := strings.Split(s, "\n")
	numLines := len(lines)
	for idx, l := range lines {
		l := strings.TrimSpace(l)
		if l == "" {
			continue
		}
		if strings.HasPrefix(l, "PAUSE") {
			parts := strings.Split(l, " ")
			delay, err := strconv.Atoi(parts[1])
			if err != nil {
				panic(err)
			}

			time.Sleep(time.Duration(delay) * time.Millisecond)
		} else if strings.HasPrefix(l, ";") {
			// Do nothing, it's a comment
		} else if strings.HasPrefix(l, "FREQ") {
			parts := strings.Split(l, " ")
			freq, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Error reading", l)
				panic(err)
			}

			dur, err := strconv.Atoi(parts[2])
			if err != nil {
				fmt.Println("Error reading", l)
				panic(err)
			}

			beeper.Beep(float32(freq), dur)
		} else {
			parts := strings.Split(l, " ")
			note := parts[0]
			if _, ok := notes[note]; !ok {
				panic("Unknown note: " + note)
			}
			duration, err := strconv.Atoi(parts[1])

			if err != nil {
				panic(err)
			}

			beeper.Beep(float32(notes[note]), duration)
		}

		progress := int(float64(idx) / float64(numLines) * 100)

		fmt.Printf("\r[")
		for i := 0; i < progress-1; i++ {
			fmt.Printf("=")
		}
		fmt.Printf(">")

		for i := progress; i < 100; i++ {
			fmt.Printf(" ")
		}

		fmt.Printf("] %d%% ", progress)
	}
	fmt.Print("\n")
}
