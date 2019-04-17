# PSA core

PSA core is basically a PDPTW solver based on 2-phase heurstic. It consists of the
construction and optimization phase. This repository is part of the diploma thesis of the author.

## Building and running

Firstly, update `config.yaml` if needed (see [Configuration](#configuration)).

To build the `psa-core` type the following:

```sh
$ make build
```

Finally, to run:

```sh
$ ./psa-core
```

To see available flags type:

```sh
$ ./psa-core --help
```

# Configuration

This table describes available configuration options:

| Name             | Description                                                             |
| ---------------- | ----------------------------------------------------------------------- |
| `common.iterMax`    | Maximum iteretion of the overall algorithm                           |
| `common.maxTime`    | Maximum execution time in seconds                                    |
| `construction.strategy`  | Strategy used to create the first posibly unfeasible solution. Available choices are `random`, `greedy`, `sortedBydueDate` and `sortedByTW` |
| `construction.levelMax`  | Maximum level of perturbation in constraction part               |
| `construction.iterMax`   | Maximum iteretion in construction part                           |
| `construction.penalty.timeWindows`    | Weight of time windows penalty                      |
| `construction.penalty.pickupDelivery`    | Weight of pickup and delivery penalty            |
| `construction.penalty.capacity`    | Weight of capacity penalty                             |
| `optimization.objective`  | Objective function in optimization phase. Available choices are `span` and `time` |
| `optimization.asymetric`  | Whether the instance is asymetric or not                       |
| `optimization.vns`       | If specified VNS is used as optimzation phase                 |
| `optimization.vns.levelMax`  | Maximum level of perturbation in optimization part                 |
| `optimization.vns.iterMax`   | Maximum iteretion in optimization part                             |
| `optimization.vns.localSearch` | Local search in vns applied. Available choices are `vnd`, `2opt` and `shifting` |
| `optimization.sa` | If `optimization.vns` is not specified, simmulated annealing is applied in optimization phase. |
| `optimization.sa.iterMax` | Maximum iteration of the annealing |

Example config:

```yaml
common:
  iterMax: 4
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
  vns:
    iterMax: 10
    levelMax: 50
    localSearch: vnd
```
