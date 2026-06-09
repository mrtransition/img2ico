package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"img2ico/internal/converter"
	"img2ico/internal/sizes"
	"img2ico/internal/utils"

	"github.com/spf13/cobra"
)

var (
	output    string
	sizeList  []int
	useAll    bool
	verbose   bool
	overwrite bool
)

var rootCmd = &cobra.Command{
	Use:   "img2ico [flags] <input> [input2 ...]",
	Short: "Convert PNG, JPEG, GIF, BMP images to Windows ICO icon files",
	Long: `img2ico converts one or more images to Windows ICO format.
    Supports PNG, JPEG, GIF, BMP input formats.
    Can generate multiple icon sizes in a single ICO file.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// 处理 --all 覆盖 --sizes
		if useAll {
			sizeList = sizes.Standard
		}

		// 验证尺寸
		if err := sizes.Validate(sizeList); err != nil {
			utils.FatalError("Invalid sizes: %v", err)
		}

		// 展开通配符和文件列表
		var inputFiles []string
		for _, arg := range args {
			files, err := utils.ExpandWildcard(arg)
			if err != nil {
				utils.FatalError("Failed to expand pattern %s: %v", arg, err)
			}
			inputFiles = append(inputFiles, files...)
		}

		if len(inputFiles) == 0 {
			utils.FatalError("No input files found")
		}

		if verbose {
			utils.EnableVerbose()
			fmt.Printf("Input files: %v\n", inputFiles)
			fmt.Printf("Sizes: %v\n", sizeList)
		}

		// 确定输出目录或文件
		var outputDir string
		var singleOutputFile string

		if len(inputFiles) == 1 && output != "" && !isDirPath(output) {
			// 单文件且输出不是目录 -> 直接输出到该文件
			singleOutputFile = output
			outputDir = filepath.Dir(output)
		} else {
			// 多文件或输出是目录
			outputDir = output
			if outputDir == "" {
				outputDir = "."
			}
			// 确保输出目录存在
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				utils.FatalError("Cannot create output directory: %v", err)
			}
		}

		// 创建转换器
		opts := converter.ConvertOptions{
			Sizes:      sizeList,
			Overwrite:  overwrite,
			Verbose:    verbose,
			OutputDir:  outputDir,
			SingleFile: singleOutputFile,
		}
		c := converter.NewConverter(opts)

		// 执行转换
		if err := c.ConvertAll(inputFiles); err != nil {
			utils.FatalError("Conversion failed: %v", err)
		}

		if verbose {
			fmt.Println("All conversions completed successfully.")
		}
	},
}

func init() {
	rootCmd.Flags().StringVarP(&output, "output", "o", "", "Output file (single input) or directory (multiple inputs)")
	rootCmd.Flags().IntSliceVarP(&sizeList, "sizes", "s", []int{16, 32, 48, 256}, "Icon sizes (comma-separated, e.g. 16,32,48,256)")
	rootCmd.Flags().BoolVarP(&useAll, "all", "m", false, "Generate all standard sizes (16,24,32,48,64,72,96,128,256)")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	rootCmd.Flags().BoolVarP(&overwrite, "overwrite", "w", false, "Overwrite existing output files")
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// isDirPath checks if a path exists and is a directory, or if it ends with a separator.
func isDirPath(p string) bool {
	info, err := os.Stat(p)
	if err == nil {
		return info.IsDir()
	}
	// If path doesn't exist, treat as file if it has an extension, else directory
	return !strings.Contains(filepath.Base(p), ".")
}
