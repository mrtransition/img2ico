package converter

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"img2ico/internal/utils"

	"github.com/disintegration/imaging"
)

type ConvertOptions struct {
	Sizes      []int
	Overwrite  bool
	Verbose    bool
	OutputDir  string
	SingleFile string
}

type Converter struct {
	opts ConvertOptions
}

// imageEntry 用于存储每个尺寸的 PNG 数据
type imageEntry struct {
	Size    int
	PNGData []byte
	Width   int
	Height  int
}

func NewConverter(opts ConvertOptions) *Converter {
	return &Converter{opts: opts}
}

func (c *Converter) ConvertAll(inputs []string) error {
	for _, inputPath := range inputs {
		if err := c.convertOne(inputPath); err != nil {
			return fmt.Errorf("failed to convert %s: %w", inputPath, err)
		}
	}
	return nil
}

func (c *Converter) convertOne(inputPath string) error {
	outputPath := c.getOutputPath(inputPath)

	if !c.opts.Overwrite {
		if _, err := os.Stat(outputPath); err == nil {
			if c.opts.Verbose {
				fmt.Printf("Skipping %s: output file already exists (use --overwrite to replace)\n", inputPath)
			}
			return nil
		}
	}

	img, err := utils.DecodeImage(inputPath)
	if err != nil {
		return err
	}

	if c.opts.Verbose {
		bounds := img.Bounds()
		fmt.Printf("Converting %s (%dx%d) to %s with sizes %v\n",
			inputPath, bounds.Dx(), bounds.Dy(), outputPath, c.opts.Sizes)
	}

	var entries []imageEntry

	for _, size := range c.opts.Sizes {
		resized := resizeToSquare(img, size)
		rgba := toRGBA(resized)
		var buf bytes.Buffer
		if err := png.Encode(&buf, rgba); err != nil {
			return fmt.Errorf("encode PNG size %d: %w", size, err)
		}
		entries = append(entries, imageEntry{
			Size:    size,
			PNGData: buf.Bytes(),
			Width:   size,
			Height:  size,
		})
		if c.opts.Verbose {
			fmt.Printf("  - size %d encoded (%d bytes)\n", size, len(buf.Bytes()))
		}
	}

	if err := encodeICOWithPNG(entries, outputPath); err != nil {
		return err
	}

	if c.opts.Verbose {
		fmt.Printf("Created: %s (contains %d images)\n", outputPath, len(entries))
	}
	return nil
}

func (c *Converter) getOutputPath(inputPath string) string {
	baseName := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	icoName := baseName + ".ico"
	if c.opts.SingleFile != "" {
		return c.opts.SingleFile
	}
	return filepath.Join(c.opts.OutputDir, icoName)
}

// encodeICOWithPNG 手动写入ICO文件，支持多尺寸PNG
func encodeICOWithPNG(entries []imageEntry, outputPath string) error {
	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// ICONDIR header
	if err := writeUint16(f, 0); err != nil { // reserved
		return err
	}
	if err := writeUint16(f, 1); err != nil { // type 1 = ICO
		return err
	}
	if err := writeUint16(f, uint16(len(entries))); err != nil {
		return err
	}

	// 记录每个entry的位置，以便后续填充offset
	entryOffsets := make([]int64, len(entries))
	for i, e := range entries {
		entryOffsets[i], _ = f.Seek(0, io.SeekCurrent)
		widthByte := byte(e.Width)
		heightByte := byte(e.Height)
		if e.Width == 256 {
			widthByte = 0
		}
		if e.Height == 256 {
			heightByte = 0
		}
		if err := writeUint8(f, widthByte); err != nil {
			return err
		}
		if err := writeUint8(f, heightByte); err != nil {
			return err
		}
		if err := writeUint8(f, 0); err != nil { // color count
			return err
		}
		if err := writeUint8(f, 0); err != nil { // reserved
			return err
		}
		if err := writeUint16(f, 1); err != nil { // planes
			return err
		}
		if err := writeUint16(f, 32); err != nil { // bit count
			return err
		}
		if err := writeUint32(f, uint32(len(e.PNGData))); err != nil {
			return err
		}
		if err := writeUint32(f, 0); err != nil { // offset placeholder
			return err
		}
	}

	// 写入图像数据并回填offset
	for i, e := range entries {
		offset, _ := f.Seek(0, io.SeekCurrent)
		// 回填offset
		if _, err := f.Seek(entryOffsets[i]+12, io.SeekStart); err != nil {
			return err
		}
		if err := writeUint32(f, uint32(offset)); err != nil {
			return err
		}
		if _, err := f.Seek(offset, io.SeekStart); err != nil {
			return err
		}
		if _, err := f.Write(e.PNGData); err != nil {
			return err
		}
	}
	return nil
}

// 小端写入辅助函数
func writeUint8(w io.Writer, v byte) error {
	_, err := w.Write([]byte{v})
	return err
}
func writeUint16(w io.Writer, v uint16) error {
	_, err := w.Write([]byte{byte(v), byte(v >> 8)})
	return err
}
func writeUint32(w io.Writer, v uint32) error {
	_, err := w.Write([]byte{byte(v), byte(v >> 8), byte(v >> 16), byte(v >> 24)})
	return err
}

// resizeToSquare 将图片缩放到指定尺寸的正方形
func resizeToSquare(img image.Image, size int) *image.NRGBA {
	bounds := img.Bounds()
	srcW, srcH := bounds.Dx(), bounds.Dy()
	if srcW == srcH {
		return imaging.Resize(img, size, size, imaging.Lanczos)
	}
	var resized *image.NRGBA
	if srcW < srcH {
		resized = imaging.Resize(img, size, 0, imaging.Lanczos)
		return imaging.CropCenter(resized, size, size)
	}
	resized = imaging.Resize(img, 0, size, imaging.Lanczos)
	return imaging.CropCenter(resized, size, size)
}

// toRGBA 将任何图片转换为 *image.RGBA
func toRGBA(img image.Image) *image.RGBA {
	if rgba, ok := img.(*image.RGBA); ok {
		return rgba
	}
	bounds := img.Bounds()
	dst := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			dst.Set(x, y, img.At(x, y))
		}
	}
	return dst
}
