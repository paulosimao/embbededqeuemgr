package eqm

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

func Test_Basic(t *testing.T) {
	os.RemoveAll("db")
	q, err := New()
	if err != nil {
		t.Fatal(err)
	}
	err = q.Push([]byte("entry1"))
	if err != nil {
		t.Fatal(err)
	}
	q.Push([]byte("entry2"))
	if err != nil {
		t.Fatal(err)
	}
	q.Push([]byte("entry3"))
	if err != nil {
		t.Fatal(err)
	}
	bs, err := q.Peek()
	if err != nil {
		t.Fatal(err)
	}
	if string(bs) != "entry1" {
		t.Fatal("Not matching: " + string(bs) + " entry1")
	}

	for err == nil && len(bs) > 0 {
		bs, err = q.Pop()
		if err != nil {
			t.Fatal(err)
		}
		log.Printf(string(bs))
	}

}

func Test_Multi(t *testing.T) {
	os.RemoveAll("db")
	q, err := New()
	if err != nil {
		panic(err)
	}
	wcount := 0
	rcount := 0
	go func() {
		for {
			q.Push([]byte(fmt.Sprintf("Message: %d", wcount)))
			if wcount > 50 {
				time.Sleep(time.Millisecond * 200)
			} else {
				time.Sleep(time.Millisecond * 50)
			}
			wcount++
		}
	}()
	go func() {
		time.Sleep(time.Millisecond * 200)
		for {
			bs, err := q.Pop()
			if err != nil {
				log.Printf(err.Error())
				if err.Error() == "Queue Empty" {
					panic(err)
				}
			}
			log.Printf("R1 - Got from queue: %s", string(bs))
			if rcount < 50 {
				time.Sleep(time.Millisecond * 200)
			} else {
				time.Sleep(time.Millisecond * 100)
			}
			rcount++
		}
	}()

	go func() {
		time.Sleep(time.Millisecond * 200)
		for {
			bs, err := q.Pop()
			if err != nil {
				log.Printf(err.Error())
				if err.Error() == "Queue Empty" {
					panic(err)
				}
			}
			log.Printf("R2 - Got from queue: %s", string(bs))
			if rcount < 50 {
				time.Sleep(time.Millisecond * 200)
			} else {
				time.Sleep(time.Millisecond * 100)
			}
			rcount++
		}
	}()

	for {
		time.Sleep(time.Minute)
	}

}
