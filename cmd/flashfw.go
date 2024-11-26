/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/Kei-K23/spinix"
	"github.com/spf13/cobra"
	"github.com/sydneyowl/g90toolkit/lib/g90updatefw"
	"github.com/sydneyowl/g90toolkit/tools"
	"os"
	"strings"
	"time"
)

var deviceFile string

// flashfwCmd represents the flashfw command
var flashfwCmd = &cobra.Command{
	Use:   "flashfw",
	Short: "Flash firmware into the device",
	Long: `Same as g90updatefw:
This program is designed to write a firmware file to a Xiegu radio.
It can be used to update either the main unit or the display unit.

    Usage: ./g90toolkit --device <serial_device> --firmware <firmware_file>

where <firmware_file> is the name of a firmware file for either the
main unit or for the display unit and <serial_device> is the name of
the serial port connected to the Xiegu radio.  On non-windows machines
the <serial_device> is typically similar to /dev/ttyUSB2. On windows
machines it will be similar to COMM2.

You should start the program with the programming cable plugged in
and the power disconnected from the radio.`,
	Run: func(cmd *cobra.Command, args []string) {
		serial, err := g90updatefw.SerialOpen(deviceFile, 115200)
		if err != nil {
			fmt.Printf("Error opening device: %v", err)
			return
		}
		defer serial.Close()

		data, err := os.ReadFile(FirmwarePath)
		if err != nil {
			fmt.Printf("Error reading firmware file: %v", err)
			return
		}
		if tools.CheckDecrypted(data) {
			fmt.Println(`THIS IS A DECRYPTED FIRMWARE ANT CANNOT BE FLASHED INTO THE DEVICE.
YOU MUST ENCRYPT IT FIRST THEN FLASH IT OR THE DEVICE WON'T FUNCTION PROPERLY.
ARE YOU SURE YOU WANT TO CONTINUE? [y/n]`)
			choice := "n"
			_, _ = fmt.Scanln(&choice)
			if strings.ToUpper(choice) != "Y" {
				return
			}
		}
		fmt.Println(`
To flash your device, please:
> 1. Disconnect power cable from the radio.
> 2. Reconnect power cable to the radio.
> 3. Press the volume button and while holding it in,
> 4. Press the power button until the radio begins erasing the existing firmware.
`)
		// start update radio, 0: writing 1:done
		progChan := make(chan uint, 4)

		go g90updatefw.UpdateRadio(serial, data, progChan)
		<-progChan
		fmt.Print("Waiting for device ready...")
		<-progChan
		fmt.Print("\rWaiting for device ready...Done\n")
		fmt.Print("Erasing and waiting for fw...")
		<-progChan
		fmt.Print("\rErasing and waiting for fw...Done\n")
		spinner := spinix.NewSpinner().
			SetSpinnerColor("\033[34m").
			SetMessage("Uploading firmware...").
			SetMessageColor("\033[36m").
			SetSpeed(100 * time.Millisecond).
			SetLastFrame("✔").
			SetLastFrameColor("\033[34m").
			SetLastMessage("Uploading firmware...Done").
			SetLastMessageColor("\033[36m")
		spinner.Start()
		<-progChan
		spinner.Stop()
	},
}

func init() {
	rootCmd.AddCommand(flashfwCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// flashfwCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	flashfwCmd.Flags().StringVar(&deviceFile, "device", "", "name of the serial port connected to the Xiegu radio.")
}
