package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	FirmwarePath string
	Key          string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "G90ToolKit",
	Short: "A simple application to modify g90 series firmware",
	Long:  `This software allows you to modify the g90 series firmware.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		if FirmwarePath == "" {
			return errors.New("firmware path is missing")
		}
		if _, err := os.Stat(FirmwarePath); err != nil {
			return fmt.Errorf("failed to read firmware: %v", err)
		}
		return nil
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//Run: func(cmd *cobra.Command, args []string) {
	//
	//},
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&FirmwarePath, "firmware", "", "Specify a firmware to modify.")
	rootCmd.PersistentFlags().StringVar(&Key, "key", "", "Specify a key to decrypt/encrypt firmware.")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
