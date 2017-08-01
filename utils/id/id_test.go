package idutils

import (
	"fmt"
	"testing"
	"time"
)

func TestID(t *testing.T) {
	max := 100
	g := NewGenerator(max)
	fmt.Println("startpoint >>> ", g.startPoint)
	idLookupMap := make(map[string]bool)
	for i := 0; i < max; i++ {
		id := g.GetID()
		if _, ok := idLookupMap[id]; ok {
			t.Fatal("id exists")
		} else {
			idLookupMap[id] = true
		}
	}
	// fmt.Println(idLookupMap)
	fmt.Println("----------------------------- \n \n")
	time.Sleep(3 * time.Second)
	fmt.Println("after sleep  startpoint >>> ", g.startPoint)
	time.Sleep(1 * time.Second)
	fmt.Println("id >>> ", g.GetID())
}
