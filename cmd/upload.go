package cmd

import (
	"fmt"
	"opfile/module/upload"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var ossCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload module",
	Long:  "To upload files to cloud object storage, you need to enter a verification code",
	Run: func(cmd *cobra.Command, args []string) {
		if passcode == "" || uploadFileName == "" {
			cmd.Help()
			return
		} else {
			targetURL := upload.GetServerURL()
			ossCfg, err := upload.GetOssCfg(targetURL, passcode)
			if err != nil {
				return
			}
			key, err := upload.UploadFile(ossCfg, uploadFileName)
			if err != nil {
				fmt.Println("UploadFile failed:", err)
				return
			}
			//获取预签名URL
			upload.GetPresignedURL(ossCfg, key)
		}
	},
}

var (
	uploadFileName string
	passcode       string
)

func init() {
	rootCmd.AddCommand(ossCmd)
	ossCmd.Flags().StringVarP(&passcode, "pcode", "p", "", "Google TOPT code")
	ossCmd.Flags().StringVarP(&uploadFileName, "file", "f", "", "source file name")
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
