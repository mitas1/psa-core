# PSA core

PSA core is basically a solver based on variable neighbourhood search. It consists of the
construction and optimization part where in both the VNS is applied. This repository is part of the
diploma thesis of the author.

## Building and running

Firstly, clone this repository to `${GOPATH}/src/github.com/mitas1/pps-core`. 

Then, update `config.yaml` if needed (see [Configuration](#configuration)).

To build the `psa-core` type the following:

```sh
$ make build
```

Finally, to run:

```sh
$ make run
```

To see available flags type:

```sh
$ ./psa-core --help
```

# Configuration

This table describes available configuration options:

| Name             | Description                                                             |
| ---------------- | ----------------------------------------------------------------------- |
| `common.iterMax`    | Maximum iteretion of the overall Algorithm                           |
| `common.maxTime`    | Maximum execution time in seconds                                    |
| `construction.strategy`  | Strategy used to create the first posibly unfeasible solution. Available choices are `random`, `greedy`, `sortedBydueDate` and `sortedByTW` |
| `construction.levelMax`  | Maximum level of perturbation in constraction part               |
| `construction.iterMax`   | Maximum iteretion in construction part                           |
| `construction.penalty.timeWindows`    | Weight of time windows penalty                      |
| `construction.penalty.pickupDelivery`    | Weight of pickup and delivery penalty            |
| `construction.penalty.capacity`    | Weight of capacity penalty                             |
| `optimization.objective`  | Objective function in optimization phase. Available choices are `span` and `time` |
| `optimization.asymetric`  | Whether the instance is asymetric or not                       |
| `optimization.gvns`       | If specified GVNS is used as optimzation phase                 |
| `optimization.gvns.levelMax`  | Maximum level of perturbation in optimization part                 |
| `optimization.gvns.iterMax`   | Maximum iteretion in optimization part                             |

Example config:

```yaml
common:
  iterMax: 2
  maxTime: 100
construction:
  strategy: random
  levelMax: 10
  penalty:
    timeWindows: 100
    pickupDelivery: 10
    capacity: 1
optimization:
  objective: span
  gvns:
    iterMax: 4
    levelMax: 40
```

# Benchmarks
