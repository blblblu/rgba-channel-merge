package cmd

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg" // allow use of jpegs
	"image/png"    // allow use of pngs
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type inputImage struct {
	filepath    string
	channelMask string // something like "rxxb" would mean that the first channel will be used as red channel, and the alpha channel will be used als blue channel
	rgba        *image.RGBA
}

type inputImages []inputImage

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "rgba-channel-merge <image-filepath> <image-channel-mask> [(<image-filepath> <image-channel-mask>)...] <output-png-file-path>",
	Short: "image channel merge tool",
	Long: `A tool to merge specific color channels of multiple images into one rgba image.

The channel masks for each image should match the regex [rgbax]{4}, with r, g, b, and a representing the red, green, blue and alpha channel, and x meaning that this channel should be ignored. E.g., the channel mask "rbax" would mean that the first (red) channel of the input image will be used as the red channel for the output image, the second (green) will be used as the blue channel, the third (blue) will be used as alpha channel, and the fourth (alpha) will be ignored.`,
	Run: func(cmd *cobra.Command, args []string) {
		imgs, outputPath, err := parseArgs(args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		maxSize, err := imgs.openImages()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(2)
		}

		outImg := image.NewNRGBA(image.Rectangle{Min: image.Point{0, 0}, Max: maxSize})
		draw.Draw(outImg, outImg.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 255}}, image.ZP, draw.Src)

		imgs.mergeChannels(outImg)

		outFile, err := os.Create(outputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(3)
		}
		defer outFile.Close()

		png.Encode(outFile, outImg)
	},
}

func parseArgs(args []string) (imgs inputImages, outputPath string, err error) {
	if len(args) < 3 || len(args)%2 != 1 {
		err = fmt.Errorf("wrong input format")
		return
	}

	for i := 0; i < (len(args)-1)/2; i++ {
		channelMask := args[i*2+1]
		if len(channelMask) != 4 {
			err = fmt.Errorf("channel mask \"%s\" has wrong length", channelMask)
			return
		}

		for _, c := range channelMask {
			if !(c == 'r' || c == 'g' || c == 'b' || c == 'a' || c == 'x') {
				err = fmt.Errorf("illegal character \"%c\" in channel mask \"%s\"", c, channelMask)
				return
			}
		}

		imgs = append(imgs, inputImage{
			filepath:    args[i*2],
			channelMask: channelMask,
		})
	}

	outputPath = args[len(args)-1]

	if filepath.Ext(outputPath) != ".png" {
		err = fmt.Errorf("only .png files are supported as output files")
	}

	return
}

func (img *inputImage) openImage() error {
	inFile, err := os.Open(img.filepath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	inImg, _, err := image.Decode(inFile)
	if err != nil {
		return err
	}

	img.rgba = image.NewRGBA(inImg.Bounds())
	draw.Draw(img.rgba, img.rgba.Bounds(), inImg, image.Point{0, 0}, draw.Src)

	return nil
}

func (images inputImages) openImages() (maxSize image.Point, err error) {
	for i := range images {
		if err = images[i].openImage(); err != nil {
			return
		}
		imgSize := images[i].rgba.Bounds().Size()
		if imgSize.X > maxSize.X {
			maxSize.X = imgSize.X
		}
		if imgSize.Y > maxSize.Y {
			maxSize.Y = imgSize.Y
		}
	}
	return
}

func (img inputImage) mergeChannels(outImg *image.NRGBA) {
	if len(img.rgba.Pix) > len(outImg.Pix) {
		// should not happen because the maximum size of all input images will be used for the output image
		panic(fmt.Sprintf("input image is bigger than output image: input: %v, output: %v", img.rgba.Bounds(), outImg.Bounds()))
	}

	for i := 0; i < len(img.rgba.Pix)/4; i++ {
		for j := 0; j < 4; j++ {
			switch img.channelMask[j] {
			case 'r':
				outImg.Pix[i*4+0] = img.rgba.Pix[i*4+j]
			case 'g':
				outImg.Pix[i*4+1] = img.rgba.Pix[i*4+j]
			case 'b':
				outImg.Pix[i*4+2] = img.rgba.Pix[i*4+j]
			case 'a':
				outImg.Pix[i*4+3] = img.rgba.Pix[i*4+j]
			case 'x':
				// do nothing
			default:
				// should not happen, because parseArgs should already take care of that
				panic(fmt.Sprintf("invalid character in channel mask: %c", img.channelMask[j]))
			}
		}
	}
}

func (images inputImages) mergeChannels(outImg *image.NRGBA) {
	for _, img := range images {
		img.mergeChannels(outImg)
	}
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
}
