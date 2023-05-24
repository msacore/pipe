# pipe

Power of Go channels with io.Pipe usability.
Build multithread tools easily.

- Thread safe
- `io` and `lo` like syntax (Tee, Reduce, Map, etc) but concurrently

## Pipeline Methods

### Map

![Map](assets/methods/map.svg)

> Take message and convert it into another type by map function.
> If input channel is closed then output channel is closed.

- `Input` single of T1
- `Output` single of T2
- `Capacity` Same as the input
- `Closing` Input closed

### Filter

![Filter](assets/methods/filter.svg)

> Take message and forward it if filter function return positive.
> If input channel is closed then output channel is closed.

- `Input` single of T
- `Output` single of T
- `Capacity` Same as the input
- `Closing` Input closed

### Split

![Split](assets/methods/split.svg)

> Take next message and forward to all output channels.
> If input channel is closed then all output channels are closed.

- `Input` single of T
- `Output` multiple of T
- `Capacity` Same as the input
- `Closing` Input closed

### Spread

![Spread](assets/methods/spread.svg)

> Take next message and forward it to next output channel.
> If input channel is closed then all output channels are closed.
> Randomization algorithm is Round Robin or random.

- `Input` single of T
- `Output` multiple of T
- `Capacity` Same as the input
- `Closing` Input closed

### Join

![Join](assets/methods/join.svg)

> Take next available message from any input and forward it to output.
> If all input channels are closed then output channel is closed.

- `Input` multiple of T
- `Output` single of T
- `Capacity` Sum of inputs
- `Closing` All input channels closed

### Merge

![Merge](assets/methods/merge.svg)

> Take next message from all channels (wait for data) and send new message into output.
> If one of input channels is closed then output channel is closed.
> All other input channels will be read till end in background.

- `Input` multiple of T1
- `Output` single of T2
- `Capacity` Minimal of inputs
- `Closing` Any of input channels closed

### Route

![Route](assets/methods/route.svg)

> Take next message from input and forward it to one of output channels by route function.
> If input channel is closed then all output channels are closed.

- `Input` single of T
- `Output` multiple of T
- `Capacity` Same as the input
- `Closing` Input closed

### Replicate

![Replicate](assets/methods/replicate.svg)

> Take next message from input and forward copies to output.
> If input channel is closed then all output channels are closed.

- `Input` single of T
- `Output` single of T
- `Capacity` Same as the input
- `Closing` Input closed

### Reduce

![Reduce](assets/methods/reduce.svg)

> Take several next messages from input and send new message to output.
> > If input channel is closed then all output channels are closed.

- `Input` single of T1
- `Output` single of T2
- `Capacity` Same as the input
- `Closing` Input closed
