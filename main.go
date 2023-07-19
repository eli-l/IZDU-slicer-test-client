package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
	"os"
	"time"
)

/*
TODO: Move logic to separate functions
*/
func main() {
	service_url := os.Args[1]
	image_url := os.Args[2]

	if service_url == "" || image_url == "" {
		panic("Please pass service URL as first arg and image URL as second arg")
	}

	start := time.Now()
	res, err := http.Post(service_url,
		"application/json",
		bytes.NewBuffer([]byte(`{"image_url": "`+image_url+`"}`)))
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	fmt.Println(fmt.Sprintf("Response time: %s", time.Since(start)))

	var rcvd_total int
	var response = make([]byte, 1200*1200*3*4)
	for {
		b := make([]byte, 1200*1200*3)
		rcvd_chunk, err := res.Body.Read(b)
		rcvd_total += rcvd_chunk
		if err == io.EOF {
			break
		}
		if err != nil && err != io.EOF {
			panic(err)
		}
		response = append(response, b[:rcvd_chunk]...)
	}
	fmt.Println(fmt.Sprintf("Received bytes: %d", rcvd_total))

	var count = 0
	for {
		index := bytes.LastIndex(response, []byte{0x89, 0x50, 0x4E, 0x47})
		fmt.Println(fmt.Sprintf("Index: %d", index))
		if index < 0 {
			break
		}

		if count > 3 {
			fmt.Println("Fuse!")
			break
		}

		rdr := bytes.NewReader(response[index:])
		img, format, err := image.Decode(rdr)
		_ = format
		if err != nil {
			fmt.Println(err)
			break
		}

		out, err := os.Create(fmt.Sprintf("%d.png", count))
		if err != nil {
			panic(err)
		}

		if err := png.Encode(out, img); err != nil {
			panic(err)
		}
		response = response[:index]
		//
		count++
	}

	fmt.Println(fmt.Sprintf("Received chunks: %d", count))

	fmt.Println(fmt.Sprintf("Total time: %s", time.Since(start)))
}
