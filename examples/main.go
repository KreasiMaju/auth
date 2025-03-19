package main

import (
	"fmt"
	"log"

	"github.com/kreasimaju/auth"
)

func main() {
	// Mendapatkan versi
	version := auth.GetVersion()
	fmt.Printf("Kreasimaju Auth Package Version: %s\n", version)
	log.Println("Package berhasil diintegrasikan!")
}
