package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
)

type image struct {
	Label int
	Data  []byte
}

func read4ByteInt(inf *os.File) (int32, error) {
	var output int32
	binary.Read(inf, binary.BigEndian, &output)
	return output, nil
}

func read28x28Image(inf *os.File) []byte {
	var row = make([]byte, 28*28)
	// 28x28 pixels full image read
	n, err := inf.Read(row)
	if n != (28 * 28) {
		fmt.Println("did not read 28x28 bytes!")
	}
	if err != nil {
		fmt.Println(err)
	}
	return row
}

func printImage(b []byte) {
	var count int
	for i := 0; i < 28; i++ {
		for j := 0; j < 28; j++ {
			fmt.Printf("%2.x", b[count])
			count++
		}
		fmt.Println()
	}
}

func readLabelFile(imgs []image) {
	inf, err := os.Open("train-labels.idx1-ubyte")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer inf.Close()
	// read the magic number and verify first
	output, _ := read4ByteInt(inf)
	fmt.Println("output is", output)
	if output == 2049 {
		fmt.Println("valid image format!")
	}

	// read the number of images
	imageCount, _ := read4ByteInt(inf)
	fmt.Println("output", imageCount)
	if imageCount == 60000 {
		fmt.Println("valid number of images for file")
	}

	var buf = make([]byte, 1)
	// now read labels
	for i := 0; i < int(imageCount); i++ {
		inf.Read(buf)
		imgs[i].Label = int(buf[0])
	}
}

func readImageFile(imgs []image) {
	inf, err := os.Open("train-images.idx3-ubyte")
	if err != nil {
		fmt.Println("could not open file!")
		return
	}
	defer inf.Close()

	// read the magic number and verify first
	output, _ := read4ByteInt(inf)
	fmt.Println("output is", output)
	if output == 2051 {
		fmt.Println("valid image format!")
	} else {
		fmt.Println("invalid format for image", output)
		return
	}

	// read the number of images
	imageCount, _ := read4ByteInt(inf)
	fmt.Println("output", imageCount)
	if imageCount == 60000 {
		fmt.Println("valid number of images for file")
	}

	// read num of rows
	output, _ = read4ByteInt(inf)
	fmt.Println("output", output)
	if output == 28 {
		fmt.Println("valid number of rows for file")
	}

	// read num of cols
	output, _ = read4ByteInt(inf)
	fmt.Println("output", output)
	if output == 28 {
		fmt.Println("valid number of cols for file")
	}

	// read all images
	//var bs []byte
	for i := 0; i < int(imageCount); i++ {
		imgs[i].Data = read28x28Image(inf)
	}
	//printImage(bs)

}

func main() {
	var imgs []image = make([]image, 60000)
	readLabelFile(imgs)
	readImageFile(imgs)
	fmt.Println(imgs[1].Label)
	printImage(imgs[1].Data)
	data, _ := json.MarshalIndent(imgs[1], "", " ")
	//fmt.Println(err)
	fmt.Println(string(data))
}
