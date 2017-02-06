package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/disintegration/imaging"
)

var configuration Configuration

func main() {
	config := flag.String("c", "./config.json", "configuration file")
	//interval := flag.Int("i", 500, "Increment (in seconds)")

	configuration = getConfiguration(*config)
	image, err := GetAndConvertFrame()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Image extracted and saved: %v", image)

	val, err := cropImage(image)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%v", val)

}

func GetAndConvertFrame() (string, error) {
	vals := strings.Split(configuration.CaptureFrameCommand, " ")
	out, err := exec.Command(vals[0], vals[1:]...).Output()
	if err != nil {
		return "", err
	}
	log.Printf("%s", out)

	vals = strings.Split(configuration.ConvertFrameCommand, " ")

	out, err = exec.Command(vals[0], vals[1:]...).Output()
	if err != nil {
		return "", err
	}

	log.Printf("%s", out)

	ok, err := exists("/tmp/images")
	if err != nil {
		return "", err
	}
	if !ok {
		os.MkdirAll("/tmp/images", 0777)
	}

	filename := "/tmp/images/" + time.Now().Format(time.RFC3339) + ".png"
	err = os.Rename(configuration.OutputFile, filename)
	if err != nil {
		return "", err
	}

	return filename, nil
}

func cropImage(path string) (string, error) {
	img, err := imaging.Open(path)
	if err != nil {
		return "", err
	}

	rect := img.Bounds()

	x := rect.Dx()
	y := rect.Dy()

	log.Printf("x: %v", x)
	log.Printf("y: %v", y)

	return "", nil
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
