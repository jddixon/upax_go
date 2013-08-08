upax_go
=======

A distributed content-keyed data storage system.

In the short term, the implementation will have three parts: server,
client, and fault-tolerant log (ftlog).  All of these will be written
in Go and will make use of the xlattice_go communications library.

Upax accepts and stores data files of arbitrary size.  These are
identified by their SHA3-256 content key.  Early versions of Upax
will store a copy of a file on each participating server node.
Later versions will store fewer copies, probably at least three,
deleting the least recently used where storage capacity is neared
and in effect migrating data to machines where it is more frequently
used.

Upax servers and clients are XLattice nodes and so are identified by 
unique 256 bit keys, NodeIDs.  Each has a pair of RSA keys, one used 
for encryption and one used for digital signatures.  At least in the
near term Upax servers will only accept messages from known hosts.
Specifically, a client must prove that it has the private RSA key
corresponding to its public RSA key when opening any connection to
a server.

In the Python implementation of Upax, files are owned by the first
to commit them.  There is no notion of privacy.  Any client may 
request a copy of any file - but only the owner may delete it.  It
is likely that upax_go will retain these characteristics in the near
term.
