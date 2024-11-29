package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	FirmwarePath string
	Key          string
	DeviceFile   string
	LogoPath     string
	Text         string
	Output       string
	NoRootCheck  bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "G90ToolKit",
	Short: "A simple application to modify g90 series firmware",
	Long:  `This software allows you to modify the g90 series firmware.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if FirmwarePath == "" {
			return errors.New("firmware path is missing")
		}
		if _, err := os.Stat(FirmwarePath); err != nil {
			return fmt.Errorf("failed to read firmware: %v", err)
		}
		return nil
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

}
