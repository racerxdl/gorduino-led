package main

import (
	"fmt"
	"github.com/jacobsa/go-serial/serial"
	"log"
	"time"
)

func main() {
	options := serial.OpenOptions{
		PortName:        "/dev/ttySAC0",
		BaudRate:        115200,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}

	// Open the port.
	port, err := serial.Open(options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}

	// Make sure to close it later.
	defer port.Close()

	lc := makeLedController(3, port)
	fmt.Println("Sending breath")
	lc.c <- ledInstruction{
		cmd: LedError,
	}

	//for i := 0; i < 10; i++ {
	//    level := getCPUPercent()
	//    lc.c <- ledInstruction{
	//        cmd: LedLevel,
	//        level: level,
	//    }
	//    fmt.Printf("Level: %03.0f%%\n", math.Floor(float64(level * 100)))
	//}

	time.Sleep(time.Second * 5)
	fmt.Println("Sending stop")
	lc.c <- ledInstruction{
		cmd: LedStop,
	}

	time.Sleep(time.Second * 2)
	fmt.Println("All done")
}
