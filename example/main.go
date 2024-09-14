package main

import (
	"context"
	"log"
	"os"

	"github.com/aditya37/logger"
)

func init() {
	log.Println(os.Getenv("ELASTIC_APM_SERVICE_NAME"))
}

type Body struct {
	Name string
}

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, logger.ContextKey(logger.TraceId), "111222121")
	body := Body{
		Name: "anis",
	}
	logger.ErrorWithContext(ctx, body, nil, "ganbate")
}

func Agus() {
	logger.ErrorWithContext(context.Background(), nil, "ganbate")
}
