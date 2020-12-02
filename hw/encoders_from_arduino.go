package main

import (
	"fmt"
	"github.com/bugst/go-serial"
	"github.com/eranet/rhombus/rhomgo"
	"log"
	"strconv"
	"strings"
)

type Encoders struct {
	Pos1 float64
	Vel1 float64
	Pos2 float64
	Vel2 float64
}

func main() {
	c := rhomgo.LocalJSONConnection()
	defer c.Close()

	mode := &serial.Mode{
		BaudRate: 115200,
	}
	port, err := serial.Open("", mode) //TODO: pass value
	if err != nil {
		log.Fatal(err)
	}

	buff := make([]byte, 100)
	for {
		// Reads up to 100 bytes
		n, err := port.Read(buff)
		if err != nil {
			log.Fatal(err)
		}
		if n == 0 {
			fmt.Println("\nEOF")
			break
		}

		data := string(buff[:n])
		data = strings.TrimSpace(data)
		fmt.Printf("%s", data)
		vals := strings.Split(data, " ")

		err = c.Publish("/encoders", Encoders{Pos1: toF(vals[0]), Vel1: toF(vals[1]), Pos2: toF(vals[2]), Vel2: toF(vals[3])})
		if err != nil {
			println("error pub:", err)
		}
	}

}
func toF(s string) float64 {
	res, _ := strconv.ParseFloat(s, 64)
	return res
}
