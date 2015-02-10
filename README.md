# `free-slots` [![Build Status](https://travis-ci.org/rmcgibbo/free-slots.svg?branch=master)](https://travis-ci.org/rmcgibbo/free-slots) [![GitHub release](https://img.shields.io/github/release/rmcgibbo/free-slots.svg)](https://github.com/rmcgibbo/free-slots/releases)
*Report the number of available compute nodes for SLURM, TORQUE, and SGE cluster schedulers.*


This is a little command line script that reports the number of available slots (e.g. free
processors) on the compute nodes on a cluster. It works with either the SLURM, TORQUE/MAUI,
or SGE schedulers.

# Installation
Download the binary from the [github release](https://github.com/rmcgibbo/free-slots/releases) page,
unpack it, and drop it somewhere in your `$PATH`.

# Example
```
$ free-slots
Summary of PBS nodes with free slots

Num Nodes    Properties    Utilization    Free slots
---------    ----------    -----------    ----------
4            GPU           0/1            1
3            SP,MP         2/16           14
5            SP,MP         3/16           13
14           SP,MP         4/16           12
8            SP,MP         5/16           11
5            SP,MP         6/16           10
3            SP,MP         7/16           9
3            SP,MP         8/16           8
4            SP,MP         9/16           7
10           SP,MP         10/16          6
1            SP,MP         11/16          5
1            SP,MP         12/16          4
3            SP,MP         14/16          2
2            SP,MP         15/16          1
```
