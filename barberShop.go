package main

import (
	"fmt"
	"time"
)

type BarberShop struct {
	SeatingCapacity int
	Barbers         int
	HaircutDuration time.Duration
	Open            bool
	ClientChannel   chan int
	BarberDoneChan  chan bool
}

func (shop *BarberShop) barberArrival(barberID int) {

	// A barber just arrived
	fmt.Printf("Barber%d just arrived at the shop\n", barberID)
	shop.Barbers++
	
	go func() {
		isNapping := false
		// Acitvate napping ability

		// Start an endless loop of checking and barbing
		for {

			// Check if there are client in the waiting room
			if len(shop.ClientChannel) == 0 {
				// There is no one around so he takes a nap
				fmt.Printf("There is no client waiting to get an haircut so barber%d takes a nap\n", barberID)

				// Toggle napping boolean
				isNapping = true
			}

			// Arrival of a client
			// If data is able to travel thru the channel, means the shop is open
			newClient, shopOpen := <-shop.ClientChannel

			if shopOpen {
				// Check if the barber is sleeping and wakw him up
				if isNapping {
					fmt.Printf("Client#%d wakes barber%d up for an haircut.\n", newClient, barberID)
					isNapping = false
				}

				// The client takes an haircut
				shop.cutHair(barberID, newClient)

			} else {
				// This means the shop is not open or the shop is closing
				// Get the Barber to finish up and go home
				shop.closeForDay(barberID)
				return

			}

		}
	}()

}

func (shop *BarberShop) cutHair(barberID, client int) {
	fmt.Printf("Barber%d is cutting client#%d's hair\n", barberID, client)
	time.Sleep(haircutDuration)
	fmt.Printf("client#%d has finished getting a haircut\n", client)

}

func (shop *BarberShop) closeForDay(barberID int) {
	shop.BarberDoneChan <- true
	fmt.Printf("Barber%d is done for the day\n", barberID)

}

func (shop *BarberShop) closingSoon() {
	fmt.Println("Shop is closing soon")
	close(shop.ClientChannel)
	shop.Open = false

	// Wait for the barbers to send their done signal
	for i := 0; i < shop.Barbers; i++ {
		<-shop.BarberDoneChan
	}

	close(shop.BarberDoneChan)
	fmt.Println("The Barber shop is now closed!!!")

}

func (shop *BarberShop) clientArrival(clientId int) {
	fmt.Println("Client#", clientId, " arrives")

	if shop.Open {
		select {
		case shop.ClientChannel <- clientId:
			fmt.Println("Client#", clientId, " takes a seat")
		default:
			fmt.Println("The waiting area is filled client#", clientId, " takes their leave.")

		}

	} else {
		fmt.Println("The barber shop is closed so client", clientId, " leaves.")
	}
}
