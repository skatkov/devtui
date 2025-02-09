package main

import (
	"fmt"
	"log"

	"github.com/google/uuid"
)

func main() {
	// version 1 uuid
	v1, err := uuid.NewUUID()
	if err != nil {
		log.Fatal("cannot generate v1 uuid")
	}
	fmt.Printf("v1 uuid: %v\n", v1)

	// version 2 uuid
	v2, err := uuid.NewDCEGroup()
	if err != nil {
		log.Fatal("cannot generate v2 uuid")
	}
	fmt.Printf("v2 uuid: %v\n", v2)

	// version 3 uuid
	v3 := uuid.NewMD5(uuid.NameSpaceURL, []byte("https://example.com"))
	fmt.Printf("v3 uuid: %v\n", v3)

	// version 4 uuid
	v4, err := uuid.NewRandom()
	if err != nil {
		log.Fatal("cannot generate v4 uuid")
	}
	fmt.Printf("v4 uuid: %v\n", v4)

	// version 5 uuid
	v5 := uuid.NewSHA1(uuid.NameSpaceURL, []byte("https://example.com"))
	fmt.Printf("v5 uuid: %v\n", v5)

	// version 6 uuid
	v6, err := uuid.NewV6()
	if err != nil {
		log.Fatal("cannot generate v6 uuid")
	}
	fmt.Printf("v6 uuid: %v\n", v6)

	// version 7 uuid
	v7, err := uuid.NewV7()
	if err != nil {
		log.Fatal("cannot generate v7 uuid")
	}
	fmt.Printf("v7 uuid: %v\n", v7)
}
