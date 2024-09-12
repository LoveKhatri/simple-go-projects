package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/LoveKhatri/encryption-go/filecrypt"
	"golang.org/x/term"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(0)
	}

	function := os.Args[1]

	switch function {
	case "help":
		printHelp()
	case "encrypt":
		encryptHandle()
	case "decrypt":
		decryptHandle()
	default:
		fmt.Println("Unknown function: " + function)
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("File Encryption")
	fmt.Println("Usage:")
	fmt.Println("filecrypt help")
	fmt.Println("filecrypt encrypt <file>")
	fmt.Println("filecrypt decrypt <file>")
}

func encryptHandle() {
	if len(os.Args) < 3 {
		fmt.Println("Missing file path")
		os.Exit(0)
	}

	file := os.Args[2]
	if !validateFile(file) {
		panic("File not found")
	}

	password := getPassword()

	fmt.Println("\nEncrypting file...")
	filecrypt.Encrypt(file, password)

	fmt.Println("\nFile encrypted successfully")
}

func decryptHandle() {
	if len(os.Args) < 3 {
		fmt.Println("Missing file path")
		os.Exit(0)
	}

	file := os.Args[2]
	if !validateFile(file) {
		panic("File not found")
	}

	fmt.Println("Enter Password:")
	password, _ := term.ReadPassword(0)

	fmt.Println("\nDecrypting file...")
	filecrypt.Decrypt(file, password)

	fmt.Println("\nFile decrypted successfully")
}

func getPassword() []byte {
	fmt.Println("Enter password:")
	password, _ := term.ReadPassword(0)

	fmt.Println("\nConfirm password:")
	confirm, _ := term.ReadPassword(0)

	if !validatePassword(password, confirm) {
		fmt.Println("\nPasswords do not match")
		return getPassword()
	}

	return password
}

func validatePassword(password []byte, confirm []byte) bool {
	return bytes.Equal(password, confirm)
}

func validateFile(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}
