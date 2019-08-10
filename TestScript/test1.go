package main

import (
	"fmt"
	"github.com/qamarian-dtp/rnet"
	"runtime"
	"time"
)

func sender (i *rnet.Interface, r string, n *rnet.Network) {
	for c := 1; c <= 20; c ++ {
		fmt.Println ("Sending...")
		errX := i.Send (c, r)
		if errX != nil {
			fmt.Println ("Message sending failed:", errX.Error ())
		}
		time.Sleep (time.Second * 1)
		if c == 10 {
			errT := n.Disconnect (r)
			fmt.Println (errT)
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
		fmt.Println (i.User (), i.IntfID (), i.NetAddr (), i.NetState (),
		mssg.(int))
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
	intf3, _    := net.NewIntf ("rt3", "anth")
	intf2, errR := net.NewIntf ("rt2", "recv")
	if errS != nil || errR != nil {
		fmt.Println ("A interface could not be created:", errS, errR)
		return
	}
	go reader (intf2)
	go reader (intf3)
	sender (intf1, "recv", net)
	sender (intf1, "anth", net)
	for {
		runtime.Gosched ()
	}
}
