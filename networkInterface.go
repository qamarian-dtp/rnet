package rnet

type NetworkI struct {}

func (i *NetworkI) Send () {}

func (i *NetworkI) Read () {}

func (i *NetworkI) DirectLink () {}

func (i *NetworkI) Close () {}

func (i *NetworkI) Open () {}

func (i *NetworkI) Release () {}
