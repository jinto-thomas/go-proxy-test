package main

import "encoding/json"
import "fmt"
import "sync"

type CircularBuffer struct {
	max          int
	start        int
	end          int
	front        int
	rear         int
	buffer       [10 * 1024]byte
	mutex        *sync.Mutex
}

func initBuffer() CircularBuffer {
	fmt.Println("new buffer..")
	buffer := CircularBuffer{}
	buffer.max = 10 * 1024
	buffer.start = 0
	buffer.end = 0
	buffer.front = 0
	buffer.rear = 1
	buffer.mutex = new(sync.Mutex)
	return buffer
}

func (buf *CircularBuffer) capacity() int {
	return buf.max
}

func (buffer *CircularBuffer) reset() {
	fmt.Println("###############RESET##############")
	buffer.max = 10 * 1024
	buffer.start = 0
	buffer.end = 0
	buffer.front = 0
	buffer.rear = 1
	for i := range buffer.buffer {
		buffer.buffer[i] = 0
	}
}

func (buf *CircularBuffer) clearBuffer(start int, end int) {
	for i := start; i < end; i++ {
		buf.buffer[i] = 0
	}
}

func (buf *CircularBuffer) write(rcvd []byte, size int) {
	buf.mutex.Lock()
	if buf.rear+size <= buf.max {
		copy(buf.buffer[buf.rear:], rcvd[:size])
		buf.rear = (buf.rear + size) % buf.max
	} else {
		writeableByte := buf.max - buf.rear

		copy(buf.buffer[buf.rear:], rcvd[:writeableByte])
		copy(buf.buffer[0:], rcvd[writeableByte:size])
		buf.rear = size - writeableByte
	}
	buf.mutex.Unlock()
	//fmt.Println("Actual Bufffer ", buf.buffer )
}

func (buf *CircularBuffer) spaceAvailable() int {
	var size int
	if buf.rear > buf.front {
		size = buf.front + buf.max - buf.rear
	} else {
		size = buf.front - buf.rear
	}
	return size
}

func (buf *CircularBuffer) dataAvailable() int {
	return buf.max - 1 - buf.spaceAvailable()
}

func (buf *CircularBuffer) process(ch chan<- JsonQuote) {
	var quote JsonQuote
	var start int

	partialRead := true

	buf.mutex.Lock()
	for i := buf.end; i < buf.max; i++ {

		if buf.buffer[i] == 10 {
			if buf.end == 0 {
				start = 1
			} else {
				start = buf.end
			}

			err := json.Unmarshal(buf.buffer[start:i], &quote)
			buf.clearBuffer(start, i)
			if err == nil {
				fmt.Println(quote)
				ch <- quote
			} else {
				fmt.Println("[A].Unmarshall error ", err)
				fmt.Printf("[A], start %d, end %d\n", start, i)
				fmt.Println("[A].Error json ", buf.buffer)
			}

			buf.end = i + 1
			fmt.Println("[original] ", buf.end)
			partialRead = false
			break
		}

	}

	if partialRead == true {
		temp := make([]byte, 1024*2)
		size := len(buf.buffer[buf.end:])
		copy(temp, buf.buffer[buf.end:]) // or buf.buffer[buf.end:]

		buf.clearBuffer(buf.end, buf.max)

		for i := range buf.buffer {
			temp[size+i] = buf.buffer[i]
			if buf.buffer[i] == 10 {
				//	fmt.Printf("new temp start %d, (init %d)\n", size + i, size )
				err := json.Unmarshal(temp[:size+i], &quote)
			  buf.clearBuffer(0, i)
				if err == nil {
					fmt.Println("++", quote)
					ch <- quote

				} else {
					fmt.Printf("[B]Error : first char %d, last char %d, newline %d\n", temp[0], buf.buffer[i-1], buf.buffer[i])
					fmt.Printf("[B]Error : first char %d, last char %d ,2nd last char %d in temp\n", temp[0], temp[size+i], temp[size+i-1])
					fmt.Println("[B].Error json ", buf.buffer)
					fmt.Println("[B]++Unmarshall error ", err)
				}
				buf.end = i + 1
				fmt.Println("[partial] ", buf.end)
				break
			}
		}
	}
	buf.mutex.Unlock()
}
