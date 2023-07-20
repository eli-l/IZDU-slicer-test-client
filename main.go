package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

/*
TODO: Move logic to separate functions
*/
func main() {
	if len(os.Args) < 4 {
		log.Fatal("Pass service url, image url and scale(px) as arguments")
	}

	serviceUrl := os.Args[1]
	imageUrl := os.Args[2]
	scale := os.Args[3]
	scaleInt, err := strconv.ParseInt(scale, 10, 32)
	if err != nil {
		log.Fatal(err)
	}

	responseTime, count := processResponse(serviceUrl, imageUrl, scaleInt)
	log.Println(
		fmt.Sprintf("Summary:\r\n Response: %s\r\n Image count: %d",
			responseTime,
			count))
}

func processResponse(serviceUrl, imageUrl string, scale int64) (time.Duration, int) {
	start := time.Now()
	res, err := http.Post(
		fmt.Sprintf("%s/?scale=%d", serviceUrl, scale),
		"application/json",
		bytes.NewBuffer([]byte(`{"image_url": "`+imageUrl+`"}`)))
	if err != nil {
		panic(err)
	}

	responseTime := time.Since(start)
	defer res.Body.Close()
	var receivedTotal int
	var response = make([]byte, 0)
	for {
		b := make([]byte, 1200*1200*3)
		rcvdChunk, err := res.Body.Read(b)
		response = append(response, b[:rcvdChunk]...)
		receivedTotal += rcvdChunk
		if err == io.EOF {
			break
		}
		if err != nil && err != io.EOF {
			panic(err)
		}
	}
	log.Println(fmt.Sprintf("Received bytes: %d", receivedTotal))

	var processedCount = 0
	for {
		index := bytes.LastIndex(response, []byte{0x89, 0x50, 0x4E, 0x47})
		if index < 0 {
			break
		}
		if processedCount > 10 {
			log.Fatal("Fuse!")
		}

		r := bytes.NewReader(response[index:])
		img, _, err := image.Decode(r)
		if err != nil {
			log.Fatal(err)
		}

		out, err := os.Create(fmt.Sprintf("%d.png", processedCount))
		if err != nil {
			panic(err)
		}

		if err := png.Encode(out, img); err != nil {
			panic(err)
		}
		response = response[:index]
		processedCount++
	}

	if processedCount == 0 {
		log.Printf("No images received. Response: %s\n", response[:receivedTotal])
	}

	return responseTime, processedCount
}
