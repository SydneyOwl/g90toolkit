package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/sydneyowl/g90toolkit/tools"
	"os"
	"regexp"
)

// patchtextCmd represents the patchtext command
var patchtextCmd = &cobra.Command{
	Use:   "patchtext",
	Short: "Patch the default boot text",
	Long:  `Patch the default boot text`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Note: you can also change this by long-press vm key on your rig!")
		data, err := os.ReadFile(FirmwarePath)
		if err != nil {
			fmt.Printf("failed to read firmware: %v", err)
			return
		}
		if Output == "" {
			fmt.Println("Please provide a output path!")
			return
		}
		if Text == "" {
			fmt.Println("Please specify text to be patched using --text!")
			return
		}
		if !tools.CheckDecrypted(data) {
			fmt.Println("THE FIRMWARE HAS NOT BEEN DECRYPTED. PLEASE DECRYPT THE FIRMWARE FIRST!")
			return
		}

		re := regexp.MustCompile("^[a-zA-Z0-9]+$")
		if !re.MatchString(Text) {
			fmt.Println("Only numbers or letters are allowed.")
			return
		}

		if len(Text) > 6 {
			fmt.Println("Only six-digit numbers or letters are allowed")
			return
		}

		tmp := []byte(Text)
		textData := make([]byte, 8)
		copy(textData, tmp)
		if err := tools.PatchBootText(textData, data); err != nil {
			fmt.Printf("failed to patch boot text: %v\n", err)
			return
		}
		err = os.WriteFile(Output, data, 0777)
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
	patchtextCmd.Flags().StringVar(&Text, "text", "", "Specify the text you want to apply to the firmware.")
	patchtextCmd.Flags().StringVar(&Output, "output", "", "Specify a path to save patched firmware.")
	patchtextCmd.Flags().StringVar(&FirmwarePath, "firmware", "", "Specify a decrypted firmware to path.")

	patchtextCmd.MarkFlagRequired("text")
	patchtextCmd.MarkFlagRequired("output")
	patchtextCmd.MarkFlagRequired("firmware")
}
