package main

import (
	"fmt"
	"github.com/eranet/rhombus/rhomgo"
)

type Position struct {
	Value float64
}

func main() {
	c := rhomgo.LocalJSONConnection()
	defer c.Close()

	c.Subscribe("/current", func(p *Position) {
		fmt.Printf("Received a position: %+v\n", p)
	})

	rate := rhomgo.NewRate(1)
	for i := 1.0; i > -1; i -= 0.001 {
		err := c.Publish("/pendabot/shoulder_torque_controller/command", Position{Value: i})
		if err != nil {
			println("error pub:", err)
		}
		rate.Sleep()
	}

}
