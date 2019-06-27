#####**New:** Qubic Game Of Life on GOQ: [gol.tanglebeat.com](http://gol.tanglebeat.com)

# GOQ: Qupla library for Go

This repository contains Go code for working with IOTA Qupla, a QUbic Programming Language 
as it is defined in the [reference implementation](https://github.com/iotaledger/qupla).
 
Note 1: GOQ is work in progress therefore bugs and unexpected behavior is highly probable. 
Please contact author as @lunfardo in IOTA Discord.

Note 2: I made efforts GOQ to be compatible with reference Abra spec and Qupla implementation.
All Qupla tests pass.
However, sometimes behavior may be different from what is expected, with execution of 
_eval_ and _test_ statements in particular.

Repository contains the following packages:

- [readyaml](https://github.com/lunfardo314/goq/tree/master/readyaml) library allows 
to read YAML representation of the Qupla module into static Go structures without much parsing. 
YAML representation of any Qupla module can be created by running [reference Qupla translator](https://github.com/iotaledger/qupla) 
with _-yaml_ flag. Examples how to use this package in Go and how to load YAML file 
into Python program can be found in [examples/readyaml](https://github.com/lunfardo314/goq/tree/master/examples/readyaml).
It also contains YAML representations of `QuplaTests`, `Examples`, `Fibonacci` and `Curl` modules.

- [analyzeyaml](https://github.com/lunfardo314/goq/tree/master/analyzeyaml) library to 
convert YAML module representation into interpretable Qupla representation which is completely independent from 
YAML source. It also performs necessary semantic analysis and checking.

- [qupla](https://github.com/lunfardo314/goq/tree/master/qupla) package contains 
Qupla runtime representations and Qupla interpreter

- [optimize](https://github.com/lunfardo314/goq/tree/master/qupla) package contains 
Qupla optimization functions.  

- [supervisor](https://github.com/lunfardo314/goq/tree/master/supervisor) contains Qubic 
Supervisor as it is defined in _Qubic Computational Model_ (QCM). 
Supervisor is completely independent from Qupla/Abra. 
It interacts with _entities_ using abstract `EntityCore` interface. 
_Entity_ can be Qupla function with interpreter or any other software agent able to calculate 
trit vector output (or null value) from trit vector input. 
Each _Entity_ can join (subscribe) the environment and it can affect environment in the sense of Abra specification.
Effects in the form of trit vectors can be sent to environment by placing message into the main input queue.
Each externally placed effect starts a _quant_. During a _quant_ all entities are running in parallel until _quant_ is finished.
Supervisor is locked for any changes during the quant.
Exported supervisor API can be found in the file `api.go`. 

- [examples/gol](https://github.com/lunfardo314/goq/tree/dev/examples/gol) contains **Game of Life**
proof-of-concept originaly created by Eric Hop in Qupla [reference implementation](https://github.com/iotaledger/qupla).
This PoC uses original unchanged Qupla code from reference implementation. It uses GOQ supervisor
and GolOracle written in Go. GOL qubics are running on the server and user can access the 
system from the web browser. 

- [examples/goq-cli](https://github.com/lunfardo314/goq/tree/dev/examples/goq-cli) contains 
_goq-cli_, a simple command line interface to Qupla and supervisor. 
Primary purpose of _goq-cli_ is testing of the library itself. It hopefully can be used to test and debug any Qupla modules.
Please find _goq-cli_ command reference below.

## goq-cli commands

- `verb` show verbosity level
- `verb <verbosity_level>` set verbosity level. 2 is default, 3 is for debugging, >5 is tracing
- `runtime` show memory usage
- `dir` show current directory
- `dir <directory>` set current directory
- `load <module yaml file>` load module form YAML file. Loading means reading module form YAML 
file, analyzing it and attaching to the supervisor by 'joining' and 'affecting' respective 
environments, referenced from functions.
- `save <file to save as yaml>` marshal module to YAML file (for echo testing)
- `lexe <filter substring>` numbered list of `eval` and `test` statements of the loaded module
- `lfun <filter substring>` list functions of the module, name of which contains substring.
- `lenv` list environments joined and/or affected by module's functions
- `trace [<filter substring> [<traceLevel>]]` set trace mode for all functions, names of which contains substring 
- `trace stop` stop tracing all functions 
- `run all` run all `test` and `eval` statements of the module
- `run <exec idx>` run specific statement by it's index in the numbered list
- `run <from exec idx>-<to exec idx>` run range of executable stataments
- `repeat <exec idx> <repeat times>` run specific executable statement number of times
- `post <effect decimal> <environment>` post effect to the environment

# Getting started

### Install Go
Follow the [instructions](https://golang.org/doc/install). 
Make sure to define `GOPATH` environment variable to the root where all your Go projects will land.
The `GOPATH` directory should contain at least `src` (for sources) and `bin` 
(for executable binaries) subdirectories. 
Set `PATH` to your `GOPATH/bin`

### Download GOQ

Run `go get github.com/lunfardo314/goq/examples/goq-cli`

### Run supervisor and Qupla tests

Make directory `GOPATH/src/lunfardo314/goq/tests` current.

Run test: `go test`

### Run goq-cli

Make directory `GOPATH/src/lunfardo314/goq/examples/goq-cli` current.

Run `go install`

Run `goq-cli`

`>>> load ../modules/QuplaTest.yml`

`>>> run all`



