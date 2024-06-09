package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"log"
	"math/rand"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sherpaurgen/msvc/src/federation/dto"
	"github.com/sherpaurgen/msvc/src/federation/utils"
)

var name = flag.String("name", "sensor", "name of sensor")
var freq = flag.Uint("freq", 5, "update frequency in cycle/sec")
var max = flag.Float64("max", 5., "max value for generated unit")
var min = flag.Float64("min", 1., "min value for generated unit")
var stepSize = flag.Float64("step", 0.1, "max allowable value to deviate")
var url = "amqp://guest:guest@localhost:5672/"

var r = rand.New(rand.NewSource(time.Now().UnixNano()))
var value = r.Float64()*(*max-*min) + *min
var nom = (*max-*min)/2 + *min

func calcValue() {
	var maxStep, minStep float64
	if value < nom {
		maxStep = *stepSize
		minStep = -1 * *stepSize * (value - *min) / (nom - *min)
	} else {
		maxStep = *stepSize * (*max - value) / (*max - nom)
		minStep = -1 * *stepSize
	}
	value += r.Float64()*(maxStep-minStep) + minStep
}

func main() {
	flag.Parse()
	conn, ch := utils.GetChannel(url)
	defer conn.Close()
	defer ch.Close()

	dataQueue := utils.GetQueue(*name, ch)

	dur, _ := time.ParseDuration(strconv.Itoa(1000/int(*freq)) + "ms") //eg here 1 second or 1000ms divided by 5ms

	signal := time.Tick(dur)
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	for range signal { //loops every X millisecond
		calcValue()
		log.Printf("Reading sent value %v\n", value)
		reading := dto.SensorMessage{
			Name:      *name,
			Value:     value,
			Timestamp: time.Now(),
		}
		buf.Reset()
		enc.Encode(reading) //encodes the reading struct into the binary format and writes it to the buffer buf.

		msg := amqp.Publishing{
			Body: buf.Bytes(),
		}
		ch.Publish(
			"", //direct exchange as its empty
			dataQueue.Name,
			false,
			false,
			msg,
		)

	}
}
