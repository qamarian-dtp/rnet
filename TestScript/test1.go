package main

import (
	"fmt"
	"github.com/qamarian-dtp/rnet"
	"runtime"
	"time"
)

func sender (i *rnet.Interface) {
	for c := 1; c <= 20; c ++ {
		fmt.Println ("Sending...")
		errX := i.Send (c, "recv")
		if errX != nil {
			fmt.Println ("Message sending failed:", errX.Error ())
			return
		}
	}
}

func reader (i *rnet.Interface) {
	for {
		mssg, errX := i.Read ()
		if errX != nil {
			fmt.Println ("Message could not be read:", errX.Error ())
			return
		}
		if mssg == nil {
			fmt.Println ("No message.")
			time.Sleep (time.Second * 1)
			continue
		}
		fmt.Println (mssg.(int))
	}
}

func main () {
	fmt.Println ("Test has started.")
	net, errX := rnet.New ()
	if errX != nil {
		fmt.Println ("Network creation failed.")
		return
	}
	intf1, errS := net.NewIntf ("rt1", "send")
	intf2, errR := net.NewIntf ("rt2", "recv")
	if errS != nil || errR != nil {
		fmt.Println ("A interface could not be created:", errS, errR)
		return
	}
	go sender (intf1)
	go reader (intf2)
	for {
		runtime.Gosched ()
	}
}
