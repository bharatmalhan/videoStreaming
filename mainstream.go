package main

import (
	"fmt"
	"crypto/tls"

	"gocv.io/x/gocv"
	quic "github.com/lucas-clemente/quic-go"
)

const addr = "0.0.0.0:8080"

func main() {
	quicConfig := &quic.Config{
		CreatePaths: true,
	}

	sess, err := quic.DialAddr(addr, &tls.Config{InsecureSkipVerify: true}, quicConfig)
	if err != nil {
		fmt.Println("Error Connecting")
		panic(err)
	}

	stream, err := sess.OpenStream()

    deviceID := 0
	webcam, err := gocv.OpenVideoCapture(deviceID)
 	// img := gocv.IMRead("/home/bharat/Pictures/index.jpeg", gocv.IMReadColor)
 	// defer img.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer webcam.Close()
	window := gocv.NewWindow("Face Detect")
	defer window.Close()

	img := gocv.NewMat()
	defer img.Close()
	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("cannot read device %v\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

		// fmt.Println(img.Rows(), img.Cols())
		window.IMShow(img)
		window.WaitKey(1)
		buffer := img.ToBytes()
		stream.Write(buffer)
	}
	// 
	// window := gocv.NewWindow("Imagestream")
	// defer window.Close()
	// window.IMShow(img)
	// window.WaitKey(1)
	// fmt.Println("I am here")
	// buffer := img.ToBytes()
	// stream.Write(buffer)
}
