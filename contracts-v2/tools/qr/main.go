// qr prints a text and a PNG QR code for an address (used to fund the deployer).
package main

import (
	"fmt"
	"os"

	qrcode "github.com/skip2/go-qrcode"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("usage: qr <text> <out.png>")
		os.Exit(1)
	}
	q, err := qrcode.New(os.Args[1], qrcode.Medium)
	if err != nil {
		panic(err)
	}
	fmt.Print(q.ToSmallString(false))
	if err := q.WriteFile(512, os.Args[2]); err != nil {
		panic(err)
	}
}
