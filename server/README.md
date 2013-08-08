Upax Server
===========

A Upax server is an XLattice Node and so has a unique 32-byte NodeID;
two RSA keys, one for encryption/decryption and one for digital 
signatures; some number of listening Acceptors, which accept connections;
and some local persistent storage - file space, the local file system.

A number of Upax servers form a cluster which cooperates in maintaining
a distributed store.  In early implementations of Upax, the store will
be identical on all servers, or nearly so.  Although all servers must
have identical commitment logs, some delay will be allowed before all
servers have a copy of a file.

