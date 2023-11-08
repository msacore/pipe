<!-- Badges -->

[![Go Reference](https://pkg.go.dev/badge/github.com/msacore/pipe.svg)](https://pkg.go.dev/github.com/msacore/pipe)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/msacore/pipe)](go.mod)
[![License MIT](https://img.shields.io/github/license/msacore/pipe)](LICENSE)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/msacore/pipe)](https://github.com/msacore/pipe/releases)
[![Go Report](https://goreportcard.com/badge/github.com/msacore/pipe)](https://goreportcard.com/report/github.com/msacore/pipe)
[![codecov](https://codecov.io/gh/msacore/pipe/branch/main/graph/badge.svg?token=E8OCREETC0)](https://codecov.io/gh/msacore/pipe)

<!-- Inner Badges Links -->

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

<!-- README -->

# pipe

> **Note**
> Please, support this repository with a ⭐️ if you like the concept and would like to support the developer's motivation

> **Warning**  
> This module under rapid development

Power of Go channels with io.Pipe usability.
Build multithread tools easily.

- Thread safe
- `io` and `lo` like syntax (Tee, Reduce, Map, etc) but concurrently

## TODO

| Function | Impl | Tests | Doc Comments | Doc Readme |
|:---------|:----:|:-----:|:------------:|:----------:|
| Map |✅|✅|✅|✅|
| Filter |✅|✅|✅|✅|
| Split |✅|✅|✅|✅|
| ForEach | | | | |
| Spread | | | | |
| Join | | | | |
| Merge | | | | |
| Route | | | | |
| Replicate | | | | |
| Reduce | | | | |
| Wait |✅|✅|✅|✅|

## :arrow_down_small: Installation

This module powered by GO111MODULE and generics feature.
So it supports Go 1.18 and upper.

```bash
go get -u github.com/msacore/pipe
```

## :joystick: Examples

<details>
  <summary>Welcome example</summary>

[Open in playground](https://play.golang.com/p/9hQidaDDBTF)

```go
package main

import (
  "fmt"
  "github.com/msacore/pipe"
)

func main() {
  // Initial inputs
  nums := make(chan int, 4)

  // Generator
  go func() {
    for i := 0; i < 8; i++ {
      nums <- i
    }
    close(nums)
  }()

  // Processor
  filtered := pipe.Filter(func (value int) bool {
    return value % 2 == 0
  }, nums)
  strs := pipe.Map(func(value int) string {
    return fmt.Sprintf("%d", value)
  }, filtered)

  // Consumer
  for str := range strs {
    fmt.Println(str)
  }
}
```

</details>

## :building_construction: Methods

### [Map](map.go)

![Map](assets/methods/map.svg)

[![Parallel]](#parallel)
[![Sync]](#sync)
[![Sequential]](#sequential)
[![Single]](#single)
[![Same]](#same)

Take message and convert it into another type by map function.
If input channel is closed then output channel is closed.
Creates a new channel with the same capacity as input.

<details> 
  <summary>Usage examples</summary>

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

### [Filter](filter.go)

![Filter](assets/methods/filter.svg)

[![Parallel]](#parallel)
[![Sync]](#sync)
[![Sequential]](#sequential)
[![Single]](#single)
[![Same]](#same)

Take message and forward it if filter function return positive.
If input channel is closed then output channel is closed.
Creates a new channel with the same capacity as input.

<details> 
  <summary>Usage examples</summary>

```go
// input := make(chan int, 4) with random values.
// Say, the input contains [1, 2, 3, 4]

// Parallel strategy
// Best performance (Multiple goroutines)

output := Filter(func(value int) bool {
  fmt.Print(value)
    return value % 2 == 0
}, input)
// stdout: 4 1 2 3
// output: [4 2]

// Sync strategy
// Consistent ordering (Multiple goroutines with sequential output)

output := FilterSync(func(value int) bool {
  fmt.Print(value)
    return value % 2 == 0
}, input)
// stdout: 4 1 2 3
// output: [2 4]

// Sequential strategy
// Preventing thread race (Single goroutine)

output := FilterSequential(func(value int) bool {
  fmt.Print(value)
    return value % 2 == 0
}, input)
// stdout: 1 2 3 4
// output: [2 4]
```

</details>

### [Split](split.go)

![Split](assets/methods/split.svg)

[![Parallel]](#parallel)
[![Sync]](#sync)
[![Sequential]](#sequential)
[![Single]](#single)
[![Same]](#same)

Split takes a number of output channels and input channel, and forwards the input
messages to all output channels.
There is no guarantee that the message will be sent to the output channels in the
sequence in which they are provided.
If input channel is closed then all output channels are closed.
Creates new channels with the same capacity as input.

<details> 
  <summary>Usage examples</summary>

```go
// input := make(chan int, 4) with random values.
// Say, the input contains [1, 2, 3, 4]

// Parallel strategy
// Best performance (Multiple goroutines)

outs := Split(2, input)
// The gaps demonstrate uneven recording in the channels
// outs[0]: [2,    1, 3   ]
// outs[1]: [   1, 3,    2]

// Sync strategy
// Consistent ordering (Multiple goroutines with sequential output)

outs := SplitSync(2, input)
// The gaps demonstrate uneven recording in the channels
// outs[0]: [1,    2, 3   ]
// outs[1]: [   1, 2,    3]

// Sequential strategy
// Preventing thread race (Single goroutine)

outs := SplitSequential(2, input)
// The gaps demonstrate uneven recording in the channels
// outs[0]: [1,    2,    3   ]
// outs[1]: [   1,    2,    3]

// Also we have several shortcut functions like:

out1, out2 := Split2(input)
out1, out2, out3 := Split3(input)
```

</details>

### [Wait](wait.go)

Here are 3 helper functions that are waiting for the channels to close.
Each function blocks current goroutine until channels closing condition won't be done.

`Wait(in chan T) chan struct{}` - Waits for the input channel is closed and sends a signal to the returned channel.

<details> 
  <summary>Usage examples</summary>

```go
<-Wait(input1)
select {
  case <-Wait(input2):
  case <-Wait(input3):
}
// Will executed after input1 closed and input2 or input3 closed
```

</details>

`WaitAll(in ...chan T) chan struct{}` - Waits for all input channels are closed, and sends a signal to the returned channel.

<details> 
  <summary>Usage examples</summary>

```go
<-WaitAll(input1, input2)
// Will executed after input1 AND input2 closed

// It's equal:
<-Wait(input1)
<-Wait(input2)
```

</details>

`WaitAny(in ...chan T) chan struct{}` - Waits for one of the input channels are closes, and sends a signal to the returned channel. All other channels are read to the end in the background.

<details> 
  <summary>Usage examples</summary>

```go
<-WaitAny(input1, input2)
// Will executed after input1 OR input2 closed

// It's equal:
select {
  case <-Wait(input1):
  case <-Wait(input2):
}
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

## Links

- [🐞 Bug report](https://github.com/msacore/pipe/issues/new?assignees=jkulvich&labels=bug&projects=&template=%F0%9F%90%9E-bug-report.md&title=%5BBUG%5D)
- [⭐️ Feature request](https://github.com/msacore/pipe/issues/new?assignees=jkulvich&labels=enhancement&projects=&template=%E2%AD%90%EF%B8%8F-feature-request.md&title=%E2%AD%90%EF%B8%8F+%5BFEATURE%5D%3A+)
