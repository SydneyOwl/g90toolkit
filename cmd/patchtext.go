package cmd

import (
	"fmt"
	"github.com/sydneyowl/g90toolkit/tools"
	"os"
	"unicode"

	"github.com/spf13/cobra"
)

var text string

// patchtextCmd represents the patchtext command
var patchtextCmd = &cobra.Command{
	Use:   "patchtext",
	Short: "Patch the default boot text",
	Long:  `Patch the default boot text`,
	Run: func(cmd *cobra.Command, args []string) {
		data, err := os.ReadFile(FirmwarePath)
		if err != nil {
			fmt.Printf("failed to read firmware: %v", err)
			return
		}
		if output == "" {
			fmt.Println("Please provide a output path!")
			return
		}
		if text == "" {
			fmt.Println("Please specify text to be patched using --text!")
			return
		}
		if !tools.CheckDecrypted(data) {
			fmt.Println("THE FIRMWARE HAS NOT BEEN DECRYPTED. PLEASE DECRYPT THE FIRMWARE FIRST!")
			return
		}

		for _, r := range text {
			if !unicode.IsDigit(r) && !unicode.IsLetter(r) {
				fmt.Println("Only numbers or letters are allowed.")
				return
			}
		}

		if len(text) > 6 {
			fmt.Println("Only six-digit numbers or letters are allowed")
			return
		}

		tmp := []byte(text)
		textData := make([]byte, 8)
		copy(textData, tmp)
		if err := tools.PatchBootText(textData, data); err != nil {
			fmt.Printf("failed to patch boot text: %v\n", err)
			return
		}
		err = os.WriteFile(output, data, 0777)
		if err != nil {
			fmt.Printf("failed to patch text: %v\n", err)
			return
		}
		fmt.Println("Successfully patched boot text.")
	},
}

func init() {
	rootCmd.AddCommand(patchtextCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// patchtextCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// patchtextCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	patchtextCmd.Flags().StringVar(&text, "text", "", "Specify the text you want to apply to the firmware.")
	patchtextCmd.Flags().StringVar(&output, "output", "", "Specify a path to save patched firmware.")
}
