package main

import (
	"hash/fnv"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

// Number of philosophers is simply the length of this list.
var philosophers = []string{"John", "Alice", "Michael", "Sophia", "Robert"}

const hunger = 3                // Number of times each philosopher eats
const think = time.Second / 100 // Mean think time
const eat = time.Second / 100   // Mean eat time

var logger = log.New(os.Stdout, "", 0)

var dining sync.WaitGroup

func diningProblem(philosopherName string, dominantHand, otherHand *sync.Mutex) {
	logger.Println(philosopherName, "has been seated.")
	hash := fnv.New64a()
	hash.Write([]byte(philosopherName))
	rng := rand.New(rand.NewSource(int64(hash.Sum64())))
	randomSleep := func(t time.Duration) {
		time.Sleep(t/2 + time.Duration(rng.Int63n(int64(t))))
	}
	for h := hunger; h > 0; h-- {
		logger.Println(philosopherName, "is hungry.")
		dominantHand.Lock() // pick up forks
		otherHand.Lock()
		logger.Println(philosopherName, "is eating.")
		randomSleep(eat)
		dominantHand.Unlock() // put down forks
		otherHand.Unlock()
		logger.Println(philosopherName, "is thinking.")
		randomSleep(think)
	}
	logger.Println(philosopherName, "is satisfied.")
	dining.Done()
	logger.Println(philosopherName, "has left the table.")
}

func main() {
	logger.Println("The table is empty.")
	dining.Add(5)
	fork0 := &sync.Mutex{}
	forkLeft := fork0
	for i := 1; i < len(philosophers); i++ {
		forkRight := &sync.Mutex{}
		go diningProblem(philosophers[i], forkLeft, forkRight)
		forkLeft = forkRight
	}
	go diningProblem(philosophers[0], fork0, forkLeft)
	dining.Wait() // wait for philosophers to finish
	logger.Println("The table is empty.")
}
