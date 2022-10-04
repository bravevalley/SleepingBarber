package main

import (
	"fmt"
	"time"
)

type BarberShop struct {
	SeatingArea int
	Barbers int
	BarbingTime time.Duration
	ClientsChan chan int
	ShopOpen bool
	cleanShop bool
	cleanShopChan chan bool
}

func (shop *BarberShop) BarbHair(barberID int) {
	// Increase the number of barbers each time the 
	// func is called
	shop.Barbers++ 

	// Start a go routine to simulate the present barber's
	// day
go func() {
	// Barber's arrival at the shop
	fmt.Printf("Barber%d just arrived at the shop", barberID)

	// If he is the first barber at the shop, he cleans the shop
	if !shop.cleanShop {
		fmt.Printf("Barber%d cleans the shop", barberID)

		// Updates the shop that he/she has cleaned the shop
		// Using a channel
		shop.cleanShopChan <- true

		// Close the channel
		close(shop.cleanShopChan)
	}

	// Barber sleep schedule, barber must be awake when he opens the shop
	var barberSleeping bool

	// Barber start the process of checking for clients in the waiting room
	// or sleeping
	for {
		// Checks if there is/are client in the waiting room
		// check the clientChannel for process queue
		if len(shop.ClientsChan) == 0 {
			barberSleeping = true
			fmt.Printf("No client in the waiting room so Barber%d goes to sleep", barberID)
		}

		// A client arrives
		// This means a new value has been passed thru the clientchannel
		newClient, openForBiz := <- shop.ClientsChan

		if openForBiz {
			fmt.Printf("\t\t\t\tclient#%d arrives", newClient)
			if barberSleeping {
				fmt.Printf("Client#%d wakes barber#%d up...", newClient, barberID)
			}
		}

	}

}()
	
}

