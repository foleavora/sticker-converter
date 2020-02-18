package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"io"
	"log"
	"os"

	"github.com/foobaz/lossypng/lossypng"
	"github.com/nfnt/resize"
)

func compress(input io.Reader, output io.Writer) error {

	//decode the image
	pic, _, err := image.Decode(input)
	if err != nil {
		return err
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

	//create a new buffer to check for filesize of the picture
	buf := new(bytes.Buffer)
	png.Encode(buf, pic)

	//variables for compression level, file size and the new picture
	comp := 1
	filesize := buf.Len()
	var newpic image.Image

	//compress as long as needed to reach desired file size
	for comp <= 20 && filesize > 512000 {
		//Compress image
		newpic = lossypng.Compress(pic, lossypng.NoConversion, comp)

		//write the compressed image into the buffer
		buf.Reset()
		png.Encode(buf, newpic)

		//get new file size and compression level
		filesize = buf.Len()
		comp++

		//stop if the file is still too large with maximum compression
		if comp > 20 {
			return errors.New("Picture too large, compression failed")
		}
	}

	//write the new pic into said file
	_, err = io.Copy(output, buf)
	if err != nil {
		return err
	}

	//output final parameters
	comp--
	fmt.Println("Compression successful with compression level ", comp, " and file size ", filesize)

	return nil
}

func main() {
	//open the picture
	file, err := os.Open("test.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	//create the new file
	newfile, err := os.Create("test.png")
	if err != nil {
		log.Fatal(err)
	}
	defer newfile.Close()

	//compress the file
	err = compress(file, newfile)
	if err != nil {
		log.Fatal(err)
	}
}
