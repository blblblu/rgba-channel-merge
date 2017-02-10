package cmd

import (
	"fmt"
	"image"
	_ "image/jpeg" // allow use of jpegs
	_ "image/png"  // allow use of pngs
	"os"

	"github.com/spf13/cobra"
)

type inputImage struct {
	filepath    string
	channelMask string // something like "rxxb" would mean that the first channel will be used as red channel, and the alpha channel will be used als blue channel
	image       image.Image
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "rgba-channel-merge <image-filepath> <image-channel-mask> [(<image-filepath> <image-channel-mask>)...] <output-file-path>",
	Short: "image channel merge tool",
	Long: `A tool to merge specific color channels of multiple images into one rgba image.

The channel masks for each image should match the regex [rgbax]{4}, with r, g, b, and a representing the red, green, blue and alpha channel, and x meaning that this channel should be ignored. E.g., the channel mask "rbax" would mean that the first (red) channel of the input image will be used as the red channel for the output image, the second (green) will be used as the blue channel, the third (blue) will be used as alpha channel, and the fourth (alpha) will be ignored.`,
	Run: func(cmd *cobra.Command, args []string) {
		inputImages, outputImagePath, err := parseArgs(args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		fmt.Printf("%#v %#v", inputImages, outputImagePath)
	},
}

func parseArgs(args []string) (inputImages []inputImage, outputImagePath string, err error) {
	if len(args) < 3 || len(args)%2 != 1 {
		err = fmt.Errorf("wrong input format")
		return
	}

	for i := 0; i < (len(args)-1)/2; i++ {
		inputImages = append(inputImages, inputImage{
			filepath:    args[i*2],
			channelMask: args[i*2+1],
		})
	}

	outputImagePath = args[len(args)-1]
	return
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
