/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/sydneyowl/g90toolkit/tools"
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt firmware using specified key.",
	Long:  `Encrypt firmware using specified key.`,
	Run: func(cmd *cobra.Command, args []string) {
		if Key == "" {
			fmt.Println("Please provide a key!")
			return
		}
		if output == "" {
			fmt.Println("Please provide a output path!")
			return
		}
		if err := tools.DoEncryptAndSave(Key, FirmwarePath, output); err != nil {
			fmt.Printf("Error: %v", err)
		} else {
			fmt.Printf("Firmware encrypted using specified key successfully.")
		}
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// encryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// encryptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	encryptCmd.Flags().StringVar(&output, "output", "", "Specify a path to save encrypted firmware.")
}
