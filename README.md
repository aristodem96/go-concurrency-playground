# Go Concurrency Playground

A learning project for practicing concurrency in Go: goroutines, channels, `sync.WaitGroup`, `context`, and worker pools.

## What it does

The program simulates a small "village" where several groups of workers run in parallel:

```
Miners (3) ──coal──► Blacksmiths (5) ──tools──► main
Postmen (3) ──mail─────────────────────────────► main
```

- **Miners** (`miner`) mine coal every second. Miner `n` has a power of `n * 10`.
- **Blacksmiths** (`blacksmith`) receive coal and forge a random tool for every 10 units.
- **Postmen** (`postman`) deliver a letter every second.

After 5 seconds, `context.WithTimeout` stops the miners and postmen. The blacksmiths finish forging with the coal they already have, and then the program prints a summary: coal mined, mail delivered, tools made, and coal left over.

## Concepts practiced

- Worker pools: the output channel is closed after `wg.Wait()`, once every worker has finished.
- Pipelines: the output of one pool (miners) feeds the input of another (blacksmiths).
- Graceful shutdown with `context`, including a `select` on `ctx.Done()` when sending to a channel, so goroutines never block forever.
- Directional channels (`chan<-`, `<-chan`) and `sync.WaitGroup.Go` (Go 1.25+).

## Running

Requires Go 1.26+.

```sh
git clone https://github.com/aristodem96/go-concurrency-playground.git
cd go-concurrency-playground
go run .
```

Check for data races:

```sh
go run -race .
```

## Project structure

```
.
├── main.go                  # starts the pools and collects results
├── miner/miner.go           # miner pool
├── blacksmith/blacksmith.go # blacksmith pool
└── postman/postman.go       # postman pool
```
