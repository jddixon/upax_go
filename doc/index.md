# upax_go

A distributed content-keyed data storage system written in the **Go**
programming language.

Currently the implementation has three parts: 

* [server](upaxServer.html),
* [client](upaxClient.html), and 
* [fault-tolerant log](ftLog.html).  

All of these are written
in the Go programming language and will make use of the 
[xlattice_go](https://github.com/jddixon/xlattice_go) communications library.
(More extensive information on XLattice, although somewhat dated, 
is available at the [XLattice website](http://www.xlattice.org).)

A later implementation of upax_go 
is expected to add a [UpaxMirror](upaxMirror.html)
which will be capable of dealing with client queries but not able to add
to the data store.  That is, clients can use the UpaxMirror as a fast
local cache, but to add files to the distributed storage system 
they will have to go through a UpaxServer instead.

Upax accepts and stores data files of arbitrary size.  These are
identified by their SHA3-256 content key instead of a file name. 
This is the 256-bit version of Secure Hash Algorithm 3, also known
as Keccak.  The content key is the hash of the contents of the file.

Early versions of Upax
will store a copy of a file on each participating server node.
Later versions will store fewer copies, probably at least three,
deleting the least recently used where storage capacity limits are neared
and in effect migrating data to machines where it is more frequently
used.

Upax servers and clients are XLattice nodes and so are identified by 
unique 256 bit keys, NodeIDs.  Each node has a pair of RSA keys, one used 
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
