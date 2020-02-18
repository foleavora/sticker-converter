package main

import (
	"image"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"

	"github.com/nfnt/resize"
)

func main() {
	//open the picture
	file, err := os.Open("test.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	//decode the image
	pic, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	//get size of the picture
	xsize := pic.Bounds().Max.X
	ysize := pic.Bounds().Max.Y

	//scale the image to make the bigger side 512px
	if xsize > ysize {
		pic = resize.Resize(512, 0, pic, resize.Bicubic)
	} else {
		pic = resize.Resize(0, 512, pic, resize.Bicubic)
	}

	//create the new file
	newfile, err := os.Create("test.png")
	if err != nil {
		log.Fatal(err)
	}
	defer newfile.Close()

	//write the new pic into said file
	err = png.Encode(newfile, pic)
	if err != nil {
		log.Fatal(err)
	}

}
