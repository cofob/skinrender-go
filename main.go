package main

import (
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Coords struct {
	X int `yaml:"x"`
	Y int `yaml:"y"`
}

type Map struct {
	Move Coords `yaml:"move"`
	From Coords `yaml:"from"`
	To   Coords `yaml:"to"`
}

type Schema struct {
	Map  []Map  `yaml:"map"`
	Size Coords `yaml:"size"`
}

func main() {
	content, err := ioutil.ReadFile("map.yml") // the file is inside the local directory
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	schema := Schema{}

	err = yaml.Unmarshal(content, &schema)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// This example uses png.Decode which can only decode PNG images.
	catFile, err := os.Open("orig.png")
	if err != nil {
		log.Fatal(err)
	}
	defer catFile.Close()

	// Consider using the general image.Decode as it can sniff and decode any registered image format.
	img, err := png.Decode(catFile)
	if err != nil {
		log.Fatal(err)
	}

	outImg := image.NewRGBA(image.Rect(0, 0, schema.Size.X, schema.Size.Y))

	for _, element := range schema.Map {
		for y := element.From.Y; y < element.To.Y; y++ {
			for x := element.From.X; x < element.To.X; x++ {
				at := img.At(x, y)
				_, _, _, a := at.RGBA()
				if a != 0 {
					outImg.Set(x-element.From.X+element.Move.X, y-element.From.Y+element.Move.Y, at)
				}
			}
		}
	}

	out, _ := os.Create("out.png")
	png.Encode(out, outImg)
	out.Close()
}
