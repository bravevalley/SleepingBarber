package main

import (
	"fmt"
	"sync"
	"time"
)

type BarberShop struct {
	SeatingArea         int
	Barbers             int
	BarbingTime         time.Duration
	ClientsChan         chan int
	ShopOpen            bool
	finishedHairCutchan chan bool
	cleanShop           bool
	cleanShopChan       sync.Mutex
}

func (shop *BarberShop) BarbHair(barberID int) {
	// Increase the number of barbers each time the
	// func is called
	shop.Barbers++

	// Start a go routine to simulate the present barber's
	// day
	go func() {
		// Barber's arrival at the shop
		fmt.Printf("Barber%d just arrived at the shop\n", barberID)

		// If he is the first barber at the shop, he cleans the shop
		if !shop.cleanShop {
			fmt.Printf("Barber%d cleans the shop\n", barberID)

			// Updates the shop that he/she has cleaned the shop
			// Using a channel
			shop.cleanShopChan.Lock()
			shop.cleanShop = true
			shop.cleanShopChan.Unlock()
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
				fmt.Printf("No client in the waiting room so Barber%d goes to sleep\n", barberID)
			}

			// A client arrives
			// This means a new value has been passed thru the clientchannel
			newClient, openForBiz := <-shop.ClientsChan

			// Openfor Business is the commaOk idiom to test if actual data comes thru the channel
			if openForBiz {
				// Check if the barber is sleeping and wake him/her up
				if barberSleeping {
					fmt.Printf("Client#%d wakes barber#%d up...\n", newClient, barberID)
				}

				// Get the haircut
				fmt.Printf("Client#%d is getting his haircut\n", newClient)
				time.Sleep(shop.BarbingTime)
				fmt.Printf("Client#%d has finished getting a haricut from Barber%d\n", newClient, barberID)

			} else {
				// The only case we do not get data is if the channel is closed which mean the shop is closed

				shop.finishedHairCutchan <- true

				fmt.Println("Shop is closed for today.")

				return

			}
		}

	}()

}

func (shop *BarberShop) sendClient(client int) {

	// Data came thru, inform the shop
	fmt.Printf("\t\t\t\tclient#%d arrives\n", client)
	if shop.ShopOpen {
		select {
		case shop.ClientsChan <- client:
			fmt.Printf("client#%d takes a sit", client)
		default:
			fmt.Printf("The waiting area is full, client#%d bails on the shop.", client)

		}

	}

}

func (shop *BarberShop) closingSoon() {
	fmt.Println("Shop is closing")

	close(shop.ClientsChan)

	for i := 0; i < shop.Barbers; i++ {
		<-shop.finishedHairCutchan
	}

	fmt.Printf("The shop is now closed!!!")

}
