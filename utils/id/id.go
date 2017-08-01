package idutils

import (
	"encoding/base64"
	"encoding/binary"
	"time"
)

type Generator struct {
	ids        chan string
	startPoint uint64
	max        int // channel max buf
}

func (g *Generator) GetID() string {
	return <-g.ids
}

func NewGenerator(max int) *Generator {
	g := &Generator{
		ids:        make(chan string, max),
		startPoint: uint64(time.Now().UnixNano() / int64(time.Millisecond)),
		max:        max,
	}
	go func() {
		for {
			select {
			case <-time.After(time.Millisecond):
				if len(g.ids) == g.max {
					continue
				}
				g.startPoint++
				v := uint64(1)<<63 + g.startPoint<<40 + 1<<10
				b := make([]byte, 8)
				for i := uint64(0); i < uint64(g.max); i++ {
					binary.LittleEndian.PutUint64(b, v+i)
					g.ids <- base64.StdEncoding.EncodeToString(b)
				}
			}
		}
	}()

	return g
}

var DefaultGenerator = NewGenerator(100)
