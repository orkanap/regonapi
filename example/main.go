package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/orkanap/regonapi"
	"github.com/subosito/gotenv"
)

func main() {
	// Load environment from .env file, if exists
	gotenv.Load()
	regonsvc := regonapi.NewClient(context.Background(), os.Getenv("REGON_API_KEY"))

	// Create new session
	err := regonsvc.Login()
	if err != nil {
		log.Fatal(err)
	}
	defer regonsvc.Logout()

	// Search by NIP
	nips := []string{"9999999999", "5261040828"}
	for _, nip := range nips {
		entities, err := regonsvc.SearchByNIP(nip)
		if err != nil {
			if err == regonapi.ErrNoDataFound {
				fmt.Println("data not found for NIP:", nip)
				continue
			}
			log.Fatal(err)
		}
		listEntities(entities)
	}

	// Search by REGON
	regons := []string{"123456789", "000331501"}
	for _, regon := range regons {
		entities, err := regonsvc.SearchByREGON(regon)
		if err != nil {
			if err == regonapi.ErrNoDataFound {
				fmt.Println("data not found for REGON:", regon)
				continue
			}
			log.Fatal(err)
		}
		listEntities(entities)
	}

	// List PKDs
	pkds, err := regonsvc.LegalPersonPKDList("340771731")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("PKD:")
	for _, pkd := range pkds {
		fmt.Println(pkd.Code, pkd.Name)
	}
}

func listEntities(entities []regonapi.Entity) {
	const space = "   "
	for _, entity := range entities {
		fmt.Println("[", entity.Type, "]")
		fmt.Println(space, "NIP:", entity.NIP, "REGON:", entity.REGON)
		fmt.Println(space, entity.Name)
		fmt.Println(space, entity.Street, entity.PropertyNumber, entity.ApartmentNumber)
		fmt.Println(space, entity.PostalCode, entity.City)
	}
}
