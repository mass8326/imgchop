package imgchop

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/mass8326/imgchop/lib/logger"
	"github.com/mass8326/imgchop/lib/util"
)

func Process(wg *sync.WaitGroup, c chan logger.Message, name string, intelligent bool) {
	defer wg.Done()

	lgr := logger.Logger{
		Messages: c,
		Source:   name,
	}

	fd, err := os.Open(name)
	if err != nil {
		lgr.Warn("Unable to open path!")
		return
	}
	defer fd.Close()

	stat, err := fd.Stat()
	if err != nil {
		lgr.Warn("Unable to stat path!")
		return
	}

	if stat.IsDir() {
		entries, err := fd.ReadDir(0)
		if err != nil {
			lgr.Warn("Unable to read directory!")
			return
		}
		wg.Add(len(entries))
		for _, entry := range entries {
			go Process(wg, c, filepath.Join(name, entry.Name()), true)
		}
	} else {
		wg.Add(1)
		Crop(wg, c, fd, intelligent)
	}
}

func Crop(wg *sync.WaitGroup, c chan logger.Message, fd *os.File, intelligent bool) {
	defer wg.Done()

	fname := fd.Name()
	lgr := logger.Logger{
		Messages: c,
		Source:   fname,
	}

	input, _, err := image.Decode(fd)
	if err != nil {
		lgr.Warn("Unable to decode image!")
		return
	}

	type Croppable interface {
		Bounds() image.Rectangle
		SubImage(r image.Rectangle) image.Image
	}
	croppable, ok := input.(Croppable)
	if !ok {
		lgr.Warn("Image does not support cropping!")
		return
	}

	target := croppable.Bounds()
	width, height := target.Dx(), target.Dy()
	switch {
	case width == height:
		lgr.Warn("Image is already a square!")
		return
	case width&1 == 1:
		lgr.Warn("Image width is an odd number and the image cannot be cropped into a square!")
		return
	case height&1 == 1:
		lgr.Warn("Image height is an odd number and the image cannot be cropped into a square!")
		return
	}

	if intelligent {
		var checks [2]image.Rectangle
		if width > height {
			offset := (width - height) / 2
			checks = [2]image.Rectangle{image.Rect(0, 0, offset, height), image.Rect(width-offset, 0, width, height)}
		} else {
			offset := (height - width) / 2
			checks = [2]image.Rectangle{image.Rect(0, 0, width, offset), image.Rect(0, height-offset, width, height)}
		}

		for _, check := range checks {
			bounds := check.Bounds()
			maximums := uint64(0)
			minimums := uint64(0)
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
					r, g, b := util.ColorToNRGB(input.At(x, y))
					maximum := max(r, g, b)
					minimum := min(r, g, b)
					lum := float32(maximum-minimum) * 100 / 255 / 2
					if lum > 20 {
						lgr.Info(fmt.Sprintf("Image did not pass intelligent filter (%f%% pixel luminosity at [%d, %d])", lum, x, y))
						return
					}

					maximums += uint64(maximum)
					minimums += uint64(minimum)
				}
			}
			avg := float32(maximums-minimums) / 2 / float32(bounds.Dx()) / float32(bounds.Dy())
			lum := avg * 100 / 255
			if lum > 1 {
				lgr.Info(fmt.Sprintf("Image did not pass intelligent filter (%f%% average luminosity)", lum))
				return
			}
		}
	}

	if width > height {
		offset := (width - height) / 2
		target = image.Rect(offset, 0, width-offset, height)
	} else {
		offset := (height - width) / 2
		target = image.Rect(0, offset, width, height-offset)
	}
	result := croppable.SubImage(target)

	dir := filepath.Dir(fname)
	basename := strings.TrimSuffix(filepath.Base(fname), filepath.Ext(fname))
	err = writeImage(result, filepath.Join(dir, basename+"-imgchop.png"))
	if err != nil {
		lgr.Warn("Could not save output image!")
	}
}

func writeImage(img image.Image, name string) error {
	fd, err := os.Create(name)
	if err != nil {
		return err
	}
	defer fd.Close()

	return png.Encode(fd, img)
}
