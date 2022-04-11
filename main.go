package main

import (
	"github.com/massimo-gollo/DASHpher/models"
	"github.com/massimo-gollo/DASHpher/reproduction"
	"log"
)

func main() {
	url := "http://cloud.gollo1.particles.dieei.unict.it/videofiles/624d6b6a7744c2250ff741ff/video.mpd"
	rm := models.ReproductionMetrics{Url: url}
	err := reproduction.Stream1(&rm, "h264", "conventional", 1080, 240000, 2, 5, 0)
	if err != nil {
		log.Println("error: ", err.Error())
	}
}
