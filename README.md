[Parallel]: assets/strategies/parallel.svg
[Sync]: assets/strategies/sync.svg
[Sequential]: assets/strategies/sequential.svg

[Single]: assets/strategies/single.svg
[All]: assets/strategies/all.svg
[Any]: assets/strategies/any.svg

[Same]: assets/strategies/same.svg
[Mult]: assets/strategies/mult.svg
[Min]: assets/strategies/min.svg
[Max]: assets/strategies/max.svg
[Sum]: assets/strategies/sum.svg

# pipe

Power of Go channels with io.Pipe usability.
Build multithread tools easily.

- Thread safe
- `io` and `lo` like syntax (Tee, Reduce, Map, etc) but concurrently

## Pipeline Methods

### [Map](map.go)

![Map](assets/methods/map.svg)

![Parallel] ![Sync] ![Sequential] ![Single] ![Same]

Take message and convert it into another type by map function.
If input channel is closed then output channel is closed.
Creates a new channel with the same capacity as input.

<details> 
  <summary>Usage example</summary>

```go
// input := make(chan int, 4) with random values.
// Say, the input contains [1, 2, 3]

// Parallel strategy
// Best performance (Multiple goroutines)

output := Map(input, func(value int) string { 
    fmt.Print(value)
    return fmt.Sprintf("val: %d", value) 
})
// stdout: 2 1 3
// output: ["val: 2", "val: 1", "val: 3"] 

// Sync strategy
// Consistent ordering (Multiple goroutines with sequential output)

output := MapSync(input, func(value int) string { 
    fmt.Print(value)
    return fmt.Sprintf("val: %d", value) 
})
// stdout: 2 1 3
// output: ["val: 1", "val: 2", "val: 3"] 

// Sequential strategy
// Preventing thread race (Single goroutine)

output := MapSequential(input, func(value int) string { 
    fmt.Print(value)
    return fmt.Sprintf("val: %d", value) 
})
// stdout: 1 2 3
// output: ["val: 1", "val: 2", "val: 3"] 
```

</details>

### Filter

> **Warning**
> This function under construction

![Filter](assets/methods/filter.svg)

![Parallel] ![Sync] ![Sequential] ![Single] ![Same]

Take message and forward it if filter function return positive.
If input channel is closed then output channel is closed.
Creates a new channel with the same capacity as input.

### Split

> **Warning**
> This function under construction

![Split](assets/methods/split.svg)

![Sequential] ![Single] ![Same]

Take next message and forward to all output channels.
If input channel is closed then all output channels are closed.
Creates new channels with the same capacity as input.

### Spread

> **Warning**
> This function under construction

![Spread](assets/methods/spread.svg)

![Sequential] ![Single] ![Same]

Take next message and forward it to next output channel.
If input channel is closed then all output channels are closed.
Randomization algorithm is `Round Robin` or `random`.
Creates new channels with the same capacity as input.

### Join

> **Warning**
> This function under construction

![Join](assets/methods/join.svg)

![Sequential] ![All] ![Sum]

Take next available message from any input and forward it to output.
If all input channels are closed then output channel is closed.
Creates new channel with sum of capacities of input channels.

### Merge

> **Warning**
> This function under construction

![Merge](assets/methods/merge.svg)

![Parallel] ![Sync] ![Sequential] ![Any] ![Min]

Take next message from all channels (wait for data) and send new message into output.
If one of input channels is closed then output channel is closed.
All other input channels will be read till end in background.
Creates new channel with minimal capacity of input channels.

### Route

> **Warning**
> This function under construction

![Route](assets/methods/route.svg)

![Parallel] ![Sync] ![Sequential] ![Single] ![Same]

Take next message from input and forward it to one of output channels by route function.
If input channel is closed then all output channels are closed.
Creates new channels with the same capacity as input.

### Replicate

> **Warning**
> This function under construction

![Replicate](assets/methods/replicate.svg)

![Sequential] ![Single] ![Mult]

Take next message from input and forward copies to output.
If input channel is closed then all output channels are closed.
Creates new channel with the same capacity as input multiplied by N.

### Reduce

> **Warning**
> This function under construction

![Reduce](assets/methods/reduce.svg)

![Sequential] ![Single] ![Same]

Take several next messages from input and send new message to output.
If input channel is closed then all output channels are closed.
Creates new channel with the same capacity as input.
