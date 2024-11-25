package g90updatefw

import (
	"errors"
	"github.com/dalefarnsworth/go-xmodem/xmodem"
	"io"
	"log"
	"strings"
	"time"
)

// Most of the codes here comes from Dale Farnsworth.
// Thanks!

const buflen = 64 * 1024

func readString(serial *Serial) string {
	buf := make([]byte, buflen)

	i := 0
	lastReadZeroBytes := false
	for {
		n, err := serial.Read(buf[i:])
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		if n == 0 {
			if i == 0 {
				continue
			}
			if lastReadZeroBytes {
				// only return
				break
			}
			lastReadZeroBytes = true
			continue
		}
		lastReadZeroBytes = false
		//syscall.Write(syscall.Stdout, buf[i:i+n])
		i += n
	}

	if i >= len(buf) {
		log.Fatal(errors.New("Read buffer overrun"))
	}

	return string(buf[0:i])
}

func expect(serial *Serial, expects []string) (expectIndex int) {
	previousStr := ""

	for {
		str := readString(serial)
		for i, expect := range expects {
			if strings.Contains(previousStr+str, expect) {
				return i
			}
		}
		previousStr = str
	}

	panic("unreachable")
}

func expectSend(serial *Serial, expects []string, sends []string) (whichExpect string) {
	if len(sends) != len(expects) {
		panic("length of sends array does not equal length of expects array")
	}

	//fmt.Printf("> Waiting for '%s'...\n\n", strings.Join(expects, "' or '"))

	expectIndex := expect(serial, expects)
	send := sends[expectIndex]

	if len(send) != 0 {
		_, err := serial.Write([]byte(send))
		if err != nil {
			log.Fatal(err)
		}
	}

	return expects[expectIndex]
}

func UpdateRadio(serial *Serial, data []byte, progChan chan<- uint) {
	attentionTimeout := 10 * time.Millisecond
	menuTimeout := 50 * time.Millisecond
	eraseTimeout := 50 * time.Millisecond
	uploadTimeout := 10 * time.Second
	cleanupTimeout := 500 * time.Millisecond

	banner := "Hit a key to abort"
	menu := "1.Update FW"
	waitFW := "Wait FW file"

	attentionGrabber := " "
	menuSelector := "1"

	// SPINNER1, waiting ready
	progChan <- 1

	serial.Flush()

	expects := []string{banner, menu}
	sends := []string{attentionGrabber, menuSelector}

	serial.SetReadTimeout(attentionTimeout)
	found := expectSend(serial, expects, sends)

	// fmt.Println()
	// SPINNER2, erasing and waiting fw...
	progChan <- 2

	if found != menu {
		serial.SetReadTimeout(menuTimeout)
		expectSend(serial, []string{menu}, []string{menuSelector})
		//fmt.Println()
	}

	serial.SetReadTimeout(eraseTimeout)
	expectSend(serial, []string{waitFW}, []string{""})
	//fmt.Printf("\n\n> Uploading %d bytes.\n", len(data))

	serial.SetReadTimeout(uploadTimeout)
	//counter := 0
	//previousBlock := -1
	// spinner3, start upload
	progChan <- 3
	callback := func(block int) {
		//if counter%40 == 0 {
		//	if counter != 0 {
		//		fmt.Print("\n")
		//	}
		//	fmt.Print("> ")
		//}
		//marker := "."
		//if block != previousBlock+1 {
		//	marker = "R"
		//}
		//fmt.Print(marker)
		//progChan <- 0
		//counter++
		//previousBlock = block
	}
	err := xmodem.ModemSend1K(serial, data, callback)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println("\n> Upload complete.")

	serial.SetReadTimeout(cleanupTimeout)
	readString(serial)
	// spinner 4,done
	progChan <- 4
}
