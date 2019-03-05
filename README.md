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
| `vns.iterMax`    | Maximum iteretion of the overall VNS                                    |
| `construction.initStrategy`  | Strategy used to create the first posibly unfeasible solution. Available choices are `random` and `greedy` |
| `construction.levelMax`  | Maximum level of perturbation in constraction part               |
| `construction.iterMax`   | Maximum iteretion in construction part                           |
| `optimization.strategy`  | Strategy used in optimization part. Available choices are  `2opt` and  `lexical2opt` |
| `optimization.levelMax`  | Maximum level of perturbation in optimization part                 |
| `optimization.iterMax`  | Maximum iteretion in optimization part                              |

Example config:

```yaml
vns:
  iterMax: 2
construction:
  strategy: random
  levelMax: 10
  penaltyWeight:
    timeWindow: 100
    capacity: 1
    pickupDelivery: 1
optimization:
  strategy: 2opt
  iterMax: 10
  levelMax: 10
```

# Benchmarks
