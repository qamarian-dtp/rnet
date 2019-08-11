// rNet (Goroutine Network) is a set of data types/network that allows multiple
// goroutines to easily communicate with one another.
//
// By network, I do not mean something like the internet. No, it is not. See this
// network as a more versatile implementation of the Golang's channel.
//
// Assuming you want goroutines X, Y, and Z to communicate with one another, you could
// create a communication network (an rNet) for them.
//
// After creating the network, then create a network interface for each one of the
// goroutines.
//
// A network interface is simply a data which can be used to send and receive messages
// from other goroutines on the same network.
//
//	net := rnet.New () // Creation of a new communication network
//
//	netInterfaceForX, err1 := net.NewIntf ("x") /* Creation of a network
//		interface for goroutine X. */
//
//	netInterfaceForY, err2 := net.NewIntf ("y") /* Creation of a network
//		interface for goroutine Y. */
//
//	netInterfaceForZ, err3 := net.NewIntf ("z") /* Creation of a network
//		interface for goroutine Z. */
//
// Afterwards, these routines could communicate with one another, over the network,
// using methods Send () and Read ().
//
//	err1 := netInterfaceForX.Send ("hello", "y") /* Goroutine X sending "hello"
//		to goroutine Y. */
//
//	mssg, err2 := netInterfaceForY.Read () /* Goroutine Y checking for any
//		message that might have been sent to it. */
//
// Warning!
//
// When using this package, import a specific version of the package. This is because
// the latest version of this package is not guaranteed to be always backward
// compatible with the older versions of the package. Don't know how to import a
// specific version of a Go package? Read up the http://gopkg.in tool.
package rnet
