# robin
Golang Weighted Round Robin implementation

The package implements the basic weighted round-robin balancing algorithm, it's go-routine safe and designed for concurrent use. 
See details in [Wikipedia](https://en.wikipedia.org/wiki/Weighted_round_robin).


### Classical WRR
In classical WRR the scheduler cycles over the queues.
When a queue q-i is selected, the scheduler will send packets, up to the emission of w-i packet or the end of queue.

```pseudo
Constant and variables: 
    const N             // Nb of queues 
    const weight[1..N]  // weight of each queue
    queues[1..N]        // queues
    i                   // queue index
    c                   // packet counter
    
Instructions:
while true do 
    for i in 1 .. N do
        c := 0
        while (not queue[i].empty) and (c<weight[i]) do
            send(queue[i].head())
            queue[i].dequeue()
            c:= c+1
```

```golang
package main

import "fmt"

func main() {
	wrr := &WRR[string]{}
  wrr.Add("Beer", 5)
  wrr.Add("Vodka", 3)
  wrr.Add("Vine", 1)

	for i := 0; i < 10; i++ {
		fmt.Printf("%s ", wrr.Next())
	}
}
```

In future, the package will also provide Interleaved WRR algorithm implementation.
