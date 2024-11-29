package cmd

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/sydneyowl/g90toolkit/firmware_data"
	"github.com/sydneyowl/g90toolkit/tools"
	"os"
)

// firmwareInfoCmd represents the firmwareInfo command
var firmwareInfoCmd = &cobra.Command{
	Use:   "getinfo",
	Short: "Read firmware information",
	Long:  `Read firmware information.`,
	Run: func(cmd *cobra.Command, args []string) {
		data, err := os.ReadFile(FirmwarePath)
		if err != nil {
			fmt.Printf("failed to read firmware: %v", err)
			return
		}
		decrypted := tools.CheckDecrypted(data)
		md5sum := tools.CalcMD5(data)
		fmt.Println("Firmware Information")
		fmt.Println("=============================")
		fmt.Printf("Firmware size: %d bytes\n", len(data))
		fmt.Printf("Firmware MD5: %s\n", md5sum)
		fmt.Printf("Decrypted: %t\n", decrypted)
		if Key == "" {
			fmt.Printf("Is the provided key valid: NOT PROVIDED\n")
		} else {
			md5sum = tools.CalcMD5([]byte(Key))
			fmt.Printf("Is the provided key valid: %t\n", md5sum == firmware_data.KnownKeyMD5)
		}
		if decrypted {
			index := bytes.Index(data, firmware_data.OriginalBootText)
			fmt.Printf("Has the boot text changed: %t\n", index == -1)

			index = bytes.Index(data, firmware_data.OriginalBootImage)
			fmt.Printf("Has the boot image changed: %t\n", index == -1)
		}
	},
}

func init() {
	rootCmd.AddCommand(firmwareInfoCmd)
	firmwareInfoCmd.Flags().StringVar(&FirmwarePath, "firmware", "", "Specify a firmware to read.")
	firmwareInfoCmd.Flags().StringVar(&Key, "key", "", "Specify a key to decrypt/encrypt firmware.(optional)")

	firmwareInfoCmd.MarkFlagRequired("firmware")
}
