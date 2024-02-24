package main

import (
	"fmt"

	gt "github.com/tehzwen/go-transmission"
)

func main() {
	var err error
	var r *gt.TransmissionResponse
	var magnet = "foo-bar"

	c := gt.NewTransmissionClient("user", "password", "server", 9091)
	if err := c.Login(); err != nil {
		panic(err)
	}

	r, err = c.AddTorrent(magnet)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", r)

	r, err = c.GetTorrents()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", r)

	r, err = c.StopTorrent(r.Arguments.Torrents[0].Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", r)

	r, err = c.StartTorrent(r.Arguments.Torrents[0].Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", r)
}
