package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/sydneyowl/g90toolkit/tools"
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
		if err := tools.DoDecryptAndSave(Key, FirmwarePath, output); err != nil {
			fmt.Printf("Error: %v", err)
		} else {
			fmt.Printf("Firmware decrypted using specified key successfully.")
		}
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
