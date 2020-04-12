package main

import (
	"fmt"
	"golang.org/x/image/colornames"
	"image/color"
	"io"
	"time"
)

type ledCommand uint8

const (
	LedStop   ledCommand = iota
	LedOff    ledCommand = iota
	LedSetRgb ledCommand = iota
	LedBreath ledCommand = iota
	LedLevel  ledCommand = iota
	LedWarn   ledCommand = iota
	LedError  ledCommand = iota
)

const breathInterval = time.Second
const countsPerMillis = float32(time.Millisecond) / float32(breathInterval)

type ledInstruction struct {
	cmd    ledCommand
	colors []color.RGBA
	level  float32
}

type ledController struct {
	c                chan ledInstruction
	port             io.ReadWriteCloser
	currentMode      ledCommand
	numLeds          int
	colors           []color.RGBA
	lastColors       []color.RGBA
	breathColors     []color.RGBA
	breathVal        float32
	breathIncreasing bool
	lastBreath       time.Time
}

func (lc *ledController) loop() {
	running := true
	lc.resetColors()
	for running {
		select {
		case inst := <-lc.c:
			switch inst.cmd {
			case LedStop:
				running = false
				fmt.Println("Received Stop")
				lc.resetColors()
				lc.currentMode = LedStop
			case LedSetRgb:
				lc.currentMode = LedSetRgb
				lc.colors = inst.colors
			case LedBreath:
				if len(inst.colors) != lc.numLeds {
					fmt.Printf("Expected %d leds but got %d\n", lc.numLeds, len(inst.colors))
				} else {
					lc.breathColors = inst.colors
					lc.setBreath()
				}
			case LedLevel:
				lc.setLevel(inst.level)
				lc.currentMode = LedSetRgb
			case LedWarn:
				fmt.Println("Received WARN")
				if len(lc.breathColors) != lc.numLeds {
					lc.breathColors = make([]color.RGBA, lc.numLeds)
				}

				for i := 0; i < lc.numLeds; i++ {
					lc.breathColors[i] = colornames.Yellow
				}
				lc.setBreath()
			case LedError:
				if len(lc.breathColors) != lc.numLeds {
					lc.breathColors = make([]color.RGBA, lc.numLeds)
				}

				for i := 0; i < lc.numLeds; i++ {
					lc.breathColors[i] = colornames.Red
				}
				lc.setBreath()
			}
		default:
			// Nothing
		}

		lc.update()
		time.Sleep(time.Millisecond)
	}
}

func (lc *ledController) setLevel(level float32) {
	//ledSegments := 1 / float32(lc.numLeds)

	lc.colors[0] = ledLevelFunc(1, level)

	v := float32(lc.numLeds) * level

	for i := 0; i < lc.numLeds; i++ {
		nv := v
		if nv > 1 {
			nv = 1
		}
		v -= nv
		lc.colors[i] = ledLevelFunc(1, nv)
	}

	//
	//for v > 0 && currLed < lc.numLeds {
	//    nv := v
	//    if nv > 256 * 3 {
	//        nv = 256 * 3
	//    }
	//    v -= nv
	//
	//    if nv >= 256 * 2 { // RED
	//        nv -= 256 * 2
	//        nv /= 256
	//        lc.colors[currLed].R = uint8(float32(colornames.Red.R) * nv)
	//        lc.colors[currLed].G = uint8(float32(colornames.Red.G) * nv)
	//        lc.colors[currLed].B = uint8(float32(colornames.Red.B) * nv)
	//    } else if nv > 256 {
	//        nv -= 256
	//        nv /= 256
	//        lc.colors[currLed].R = uint8(float32(colornames.Yellow.R) * nv)
	//        lc.colors[currLed].G = uint8(float32(colornames.Yellow.G) * nv)
	//        lc.colors[currLed].B = uint8(float32(colornames.Yellow.B) * nv)
	//    } else {
	//        nv /= 256
	//        lc.colors[currLed].R = uint8(float32(colornames.Green.R) * nv)
	//        lc.colors[currLed].G = uint8(float32(colornames.Green.G) * nv)
	//        lc.colors[currLed].B = uint8(float32(colornames.Green.B) * nv)
	//    }
	//
	//    if v > 0 {
	//        currLed++
	//    }
	//}

	//// Do Green
	//for i := 0; i < lc.numLeds; i++ {
	//   nv := v
	//   if nv > 255 {
	//       nv = 255
	//   }
	//   v -= nv
	//
	//   nv /= 256
	//
	//    if nv > 0 {
	//        lc.colors[i].R = uint8(float32(colornames.Green.R) * nv)
	//        lc.colors[i].G = uint8(float32(colornames.Green.G) * nv)
	//        lc.colors[i].B = uint8(float32(colornames.Green.B) * nv)
	//    }
	//}
	//
	//// Do Yellow
	//for i := 0; i < lc.numLeds; i++ {
	//    nv := v
	//    if nv > 255 {
	//        nv = 255
	//    }
	//    v -= nv
	//
	//    nv /= 256
	//
	//    if nv > 0 {
	//        lc.colors[i].R = uint8(float32(colornames.Yellow.R) * nv)
	//        lc.colors[i].G = uint8(float32(colornames.Yellow.G) * nv)
	//        lc.colors[i].B = uint8(float32(colornames.Yellow.B) * nv)
	//    }
	//}
	//
	//// Do Red
	//for i := 0; i < lc.numLeds; i++ {
	//    nv := v
	//    if nv > 255 {
	//        nv = 255
	//    }
	//    v -= nv
	//
	//    nv /= 256
	//
	//    if nv > 0 {
	//        lc.colors[i].R = uint8(float32(colornames.Red.R) * nv)
	//        lc.colors[i].G = uint8(float32(colornames.Red.G) * nv)
	//        lc.colors[i].B = uint8(float32(colornames.Red.B) * nv)
	//    }
	//}

	lc.setRGB(lc.colors)
}

