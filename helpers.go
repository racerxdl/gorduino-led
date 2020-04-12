package main

import (
	linuxproc "github.com/c9s/goprocinfo/linux"
	"github.com/lucasb-eyer/go-colorful"
	"image/color"
	"io"
	"log"
	"time"
)

func colorEqual(a, b color.RGBA) bool {
	return a.R == b.R && a.G == b.G && a.B == b.G
}

func writeRGB(port io.ReadWriteCloser, colors []color.RGBA) error {
	for i := 0; i < len(colors); i++ {
		c := colors[i]
		_, err := port.Write([]byte{c.R, c.G, c.B})
		if err != nil {
			return err
		}
	}
	return nil
}

func ledLevelFunc(intensity, level float32) color.RGBA {
	h := 120 - (120 * level)
	c := colorful.Hsv(float64(h), 1, float64(intensity))
	r, g, b := c.RGB255()

	return color.RGBA{R: r, G: g, B: b}
}

func getCPUPercent() float32 {
	stat0, err := linuxproc.ReadStat("/proc/stat")
	if err != nil {
		log.Fatal("stat read fail")
	}
	time.Sleep(time.Millisecond * 100)
	stat1, err := linuxproc.ReadStat("/proc/stat")
	if err != nil {
		log.Fatal("stat read fail")
	}

	return calcSingleCoreUsage(stat1.CPUStatAll, stat0.CPUStatAll)
}

func calcSingleCoreUsage(curr, prev linuxproc.CPUStat) float32 {

	PrevIdle := prev.Idle + prev.IOWait
	Idle := curr.Idle + curr.IOWait

	PrevNonIdle := prev.User + prev.Nice + prev.System + prev.IRQ + prev.SoftIRQ + prev.Steal
	NonIdle := curr.User + curr.Nice + curr.System + curr.IRQ + curr.SoftIRQ + curr.Steal

	PrevTotal := PrevIdle + PrevNonIdle
	Total := Idle + NonIdle
	// fmt.Println(PrevIdle, Idle, PrevNonIdle, NonIdle, PrevTotal, Total)

	//  differentiate: actual value minus the previous one
	totald := Total - PrevTotal
	idled := Idle - PrevIdle

	CPU_Percentage := (float32(totald) - float32(idled)) / float32(totald)

	return CPU_Percentage
}
