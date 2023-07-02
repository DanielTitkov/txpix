package pix

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/DanielTitkov/txpix/config"
	"github.com/disintegration/imaging"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"
)

func addLabel(img *image.RGBA, face font.Face, col color.RGBA, lines *[][]string, config config.Config) (remainingLines [][]string) {
	drawer := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: face,
	}

	lineHeight := int(face.Metrics().Height.Ceil())  // calculate the height of a line
	maxWidth := img.Bounds().Dx() - 2*config.Margin  // leave some margin on the sides of the image
	maxHeight := img.Bounds().Dy() - 2*config.Margin // leave some margin at the top and bottom of the image
	currentX := config.Margin
	currentY := config.Margin + lineHeight

	remainingWordCount := 0
	for _, line := range *lines {
		remainingWordCount += len(line)
	}

	for i, line := range *lines {
		for j, word := range line {
			wordWidth := drawer.MeasureString(word).Ceil()
			if currentX+wordWidth > maxWidth {
				// if the next word does not fit on the current line,
				// start a new line.
				currentX = config.Margin
				currentY += lineHeight + config.LineSpacing // add line spacing to the current Y position
			}

			if currentY+lineHeight > maxHeight {
				// if the next line does not fit in the image,
				// and there are more than 3 words remaining, return the remaining words.
				if remainingWordCount > 3 {
					remainingLines = make([][]string, len(*lines)-i)
					copy(remainingLines, (*lines)[i:])
					remainingLines[0] = remainingLines[0][j:]
					return remainingLines
				}
			}

			drawer.Dot = fixed.Point26_6{fixed.Int26_6(currentX * 64), fixed.Int26_6(currentY * 64)}
			drawer.DrawString(word + " ") // add a space after each word to separate words
			currentX += wordWidth + drawer.MeasureString(" ").Ceil()
			remainingWordCount--
		}
		// Force a new line after each line of text, unless it is the last line of this image
		if i < len(*lines)-1 || remainingWordCount > 0 {
			currentX = config.Margin
			currentY += lineHeight + config.LineSpacing // add line spacing to the current Y position
		}
	}
	return [][]string{}
}

func preprocessText(data []byte, config config.Config) ([]byte, error) {
	re, err := regexp.Compile(config.Preprocess.Remove)
	if err != nil {
		return nil, err
	}

	processedData := re.ReplaceAllLiteral(data, []byte{})
	return processedData, nil
}

func Build(data []byte, output string, config config.Config) error {
	data, err := preprocessText(data, config)
	if err != nil {
		return err
	}

	// Split the text into lines, then each line into words
	lines := strings.Split(string(data), "\n")
	words := make([][]string, len(lines))
	for i, line := range lines {
		words[i] = strings.Fields(line)
	}

	// Prepare the background
	var bg *image.RGBA
	if len(config.BackgroundImages) > 0 {
		// Choose random background image from the list
		rand.Seed(time.Now().UnixNano())
		randIndex := rand.Intn(len(config.BackgroundImages))

		bgImageFile, err := os.Open(config.BackgroundImages[randIndex])
		if err != nil {
			return err
		}
		defer bgImageFile.Close()

		bgImage, _, err := image.Decode(bgImageFile)
		if err != nil {
			return err
		}

		// Fill the background image
		bgImage = imaging.Fill(bgImage, config.ImageWidth, config.ImageHeight, imaging.Center, imaging.Lanczos)

		// Prepare the background image
		bg = image.NewRGBA(image.Rect(0, 0, config.ImageWidth, config.ImageHeight))
		draw.Draw(bg, bg.Bounds(), bgImage, image.Point{}, draw.Src)
	} else {
		// Use the background color
		bgColor, err := parseColor(config.BackgroundColor)
		if err != nil {
			return err
		}

		bg = image.NewRGBA(image.Rect(0, 0, config.ImageWidth, config.ImageHeight))
		draw.Draw(bg, bg.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)
	}

	fontColor, err := parseColor(config.FontColor)
	if err != nil {
		return err
	}

	var f *truetype.Font
	if config.FontFile != "" {
		fontBytes, err := ioutil.ReadFile(config.FontFile)
		if err != nil {
			panic(err)
		}
		f, err = truetype.Parse(fontBytes)
		if err != nil {
			panic(err)
		}
	} else {
		f, err = truetype.Parse(goregular.TTF)
		if err != nil {
			panic(err)
		}
	}

	opts := truetype.Options{
		Size: config.FontSize,
	}
	face := truetype.NewFace(f, &opts)

	i := 1
	for len(words) > 0 {
		// Copy the background so as not to modify the original
		img := image.NewRGBA(bg.Bounds())
		draw.Draw(img, img.Bounds(), bg, image.Point{}, draw.Src)

		words = addLabel(img, face, fontColor, &words, config)

		filename := fmt.Sprintf("%s_%02d.png", output, i)
		outputFile, err := os.Create(filename)
		if err != nil {
			panic(err)
		}

		png.Encode(outputFile, img)
		outputFile.Close()
		fmt.Println("Image saved:", filename)
		i++
	}
	return nil
}

func parseColor(s string) (color.RGBA, error) {
	var c color.RGBA
	_, err := fmt.Sscanf(s, "%d,%d,%d,%d", &c.R, &c.G, &c.B, &c.A)
	return c, err
}
