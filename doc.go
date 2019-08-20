// rNet (Goroutine Network) is a network (a set of data types) that allows multiple
// goroutines to easily communicate with one another.
//
// By network, I do not mean something like the internet. No, it is not. See this
// network as a more versatile implementation of the Golang's channel.
//
// Assuming you want goroutines X, Y, and Z to communicate with one another, you could
// create a communication network (an rNet) for them.
//
// After creating the network, then create a personal post office (PPO) for each one
// of the goroutines.
//
// A PPO is simply a data which can be used to send and receive messages from other
// goroutines on the same network.
//
//	net := rnet.New () // Creation of a new rNet (communication network)
//
//	xPPO, err1 := net.NewPPO ("x") // Creation of a PPO for goroutine X.
//
//	yPPO, err2 := net.NewPPO ("y") // Creation of a PPO for goroutine Y.
//
//	zPPO, err3 := net.NewPPO ("z") // Creation of a PPO for goroutine Z.
//
// Afterwards, these routines could communicate with one another, over the network,
// using methods Send () and Read ().
//
//	err1 := xPPO.Send ("hello", "y") // Goroutine X sending "hello world!" to Y.
//
//	mssg, err2 := yPPO.Read () /* Goroutine Y checking for any message that might
//		have been sent to it. */
//
// Warning!
//
// When using this package, import a specific version of the package. This is because
// the latest version of this package is not guaranteed to be always backward
// compatible with the older versions of the package. Don't know how to import a
// specific version of a Go package? Read up the http://gopkg.in tool.
package rnet
