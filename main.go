package main

import (
	"bufio"
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var logger = log.New(os.Stderr, "", 0)

func main() {
	if len(os.Args) < 2 {
		logger.Println("No files were passed in!")
		pause()
		os.Exit(1)
	}

	warning := false
	var wg sync.WaitGroup
	for _, file := range os.Args[1:] {
		wg.Add(1)
		go crop(&wg, &warning, file)
	}
	wg.Wait()

	if warning {
		pause()
	}
}

func crop(wg *sync.WaitGroup, warning *bool, file string) {
	defer wg.Done()

	warn := func(msg string) {
		logger.Printf("%s (%s)\n", msg, file)
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
	if width == height {
		warn("Image is already a square!")
		return
	} else if width&1 == 1 {
		warn("Image width is an odd number and the image cannot be squared!")
		return
	} else if height&1 == 1 {
		warn("Image height is an odd number and the image cannot be squared!")
		return
	}

	if width > height {
		offset := (width - height) / 2
		target = image.Rect(offset, 0, width-offset, height)
	} else {
		offset := (height - width) / 2
		target = image.Rect(0, offset, width, height-offset)
	}
	result := croppable.SubImage(target)

	dir := filepath.Dir(file)
	basename := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
	err = writeImage(result, filepath.Join(dir, basename+"-crop.png"))
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

func pause() {
	logger.Println("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
