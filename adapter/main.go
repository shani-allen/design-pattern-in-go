package main

import "fmt"

//adapter design pattern

//WHAT
//WHY
//HOW

//WHAT:
// comes under structural design pattern
// also knows as Wrapper.
//its used so that two unrelated object can work together using adapter.
// the thing that join these two unrelated object is called an adapter.
// RTE

//four participants while implementation:
//Target Interface: this is the interface which will be used by the clients
//adapter: this the wrapper which implements the target interface and modifies the specific request
//available from the adaptee class
//adaptee: this is the object which is used by the adapter to reuse the functionality and modify them for desired use.

//WHY
//when you dont need to change the existing object or interface rather
// wants to add new functionality on top of what is existing already.

// HOW:
// target
type mobile interface {
	chargeAppleMobile()
}

// concrete prototype implementation
type apple struct {
}

func (a *apple) chargeAppleMobile() {
	fmt.Println("charging an apple")
}

// adaptee
type android struct {
}

func (a *android) chargeAndroidMobile() {
	fmt.Println("charging an android")
}

// Adapter
type androidAdapter struct {
	android *android
}

func (a *androidAdapter) chargeAppleMobile() {
	a.android.chargeAndroidMobile()
}

// client
type Client struct {
}

func (c *Client) chargeMobile(mob mobile) {
	mob.chargeAppleMobile()
}

func main() {
	// initial requirement
	client := Client{}
	apple := &apple{}
	client.chargeMobile(apple)

	//extended requirement
	ad := android{}
	adapterAndroid := &androidAdapter{
		android: &ad,
	}
	client.chargeMobile(adapterAndroid)
}
