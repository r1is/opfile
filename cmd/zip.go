package cmd

import (
	"fmt"
	"log"
	"opfile/common"
	_zip "opfile/module/zip"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var zipCmd = &cobra.Command{
	Use:   "zip",
	Short: "compression module",
	Long: "Compress the target file (directory) into a zip file\n\n" +
		"eg: opfile zip -s /tmp/xxx -d /tmp/xxx.zip (默认排除压缩)\n" +
		"eg: opfile zip -a -s C:\\Users\\xxx\\Desktop\\test -d C:\\Users\\xxx\\Desktop\\test.zip (全部压缩)\n" +
		"eg: opfile zip -s /tmp/xxx -d /tmp/xxx.zip --exclude-file=jpg,png --exclude-dir=image,img (排除jpg、png文件和image、img目录)\n",
	Run: func(cmd *cobra.Command, args []string) {
		if sourceFileName == "" || destFileName == "" {
			cmd.Help()
			fmt.Println("默认排除的后缀: ", ExcludeExtensions)
			fmt.Println("默认排除的目录: ", ExcludeDirectories)
			return
		}
		start := time.Now()
		sourceFileAbsPath, _ := filepath.Abs(sourceFileName)
		destFileAbsPath, _ := filepath.Abs(destFileName)
		if allCompressor {

			if err := _zip.ZipAll(sourceFileAbsPath, destFileAbsPath); err != nil {
				log.Fatalln(err)
			}
			fmt.Printf("%s 压缩成功,存放文件在:%s\n", sourceFileAbsPath, destFileAbsPath)
			end := time.Now()
			msg := fmt.Sprintf("压缩耗时: %v.\n", end.Sub(start))
			fmt.Println(msg)
		} else {
			fmt.Println("默认排除的后缀: ", ExcludeExtensions)
			fmt.Println("默认排除的目录: ", ExcludeDirectories)
			fmt.Println("开始压缩文件...")
			if err := _zip.ZipExclude(sourceFileName, destFileName, ExcludeExtensions, ExcludeDirectories); err != nil {
				log.Fatalln(err)
			}
			fmt.Printf("%s 压缩成功,存放文件在:%s\n", sourceFileAbsPath, destFileAbsPath)
			end := time.Now()
			msg := fmt.Sprintf("压缩耗时: %v.\n", end.Sub(start))
			fmt.Println(msg)
		}
	},
}

// 需要压缩文件的名字
var (
	sourceFileName     string
	destFileName       string
	allCompressor      bool
	ExcludeExtensions  []string
	ExcludeDirectories []string
)

func init() {
	rootCmd.AddCommand(zipCmd)
	zipCmd.Flags().StringVarP(&sourceFileName, "source", "s", "", "源文件名 eg: /tmp/xxx 、C:\\Users\\xxx\\Desktop\\test")
	zipCmd.Flags().StringVarP(&destFileName, "dest", "d", "", "目标文件名 eg: /tmp/xxx.zip 、C:\\Users\\xxx\\Desktop\\test.zip")
	zipCmd.Flags().BoolVarP(&allCompressor, "all", "a", false, "压缩目标下的所有文件")
	zipCmd.Flags().StringArrayVar(&ExcludeExtensions, "exclude-file", common.ExcludeExtensions, "排除的文件后缀名")
	zipCmd.Flags().StringArrayVar(&ExcludeDirectories, "exclude-dir", common.ExcludeDirectories, "排除的目录名")
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
