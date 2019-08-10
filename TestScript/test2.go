package main

import (
	"fmt"
	"github.com/qamarian-dtp/rnet"
	"runtime"
	"time"
)

func sender (i *rnet.Interface, r string) {
	for c := 1; c <= 10; c ++ {
		fmt.Println ("Sending...")
		errX := i.Send (i.NetAddr (), r)
		if errX != nil {
			fmt.Println ("Message sending failed:", errX.Error ())
		}
		time.Sleep (time.Second * 1)
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
		fmt.Println (mssg.(string))
	}
}

func main () {
	fmt.Println ("Test has started.")
	net, errX := rnet.New ()
	if errX != nil {
		fmt.Println ("Network creation failed. Error:", errX.Error ())
		return
	}
	intf1, err1 := net.NewIntf ("net-addr-1")
	intf2, err2 := net.NewIntf ("net-addr-2")
	intf3, err3 := net.NewIntf ("net-addr-3")
	if err1 != nil || err2 != nil || err3 != nil {
		fmt.Println ("An interface could not be created:", err1, err2, err3)
		return
	}
	go reader (intf1)
	go sender (intf2, "net-addr-1")
	go sender (intf3, "net-addr-1")
	for {
		runtime.Gosched ()
	}
}
