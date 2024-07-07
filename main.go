package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
)

type Object struct {
	Id          int
	ObjectName  string
	ObjectDesc  string
	ObjectValue int
}

type ParentObject struct {
	Id         int
	ObjectName string
	ObjectDesc string
	Objects    Object
}

func createObject(length int) []Object {
	var arrObj []Object
	for i := 0; i < length; i++ {
		objA := Object{
			Id:          i + 1,
			ObjectName:  fmt.Sprintf("objA-%d", i+1),
			ObjectDesc:  fmt.Sprintf("this ini objA-%d", i+1),
			ObjectValue: rand.Intn(1000),
		}
		arrObj = append(arrObj, objA)
	}
	return arrObj
}

func sequence() {
	var arrObj []Object
	var arrObj2 []Object
	var arrObj3 []Object

	// Create arrObj 1
	arrObj = createObject(1000000)
	PrintMemUsage()

	// Create arrObj 2
	arrObj2 = createObject(1000000)
	PrintMemUsage()

	// Create arrObj 3
	arrObj3 = createObject(1000000)
	PrintMemUsage()

	fmt.Println(len(arrObj))

	fmt.Println(len(arrObj2))

	fmt.Println(len(arrObj3))
}

func parallel() {
	chanObject := make(chan []Object, 3)
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			arrObj := createObject(1000000)
			PrintMemUsage()
			chanObject <- arrObj
			wg.Done()
		}()
	}
	wg.Wait()
	close(chanObject)
	for obj := range chanObject {
		fmt.Println(len(obj))
	}
}

func parallelDevided() {
	arrQueues := [2]int{1, 2}
	var wg sync.WaitGroup

	for _, q := range arrQueues {
		chanObject := make(chan []Object, q)
		for i := 0; i < q; i++ {
			wg.Add(1)
			go func() {
				arrObj := createObject(1000000)
				PrintMemUsage()
				chanObject <- arrObj
				wg.Done()
			}()
		}
		wg.Wait()
		close(chanObject)
		for obj := range chanObject {
			fmt.Println(len(obj))
		}
	}

}

func main() {
	parallelDevided()
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("HeapAlloc = %v MiB", bToMb(m.HeapAlloc))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
