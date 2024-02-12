package common

import (
	"io"
	"log"
)

func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Println(err)
	}
}
