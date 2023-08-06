package njson

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/stretchr/testify/assert"
)

const genCount = 20000000

func Test_genfile(t *testing.T) {
	f, err := ioutil.TempFile(os.TempDir(), "*.njson")
	assert.Nil(t, err)
	// defer os.Remove(f.Name())

	for i := 0; i < genCount; i++ {
		fmt.Fprintln(f, `{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5, "f": 6, "g": 7, "h": 8, "I": 9}`)
	}
	fmt.Println(f.Name())
}

func TestCountN_json(t *testing.T) {
	fileName := "/tmp/4223751050.njson"
	f, err := os.Open(fileName)
	assert.Nil(t, err)

	count := 0
	d := json.NewDecoder(f)
	for {
		var v interface{}
		err := d.Decode(&v)
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			break
		}
		// fmt.Println(v)
		count += 1
	}
	assert.Equal(t, genCount, count)
	t.Log(count)
}

func TestCountN_sonic(t *testing.T) {
	fileName := "/tmp/4223751050.njson"
	f, err := os.Open(fileName)
	assert.Nil(t, err)

	count := 0
	d := sonic.ConfigDefault.NewDecoder(f)
	for {
		var v interface{}
		err := d.Decode(&v)
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			break
		}
		// fmt.Println(v)
		count += 1
	}
	assert.Equal(t, genCount, count)
	t.Log(count)
}

func TestCountN_scan(t *testing.T) {
	fileName := "/tmp/4223751050.njson"
	f, err := os.Open(fileName)
	assert.Nil(t, err)

	count := 0
	bfo := bufio.NewReader(f)
	for {
		_, _, err := bfo.ReadLine()
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			break
		}
		count += 1
	}
	assert.Equal(t, genCount, count)
	t.Log(count)
}

func BenchmarkCountN(b *testing.B) {
	// d := json.NewDecoder(strings.NewReader(stream))
	// for {
	// 	// Decode one JSON document.
	// 	var v interface{}
	// 	err := d.Decode(&v)

	// 	if err != nil {
	// 		// io.EOF is expected at end of stream.
	// 		if err != io.EOF {
	// 			log.Fatal(err)
	// 		}
	// 		break
	// 	}

	// 	// Do something with the value.
	// 	fmt.Println(v)
	// }
}
