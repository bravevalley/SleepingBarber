// This is a simple demonstration of how to solve the Sleeping Barber dilemma, a classic computer science problem
// which illustrates the complexities that arise when there are multiple operating system processes. Here, we have
// a finite number of barbers, a finite number of seats in a waiting room, a fixed length of time the barbershop is
// open, and clients arriving at (roughly) regular intervals. When a barber has nothing to do, he or she checks the
// waiting room for new clients, and if one or more is there, a haircut takes place. Otherwise, the barber goes to
// sleep until a new client arrives. So the rules are as follows:
//
//   - if there are no customers, the barber falls asleep in the chair
//   - a customer must wake the barber if he is asleep
//   - if a customer arrives while the barber is working, the customer leaves if all chairs are occupied and
//     sits in an empty chair if it's available
//   - when the barber finishes a haircut, he inspects the waiting room to see if there are any waiting customers
//     and falls asleep if there are none
//   - shop can stop accepting new clients at closing time, but the barbers cannot leave until the waiting room is
//     empty
//   - after the shop is closed and there are no clients left in the waiting area, the barber
//     goes home
//
// The Sleeping Barber was originally proposed in 1965 by computer science pioneer Edsger Dijkstra.
//
// The point of this problem, and its solution, was to make it clear that in a lot of cases, the use of
// semaphores (mutexes) is not needed.
package main

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	waitingRoom     = 4
	arrivalInterval = 200
	haircutDuration = 1000 * time.Millisecond
	timeOpen        = 5 * time.Second
)

func main() {
	// Seed a rondom number
	rand.Seed(time.Now().UnixNano())

	// Print welcome message
	fmt.Println("------------- Program starts -------------")

	// Create a chaneel to interface with the shop and barbers
	clientChannel := make(chan int, waitingRoom)
	barberCloseChan := make(chan bool)

	// Create the shop it self
	allenshop := BarberShop{
		SeatingCapacity: waitingRoom,
		Barbers:         0,
		HaircutDuration: haircutDuration,
		Open:            true,
		ClientChannel:   clientChannel,
		BarberDoneChan:  barberCloseChan,
	}

	allenshop.barberArrival(1)

	closingChannel := make(chan bool)
	shopClosedChannel := make(chan bool)

	go func() {
		<-time.After(timeOpen)
		closingChannel <- true
		allenshop.closingSoon()
		shopClosedChannel <- true
	}()

	// Add clients
	i := 1
	go func() {

		for {
			// Get random number
			randomization := rand.Int() % (arrivalInterval * 2)
			select {
				case <-closingChannel:
					return
				case <-time.After(time.Millisecond * time.Duration(randomization)):
					allenshop.clientArrival(i)
					i++
			}
		}
	}()

	<-shopClosedChannel
}
