# GOQ: IOTA Qubic library for Go

This repository contains Go code for working with IOTA Qupla, a QUbic Programming Language 
as it is defined in the [reference implementation](https://github.com/iotaledger/qupla).
 
Note 1: GOQ is work in progress therefore bugs and unexpected behavior is highly probable. 
Please contact author at @lunfardo in IOTA Discord._

Note 2: I made efforts GOQ to be compatible with reference Abra spec and Qupla implementation.
All Qupla tests pass.
However, sometimes behavior may be different from what is expected, with execution of 
_eval_ and _test_ statements in particular.


Repository contains following relatively independent packages:

- [readyaml](https://github.com/lunfardo314/goq/tree/master/readyaml) library allows 
to read YAML representation of the Qupla module into static Go structure without much parsing. 
Qupla YAML representation contains everything necessary to interpret the module.
YAML representation of any Qupla module can be created by running [reference Qupla translator](https://github.com/iotaledger/qupla) 
with _-yaml_ flag. Examples how to use this package in Go and how to load YAML file 
into Python program can be found in [examples/readyaml](https://github.com/lunfardo314/goq/tree/master/examples/readyaml).

- [analyzeyaml](https://github.com/lunfardo314/goq/tree/master/analyzeyaml) library to 
convert module representation into interpretable Qupla representation which is completely independent from 
YAML source. It also performs necessary semantic analysis and checking.

- [qupla](https://github.com/lunfardo314/goq/tree/master/qupla) library contains 
Qupla runtime representations and Qupla interpreter

- [supervisor](https://github.com/lunfardo314/goq/tree/master/supervisor) contains Qubic 
Supervisor how it is defined in _Qubic Computational Model_ (QCM). 
Supervisor library is completely independent from Abra/Qupla implementation. 
Supervisor interacts with _entities_ using abstract _EntityCore_ interface. 
_Entity_ can be Qupla function interpreter or any other software agent able to calculate 
trit vector output or null value from trit vector input.



