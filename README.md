# eqm - Embedded Queue Manager
Simple Queue Manager Implementaiton for Go Prgrams


This simple lib uses filesystem to enqueue / dequeue messages for local usage.

Portability and simplicity are priority here, not performance. Though it will work well for most of use cases.

```go
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
```
