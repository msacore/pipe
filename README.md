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

> **Warning**  
> This module under rapid development

Power of Go channels with io.Pipe usability.
Build multithread tools easily.

- Thread safe
- `io` and `lo` like syntax (Tee, Reduce, Map, etc) but concurrently

## :building_construction: Methods

### [Map](map.go)

![Map](assets/methods/map.svg)

:arrows_counterclockwise:
[![Parallel]](#parallel)
[![Sync]](#sync)
[![Sequential]](#sequential)
:lock:
[![Single]](#single)
:package:
[![Same]](#same)

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

output := Map(func(value int) string { 
    fmt.Print(value)
    return fmt.Sprintf("val: %d", value) 
}, input)
// stdout: 2 1 3
// output: ["val: 2", "val: 1", "val: 3"] 

// Sync strategy
// Consistent ordering (Multiple goroutines with sequential output)

output := MapSync(func(value int) string { 
    fmt.Print(value)
    return fmt.Sprintf("val: %d", value) 
}, input)
// stdout: 2 1 3
// output: ["val: 1", "val: 2", "val: 3"] 

// Sequential strategy
// Preventing thread race (Single goroutine)

output := MapSequential(func(value int) string { 
    fmt.Print(value)
    return fmt.Sprintf("val: %d", value) 
}, input)
// stdout: 1 2 3
// output: ["val: 1", "val: 2", "val: 3"] 
```

</details>

## :gear: Strategies

Each function has own set of strategies from all categories.
It describes how your channel data `processing`, when channels `closing`, and how calculate `capacity` of output channels.

### :arrows_counterclockwise: Processing

Some functions have different channel processing algorithms. To ensure maximum performance, it is recommended to use the original function. However, specific algorithms can help in cases where you are faced with a race of threads or you need to output data strictly in the same order in which you received them.

#### Parallel

![Parallel]  
Each handler is executed in its own goroutine and there is no guarantee that the output order will be consistent. Recommended for best performance.

#### Sync

![Sync]  
Each handler executes in its own goroutine, but the result of the youngest goroutine waits for the oldest goroutine to finish before being passed to the output stream. To prevent memory leaks, the strategy will wait if there is more waiting data than the capacity of the output channel. Recommended if you want to get the output data in the same order as the input data.

#### Sequential

![Sequential]  
Each handler is executed sequentially, one after the other. Keeps the order of the output data equal to the order of the input data. Recommended if it is necessary to exclude the race of threads between handlers.

### :lock: Closing

Each function has one of several strategies for closing output channels. Understanding will help you understand how and when your pipeline closes.

#### Single

![Single]  
Suitable only for functions with one input. If the input channel is closed, then the output channels are closed.

#### All

![All]  
If all input channels are closed, then the output channels are closed.

#### Any

![Any]  
If one of the input channels is closed, the output channels are closed. All other channels will be read to the end in the background.

### :package: Capacity

Each function creates new output channels with the capacity corresponding to a specific strategy.

#### Same

![Same]  
Suitable only for functions with one input channel. The output channels will have a capacity equal to the input channel.

#### Mult

![Mult] 
Suitable only for functions with one input channel. The output channels will have a capacity equal to the input channel multiplied by N.

#### Min

![Min]  
The output channels will have a capacity equal to the minimum capacity of the input channels.

#### Max

![Max]  
The output channels will have a capacity equal to the maximum capacity of the input channels.

#### Sum

![Sum]  
The output channels will have a capacity equal to the sum of capacities of the input channels.

## === DRAFT ===

<details> 
  <summary><b>Under Construction</b></summary>

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

</details>
