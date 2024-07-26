package imgchop

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/mass8326/imgchop/lib/logger"
)

func Crop(wg *sync.WaitGroup, warning *bool, file string) {
	defer wg.Done()

	warn := func(msg string) {
		logger.Logger.Printf("%s (%s)\n", msg, file)
		*warning = true
	}

	img, err := readImage(file)
	if err != nil {
		warn("Unable to read image!")
		return
	}

	type Croppable interface {
		Bounds() image.Rectangle
		SubImage(r image.Rectangle) image.Image
	}
	croppable, ok := img.(Croppable)
	if !ok {
		warn("Image does not support cropping!")
		return
	}

	target := croppable.Bounds()
	width, height := target.Dx(), target.Dy()
	switch {
	case width == height:
		warn("Image is already a square!")
		return
	case width&1 == 1:
		warn("Image width is an odd number and the image cannot be squared!")
		return
	case height&1 == 1:
		warn("Image height is an odd number and the image cannot be squared!")
		return
	case width > height:
		offset := (width - height) / 2
		target = image.Rect(offset, 0, width-offset, height)
	default:
		offset := (height - width) / 2
		target = image.Rect(0, offset, width, height-offset)
	}
	result := croppable.SubImage(target)

	dir := filepath.Dir(file)
	basename := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
	err = writeImage(result, filepath.Join(dir, basename+"-imgchop.png"))
	if err != nil {
		warn("Could not save output image!")
	}
}

func readImage(name string) (image.Image, error) {
	fd, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	img, _, err := image.Decode(fd)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func writeImage(img image.Image, name string) error {
	fd, err := os.Create(name)
	if err != nil {
		return err
	}
	defer fd.Close()

	return png.Encode(fd, img)
}
