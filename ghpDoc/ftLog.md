ftlog
=====

The Upax log is an append-only file recording a series of commits.
Each commit line contains at least a timestamp, the SHA3-256 hash of
the file being committed, and the nodeID of the committer.  It may
also contain a POSIX path.  The commit line is CR/LF-terminated.
The hash is guaranteed to be unique to the file and the nodeID is
unique to the committer.  The path may not be unique.

This is a chained log, meaning that the first line in each log file
identifies the preceding log by its hash.  A copy of that log is
present in the data store; on production systems this will normally
be /var/Upax.

The expectation is that the Upax log will normally have two
identical copies, one in memory and one on disk.  The log specifies
the state of the data store.  All Upax servers in a cluster should
share the same log, in the sense that all copies of the log should
be identical.

Data stores may not be identical.  A minimum number of servers
(three if there are at least that many servers, all if fewer) must
have a copy of a data file before it is committed.
