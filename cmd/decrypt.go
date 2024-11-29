package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/sydneyowl/g90toolkit/tools"
)

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
		if Output == "" {
			fmt.Println("Please provide a output path!")
			return
		}
		if err := tools.DoDecryptAndSave(Key, FirmwarePath, Output); err != nil {
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
	decryptCmd.Flags().StringVar(&FirmwarePath, "firmware", "", "Specify a firmware to decrypt.")
	decryptCmd.Flags().StringVar(&Key, "key", "", "Specify a key to decrypt/encrypt firmware.")
	decryptCmd.Flags().StringVar(&Output, "output", "", "Specify a path to save decrypted firmware.")

	decryptCmd.MarkFlagRequired("firmware")
	decryptCmd.MarkFlagRequired("key")
	decryptCmd.MarkFlagRequired("output")
}
