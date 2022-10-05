// This is a simple demonstration of how to solve the Sleeping Barber dilemma, a classic computer science problem
// which illustrates the complexities that arise when there are multiple operating system processes. Here, we have
// a finite number of barbers, a finite number of seats in a waiting room, a fixed length of time the barbershop is
// open, and clients arriving at (roughly) regular intervals. When a barber has nothing to do, he or she checks the
// waiting room for new clients, and if one or more is there, a haircut takes place. Otherwise, the barber goes to
// sleep until a new client arrives. So the rules are as follows:
//
//		- if there are no customers, the barber falls asleep in the chair
//		- a customer must wake the barber if he is asleep
//		- if a customer arrives while the barber is working, the customer leaves if all chairs are occupied and
//		  sits in an empty chair if it's available
//		- when the barber finishes a haircut, he inspects the waiting room to see if there are any waiting customers
//		  and falls asleep if there are none
// 		- shop can stop accepting new clients at closing time, but the barbers cannot leave until the waiting room is
//	      empty
//		- after the shop is closed and there are no clients left in the waiting area, the barber
//		  goes home
//
// The Sleeping Barber was originally proposed in 1965 by computer science pioneer Edsger Dijkstra.
//
// The point of this problem, and its solution, was to make it clear that in a lot of cases, the use of
// semaphores (mutexes) is not needed.

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var (
	seatingCapacity = 5
	noBarbers       = 1
	hairCutTime     = 100
	shopOpenTime    = 10 * time.Second
	arrivalInterval = 200
)

func main() {

	rand.Seed(time.Now().Unix())

	clientChannel := make(chan int, seatingCapacity)
	finishedHairCut := make(chan bool)

	// Print the program starting

	fmt.Println("------------------ Program Starting ------------------")

	allenshop := BarberShop{
		SeatingArea:         seatingCapacity,
		Barbers:             noBarbers,
		BarbingTime:         time.Duration(hairCutTime * int(time.Millisecond)),
		ClientsChan:         clientChannel,
		ShopOpen:            false,
		finishedHairCutchan: finishedHairCut,
		cleanShop:           false,
		cleanShopChan:       sync.Mutex{},
	}

	// Simulate a day in the life of a barbing shop
	allenshop.BarbHair(1)
	allenshop.BarbHair(2)

	// We need a channel for when the shop is about closing

	closingChannel := make(chan bool)
	shopCloseChannel := make(chan bool)

	go func() {
		<-time.After(shopOpenTime)
		<-closingChannel
		allenshop.closingSoon()
		<-shopCloseChannel
		close(shopCloseChannel)

	}()

	go func() {
		var i int
		for {
			randomization := rand.Int() * (2 % arrivalInterval)
			select {
			case <-closingChannel:
				return
			case <-time.After(time.Millisecond * time.Duration(randomization)):
				allenshop.sendClient(i)
				i++
			}
		}
	}()

	<-shopCloseChannel

	fmt.Println("Goodbye")

}
