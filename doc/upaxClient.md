<h1 class="appTop">Upax Client</h1>

A Upax client is an XLattice Node and so has a unique 20 or 32-byte nodeID;
two RSA keys, one for encryption/decryption and one for digital 
signatures; some number of listening Acceptors, which accept connections;
and may some local persistent storage - file space, the local file system.

A client posts data as files to a Upax cluster, a group of machines 
which cooperate to maintain a reliable distributed store.  The servers 
might be all local, or might be anywhere in the world.  The client 
should in general be able to post to any of them and be confident that
the data has been stored and replicated as soon as it has received an
acknowledgement from the server.

Before posting the data the client must calculate the SHA3-256 hash 
of the data, its content hash.  The server will reply with its hash
of the data actually received.  If this differs from the client's hash,
the data will not have been stored; if they are the same, the data will
have been stored and replicated.

When posting data the client will also supply its nodeID and may 
specify a POSIX path.  The server will record the content hash, 
nodeID, a timestamp, and the path as metadata.

The data posted belongs to the client in the sense that at some point
in the future the client may delete it.  Otherwise there is no 
guarantee of privacy: any other client may fetch a copy of the data.

Clients retrieve data from the network by issuing a GET with a 
content key.  If the data is present, its metadata plus the data file
will be returned.

A client may also retrieve data from Upax mirrors.  These maintain
a copy of the Upax cluster's commit log and copies of at least some
of the files stored in the Upax cluster.
