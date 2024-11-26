package cmd

import (
	"errors"
	"fmt"
	"github.com/Kei-K23/spinix"
	"github.com/spf13/cobra"
	"github.com/sydneyowl/g90toolkit/firmware_data"
	"github.com/sydneyowl/g90toolkit/lib/g90updatefw"
	"github.com/sydneyowl/g90toolkit/tools"
	ser "go.bug.st/serial"
	"os"
	"os/user"
	"regexp"
	"runtime"
	"strings"
	"time"
)

var (
	FirmwarePath string
	Key          string
	interactive  bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "G90ToolKit",
	Short: "A simple application to modify g90 series firmware",
	Long:  `This software allows you to modify the g90 series firmware.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if interactive {
			return nil
		}
		cmd.SilenceUsage = true
		if FirmwarePath == "" {
			return errors.New("firmware path is missing")
		}
		if _, err := os.Stat(FirmwarePath); err != nil {
			return fmt.Errorf("failed to read firmware: %v", err)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if !interactive {
			return
		}
		if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
			currentUser, _ := user.Current()
			if currentUser.Uid != "0" {
				fmt.Println("WARNING: THIS PROGRAM IS NOT RUNNING WITH SUDO. SUDO IS NEEDED WHEN FLASHING FIRMWARE.")
			}
		}
		// run interactive mode
		if FirmwarePath == "" {
			fmt.Print("Input a ORIGINAL firmware path here: ")
			fmt.Scanln(&FirmwarePath)
		}
		fwdata, err := os.ReadFile(FirmwarePath)
		if err != nil {
			fmt.Printf("Failed to read %s: %v", FirmwarePath, err)
			return
		}
		if tools.CheckDecrypted(fwdata) {
			fmt.Println("Seems like you selected a decrypted firmware. Please use original firmware! ")
			return
		}
		if Key == "" {
			// use default key
			Key = firmware_data.KnownKey
		}
		// start decrypt software
		decryptedData, err := tools.DoDecrypt(Key, fwdata)
		if err != nil {
			fmt.Printf("Failed to decrypt %s: %v", FirmwarePath, err)
			return
		}
		fmt.Println("Firmware decrypted.")
		fmt.Print("Do you wish to patch a new boot logo? [y/n] ")
		choice := "n"
		fmt.Scanln(&choice)
		if strings.ToUpper(choice) == "Y" {
			fmt.Print("Input the logo path you want to patch at here, should be a png file: ")
			logopath := ""
			fmt.Scanln(&logopath)
			logoData, err := os.ReadFile(logopath)
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
		fmt.Print("Do you wish to patch a new boot text? [y/n]")
		fmt.Scanln(&choice)
		if strings.ToUpper(choice) == "Y" {
			fmt.Print("Note: Only six-digit numbers or letters are allowed. \nInput the text you want to patch at here: ")
			text := ""
			fmt.Scanln(&text)
			re := regexp.MustCompile("^[a-zA-Z0-9]+$")
			if !re.MatchString(text) {
				fmt.Println("Only numbers or letters are allowed.")
				return
			}
			if len(text) > 6 {
				fmt.Println("Only six-digit numbers or letters are allowed.")
				return
			}
			tmp := []byte(text)
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
		fmt.Print("Do you wish to save the patched firmware so that you can flash it directly next time? [y/n] ")
		fmt.Scanln(&choice)
		if strings.ToUpper(choice) == "Y" {
			output := ""
			fmt.Print("Where do you want to save the firmware? e.g. /home/user/desktop/firmware.xgf:  ")
			fmt.Scanln(&output)
			if err := os.WriteFile(output, encryptedData, 0777); err != nil {
				fmt.Printf("failed to save firmware: %v", err)
				return
			}
			fmt.Println("Firmware saved successfully.")
		}
		fmt.Print("Do you wish to flash the patched firmware into your rig now? [y/n]")
		fmt.Scanln(&choice)
		if strings.ToUpper(choice) != "Y" {
			fmt.Println("Done.")
			return
		}
		fmt.Println("Plug in the cable, then press enter...")
		ports, _ := ser.GetPortsList()
		fmt.Println("Current available ports:")
		for _, port := range ports {
			fmt.Printf("%s ", port)
		}
		fmt.Println()
		fmt.Println("Input a serial port: ")
		serport := ""
		fmt.Scanln(&serport)
		fmt.Println(`
To flash your device, please:
> 1. Disconnect power cable from the radio.
> 2. Reconnect power cable to the radio.
> 3. Press the volume button and while holding it in,
> 4. Press the power button until the radio begins erasing the existing firmware.
`)
		// start update radio, 0: writing 1:done
		progChan := make(chan uint, 4)
		serial, err := g90updatefw.SerialOpen(serport, 115200)
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
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&FirmwarePath, "firmware", "", "Specify a firmware to modify.")
	rootCmd.PersistentFlags().StringVar(&Key, "key", "", "Specify a key to decrypt/encrypt firmware.")
	rootCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Use interactive mode.")
}
