package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/sydneyowl/g90toolkit/firmware_data"
	"github.com/sydneyowl/g90toolkit/tools"
	"os"
	"strings"
)

var output string

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt firmware using specified key.",
	Long:  `Decrypt firmware using specified key.`,
	Run: func(cmd *cobra.Command, args []string) {
		if Key == "" {
			fmt.Println("Please provide a key!")
			return
		}
		if output == "" {
			fmt.Println("Please provide a output path!")
			return
		}
		md5sum := tools.CalcMD5([]byte(Key))
		if md5sum != firmware_data.KnownKeyMD5 {
			fmt.Println("WARNING: THE KEY PROVIDED MISMATCH WITH KNOWN KEY. USE AT YOUR OWN RISK.")
		}
		data, err := os.ReadFile(FirmwarePath)
		if err != nil {
			fmt.Printf("failed to read firmware: %v\n", err)
			return
		}
		if tools.CheckDecrypted(data) {
			fmt.Println("Seems like the firmware you provided is already decrypted. Do you wish to continue? [y/n]")
			choice := "n"
			_, _ = fmt.Scanln(&choice)
			if strings.ToUpper(choice) != "Y" {
				return
			}
		}
		dec, err := tools.DoDecrypt(Key, data)
		if err != nil {
			fmt.Printf("failed to decrypt firmware: %v\n", err)
			return
		}
		err = os.WriteFile(output, dec, 0777)
		if err != nil {
			fmt.Printf("failed to write to output: %v\n", err)
			return
		}
		fmt.Println("Firmware decrypted with the key provided successfully!")
	},
}

func init() {
	rootCmd.AddCommand(decryptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// decryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	decryptCmd.Flags().StringVar(&output, "output", "", "Specify a path to save decrypted firmware.")
}
