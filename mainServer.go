package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"
	"gocv.io/x/gocv"
	// "strconv"

	quic "github.com/lucas-clemente/quic-go"
)

const addr = "0.0.0.0:8080"

var fakeBuff = make([]byte, 1024*1024*100) // 100 MB

func main() {

	quicConfig := &quic.Config{
		CreatePaths: true,
	}
	listener, err := quic.ListenAddr(addr, generateTLSConfig(), quicConfig)

	if err != nil {
		print("Unable to make socket")
		panic(err)
	}

	window := gocv.NewWindow("Image")
	defer window.Close()

	sess, err := listener.Accept()
	if err != nil {
		print("Couldnt make Session")
	}
	fmt.Print("Connection made with ")
	fmt.Println(sess.RemoteAddr())

	stream1, err := sess.AcceptStream()
	if err != nil {
		fmt.Println("Couldnt make Stream 1")
	}

	buff := make([]byte, 50000)

	for {
		bytesRead, err := stream1.Read(buff)
		img, err:= gocv.NewMatFromBytes(480, 640, gocv.MatTypeCV8UC3, buff)
		if(err!=nil) {
			fmt.Println(err)
		}
		bytesRead++
		defer img.Close()
		window.IMShow(img)
		window.WaitKey(1)
	}

	// img, err := gocv.NewMatFromBytes(480, 640, gocv.MatTypeCV8UC3, buff)
	// fmt.Println(bytesRead)
	// defer img.Close()

	// window.IMShow(img)
	// window.WaitKey(1)
	// fmt.Println("Hello I am here")
}


func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{Certificates: []tls.Certificate{tlsCert}}
}