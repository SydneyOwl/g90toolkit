/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/sydneyowl/g90toolkit/firmware_data"
	"github.com/sydneyowl/g90toolkit/tools"
	"os"
	"strings"

	"github.com/spf13/cobra"
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
		md5sum := tools.CalcMD5([]byte(Key))
		if md5sum != firmware_data.KnownKeyMD5 {
			fmt.Println("WARNING: THE KEY PROVIDED MISMATCH WITH KNOWN KEY. USE AT YOUR OWN RISK.")
		}
		data, err := os.ReadFile(FirmwarePath)
		if err != nil {
			fmt.Printf("failed to read firmware: %v\n", err)
			return
		}
		if !tools.CheckDecrypted(data) {
			fmt.Println("Seems like the firmware you provided is already encrypted. Do you wish to continue? [y/n]")
			choice := "n"
			_, _ = fmt.Scanln(&choice)
			if strings.ToUpper(choice) != "Y" {
				return
			}
		}
		dec, err := tools.DoEncrypt(Key, data)
		if err != nil {
			fmt.Printf("failed to Encrypt firmware: %v\n", err)
			return
		}
		err = os.WriteFile(output, dec, 0777)
		if err != nil {
			fmt.Printf("failed to write to output: %v\n", err)
			return
		}
		fmt.Println("Firmware encrypted with the key provided successfully!")
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