func (lc *ledController) resetColors() {
	lc.colors = make([]color.RGBA, lc.numLeds)
	for i := 0; i < lc.numLeds; i++ {
		lc.colors[i] = colornames.Black
	}
	lc.setRGB(lc.colors)
}

func (lc *ledController) setBreath() {
	lc.currentMode = LedBreath
	lc.lastBreath = time.Now()
	lc.breathVal = 0
	lc.breathIncreasing = true
}

func (lc *ledController) ledsNeedUpdate() bool {
	needsUpdate := false
	for i := 0; i < lc.numLeds; i++ {
		needsUpdate = needsUpdate || !colorEqual(lc.colors[i], lc.lastColors[i])
	}

	return needsUpdate
}

func (lc *ledController) update() {
	switch lc.currentMode {
	case LedSetRgb:
		if lc.ledsNeedUpdate() {
			lc.setRGB(lc.colors)
		}
	case LedBreath:
		if len(lc.breathColors) != len(lc.colors) {
			lc.colors = make([]color.RGBA, len(lc.colors))
		}
		delta := time.Now().Sub(lc.lastBreath)
		deltaCounts := countsPerMillis * float32(delta.Milliseconds())
		if lc.breathIncreasing {
			lc.breathVal += deltaCounts
			if lc.breathVal >= 1 {
				lc.breathIncreasing = false
				lc.breathVal = 1
			}
		} else {
			lc.breathVal -= deltaCounts
			if lc.breathVal <= 0 {
				lc.breathIncreasing = true
				lc.breathVal = 0
			}
		}

		//fmt.Printf("BREATH: %f\n", lc.breathVal)

		needsUpdate := false

		for i := 0; i < len(lc.breathColors); i++ {
			lc.colors[i].R = uint8(float32(lc.breathColors[i].R) * lc.breathVal)
			lc.colors[i].G = uint8(float32(lc.breathColors[i].G) * lc.breathVal)
			lc.colors[i].B = uint8(float32(lc.breathColors[i].B) * lc.breathVal)

			needsUpdate = needsUpdate || !colorEqual(lc.colors[i], lc.lastColors[i])
		}

		if needsUpdate {
			lc.setRGB(lc.colors)
		}

		lc.lastBreath = time.Now()
	}
}

func (lc *ledController) setRGB(colors []color.RGBA) {
	writeRGB(lc.port, colors)
	copy(lc.lastColors, colors)
}

func makeLedController(numLeds int, port io.ReadWriteCloser) *ledController {
	lc := &ledController{
		c:           make(chan ledInstruction),
		port:        port,
		currentMode: LedOff,
		lastColors:  make([]color.RGBA, numLeds),
		breathVal:   0,
		numLeds:     numLeds,
		lastBreath:  time.Now(),
	}

	go lc.loop()

	return lc
}
