package main

import (
	"fmt"
)

// link of the documentation
// https://go.dev/doc/effective_go#channels

// Types of channel - i)Buffered ii)Unbuffered

/*
If the channel is unbuffered, the sender blocks until the receiver has received the value. (direct from sender to receiver)
If the channel has a buffer, the sender blocks only until the value has been copied to the buffer; if the buffer is full, this means waiting until some receiver has retrieved a value.
(sender -> buffer -> receiver)
*/

/*
When to use which channel:
Unbuffered - strict synchronization is required b/w go routines; Tightly coupled; Avoid memory overhead; when data is received or not is more imp than data is sent; Signaling/acknowledgment is required

Buffered - Async behaviour; Ideal when queueing is required
*/

func main() {
	//Channel: The channel's buffer is initialized with the specified
	//	buffer capacity. If zero, or the size is omitted, the channel is
	//	unbuffered.
	userCh := make(chan string, 2)
	//var userCh <-chan string - receive only channel

	// so when I don't input the size, it basically creates unbuffered channel and give deadlock error as there is no go routine to consume the content of the channel
	// if i put the consuming code in goroutine then there is no error

	//userCh <- "Glock"

	//name := <-userCh

	// fmt.Println(name)
	//fmt.Println(name1)

	sendMessage(userCh)
	readMessage(userCh)

	// this is a Go label
	// use to target specific loop
free:
	for i := 0; i < 10; i++ {
		fmt.Println(i)
		for j := i; j < i+1; j++ {
			fmt.Println(j)
			if i == 4 && j == 4 {
				fmt.Println("breaking")
				break free
			}
		}
		fmt.Println("outer loop")
	}
}

func sendMessage(msgCh chan<- string) { // msgCh chan<- string - this is a send only channel
	msgCh <- "Hello"
}

func readMessage(msgCh <-chan string) { // msgCh <-chan string - this is a receive only channel
	msg := <-msgCh

	fmt.Println(msg)

}
