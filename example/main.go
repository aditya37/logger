package main

import (
	"log"
	"os"

	"github.com/aditya37/logger"
	"go.elastic.co/apm"
)

func init() {
	log.Println(os.Getenv("ELASTIC_APM_SERVICE_NAME"))
}
func main() {
	tx := apm.DefaultTracer.StartTransaction(os.Getenv("ELASTIC_APM_SERVICE_NAME"), "test logget")
	defer tx.End()
	logger.Error("error")
	defer func() {
		if v := recover(); v != nil {
			e := apm.DefaultTracer.Recovered(v)
			e.SetTransaction(tx) // or e.SetSpan(span)
			e.Send()
		}
	}()

}
