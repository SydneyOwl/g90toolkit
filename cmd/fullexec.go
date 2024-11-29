package cmd

import (
	"fmt"
	"github.com/Kei-K23/spinix"
	"github.com/spf13/cobra"
	"github.com/sydneyowl/g90toolkit/lib/g90updatefw"
	"github.com/sydneyowl/g90toolkit/tools"
	"os"
	"os/user"
	"regexp"
	"runtime"
	"time"
)

// fullexecCmd represents the fullexec command
var fullexecCmd = &cobra.Command{
	Use:   "fullexec",
	Short: "Exec all procedure for you",
	Long:  `Auto decrypt firmware,patch logo and text, re-encrypt firmware, and flash it into device.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("-> FULLEXEC MODE")
		if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
			currentUser, _ := user.Current()
			if currentUser.Uid != "0" && !NoRootCheck {
				fmt.Println("WARNING: THIS PROGRAM IS NOT RUNNING WITH SUDO. SUDO IS NEEDED WHEN FLASHING FIRMWARE.")
				fmt.Println("SUPPRESS THIS WARNING WITH --no-root-check.")
			}
		}
		fwdata, err := os.ReadFile(FirmwarePath)
		if err != nil {
			fmt.Printf("Failed to read %s: %v", FirmwarePath, err)
			return
		}
		if tools.CheckDecrypted(fwdata) {
			fmt.Println("WARNING: THIS IS A DECRYPTED FIRMWARE. IT IS RECOMMENDED TO USE ORIGINAL ONE.")
		}
		decryptedData, err := tools.DoDecrypt(Key, fwdata)
		if err != nil {
			fmt.Printf("Failed to decrypt %s: %v", FirmwarePath, err)
			return
		}
		fmt.Println("Firmware decrypted.")
		if LogoPath == "" {
			fmt.Println("Boot logo is not provided. Skipping...")
		} else {
			logoData, err := os.ReadFile(LogoPath)
			if err != nil {
				fmt.Printf("failed to read logo: %v\n", err)
				return
			}
			fmt.Println("Trying to patch logo...")
			if err := tools.PatchBootLogo(logoData, decryptedData); err != nil {
				fmt.Printf("failed to patch boot logo: %v\n", err)
				return
			}
			fmt.Println("Done.")
		}

		if Text == "" {
			fmt.Println("Boot text is not provided. Skipping...")
		} else {
			re := regexp.MustCompile("^[a-zA-Z0-9]+$")
			if !re.MatchString(Text) {
				fmt.Println("Only numbers or letters are allowed.")
				return
			}
			if len(Text) > 6 {
				fmt.Println("Only six-digit numbers or letters are allowed.")
				return
			}
			tmp := []byte(Text)
			textData := make([]byte, 8)
			copy(textData, tmp)
			if err := tools.PatchBootText(textData, decryptedData); err != nil {
				fmt.Printf("failed to patch boot text: %v\n", err)
				return
			}
			fmt.Println("Boot text patched successfully.")
		}
		fmt.Println("All steps done. Encrypting firmware...")
		encryptedData, err := tools.DoEncrypt(Key, decryptedData)
		if err != nil {
			fmt.Printf("Failed to encrypt %s: %v", FirmwarePath, err)
			return
		}
		if Output == "" {
			fmt.Println("Output path is not set. Skipping...")
		} else {
			if err := os.WriteFile(Output, encryptedData, 0777); err != nil {
				fmt.Printf("failed to save firmware: %v", err)
				return
			}
			fmt.Println("Firmware saved successfully.")
		}
		if DeviceFile == "" {
			fmt.Println("Device is not set. Skipping...")
		} else {
			fmt.Println("Plug in the cable, then press enter...")
			fmt.Scanln()
			fmt.Println(`
To flash your device, please:
> 1. Disconnect power cable from the radio.
> 2. Reconnect power cable to the radio.
> 3. Press the volume button and while holding it in,
> 4. Press the power button until the radio begins erasing the existing firmware.
`)
			// start update radio, 0: writing 1:done
			progChan := make(chan uint, 4)
			serial, err := g90updatefw.SerialOpen(DeviceFile, 115200)
			if err != nil {
				fmt.Printf("Error opening device: %v", err)
				return
			}
			defer serial.Close()
			go g90updatefw.UpdateRadio(serial, encryptedData, progChan)
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
		}
	},
}

func init() {
	rootCmd.AddCommand(fullexecCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fullexecCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fullexecCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	fullexecCmd.Flags().StringVar(&FirmwarePath, "firmware", "", "Specify a firmware to modify.")
	fullexecCmd.Flags().StringVar(&Key, "key", "", "Specify a key to decrypt/encrypt firmware.")
	fullexecCmd.Flags().StringVar(&DeviceFile, "device", "", "name of the serial port connected to the Xiegu radio. (OPTIONAL)")
	fullexecCmd.Flags().StringVar(&LogoPath, "logo-path", "", "Specify the logo path you want to apply to the firmware. (OPTIONAL)")
	fullexecCmd.Flags().StringVar(&Text, "text", "", "Specify the text you want to apply to the firmware.  (OPTIONAL)")
	fullexecCmd.Flags().BoolVar(&NoRootCheck, "no-root-check", false, "Don't check if the program is running with sudo")
	fullexecCmd.Flags().StringVar(&Output, "output", "", "Specify a path to save patched firmware. (OPTIONAl)")

	fullexecCmd.MarkFlagRequired("firmware")
	fullexecCmd.MarkFlagRequired("key")
}
