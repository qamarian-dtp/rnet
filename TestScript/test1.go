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
		fmt.Println (i.NetAddr (), i.NetState (), mssg.(int))
	}
}

func main () {
	fmt.Println ("Test has started.")
	net := rnet.New ()
	intf1, err1 := net.NewIntf ("send")
	intf2, err2 := net.NewIntf ("recv")
	intf3, err3 := net.NewIntf ("anth")
	if err1 != nil || err2 != nil || err3 != nil {
		fmt.Println ("An interface could not be created:", err1, err2, err3)
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
