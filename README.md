# goq
GOQ: IOTA Qubic library for Go

This repository contains Go code for working with IOTA Qupla, a QUbic Programming Language 
as it is defined in the [reference implementation](https://github.com/iotaledger/qupla).
 
Note 1: GOQ is work in progress therefore bugs and unexpected behavior is highly probable. 
Apologies in advance. Please contact author at @lunfardo in IOTA Discord.

Note 2: I made efforts GOQ to be compatible with reference Abra spec and Qupla implementation.
All Qupla tests pass.
However, sometimes behaviour may be different from what is expected, with execution of _eval_ and _test_ statements in particular.


Repository contains following relatively independent packages:

- [readyaml](https://github.com/lunfardo314/goq/tree/master/readyaml) library allows 
to read YAML representation of the Qupla module in Go. YAML representation of any Qupla 
module can be created by running reference [Qupla translator](https://github.com/iotaledger/qupla) 
with _-yaml_ flag. Examples how to use this package in Go and how to load YAML file 
into Python program can be found in [examples/readyaml](https://github.com/lunfardo314/goq/tree/master/examples/readyaml).

- [analyzeyaml](https://github.com/lunfardo314/goq/tree/master/analyzeyaml) library to 
convert module representation into interpretable representation completely independent from 
YAML source.

- [qupla](https://github.com/lunfardo314/goq/tree/master/qupla) library contains 
Qupla runtime representations and Qupla interpreter

- [supervisor](https://github.com/lunfardo314/goq/tree/master/supervisor) contains Qubic 
supervisor which implements QCM, _Qubic Computational Model_. 
Supervisor library is completely independent from Abra/Qupla implementation. 
Supervisor interacts with _entities_ using abstract _EntityCore_ interface. 
_Entity_ can be Qupla function interpreter ro any other software agent able to calculate 
trit vector output or null value from trit vector input.



