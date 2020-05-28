package main

import (
	"flag"
	"log"

	"github.com/Project-Wartemis/pw-engine/pkg/engine"
	"github.com/sirupsen/logrus"
)

var addr = flag.String("addr", "pw-backend:80", "http service address")

func main() {
	flag.Parse()
	logrus.SetLevel(logrus.DebugLevel)
	log.Println("Execute main")

	conquest := engine.NewConquestEngine()
	conquest.Start(*addr, "/socket")
}
