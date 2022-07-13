package helpers

import (
	"fmt"
	"log"
	"sync"
)

var wg sync.WaitGroup

func Background(fn func()) {
	wg.Add(1)

	go func() {

		defer wg.Done()

		defer func() {
			if err := recover(); err != nil {
				log.Println(fmt.Errorf("%s", err))
			}
		}()

		fn()
	}()
}
