package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/sydneyowl/g90toolkit/tools"
	"os"
)

// patchimgCmd represents the patchimg command
var patchimgCmd = &cobra.Command{
	Use:   "patchimg",
	Short: "Patch the default boot logo",
	Long:  `Patch the default boot logo`,
	Run: func(cmd *cobra.Command, args []string) {
		data, err := os.ReadFile(FirmwarePath)
		if err != nil {
			fmt.Printf("failed to read firmware: %v", err)
			return
		}
		if Output == "" {
			fmt.Println("Please provide a output path!")
			return
		}
		if LogoPath == "" {
			fmt.Println("Please specify a logo path using --logo-path!")
			return
		}
		logoData, err := os.ReadFile(LogoPath)
		if err != nil {
			fmt.Printf("failed to read logo: %v\n", err)
			return
		}
		if !tools.CheckDecrypted(data) {
			fmt.Println("THE FIRMWARE HAS NOT BEEN DECRYPTED. PLEASE DECRYPT THE FIRMWARE FIRST!")
			return
		}

		if err := tools.PatchBootLogo(logoData, data); err != nil {
			fmt.Printf("failed to patch boot logo: %v\n", err)
			return
		}
		err = os.WriteFile(Output, data, 0777)
		if err != nil {
			fmt.Printf("failed to patch image: %v\n", err)
			return
		}
		fmt.Println("Successfully patched boot logo")
	},
}

func init() {
	rootCmd.AddCommand(patchimgCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// patchimgCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	patchimgCmd.Flags().StringVar(&FirmwarePath, "firmware", "", "Specify a decrypted firmware to path.")
	patchimgCmd.Flags().StringVar(&LogoPath, "logo-path", "", "Specify the logo path you want to apply to the firmware.")
	patchimgCmd.Flags().StringVar(&Output, "output", "", "Specify a path to save patched firmware.")

	patchimgCmd.MarkFlagRequired("firmware")
	patchimgCmd.MarkFlagRequired("logo-path")
	patchimgCmd.MarkFlagRequired("output")
}
