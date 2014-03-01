<h1 class="appTop">Upax Mirror</h1>

A Upax mirror is an XLattice Node and so has a unique 20 or 32-byte nodeID;
two RSA keys, one for encryption/decryption and one for digital 
signatures; some number of listening Acceptors, which accept connections;
and may some local persistent storage - file space, the local file system.

A Upax mirror differs from a Upax server in that it only accepts 
requests for copies of data.  It may not be used by a client to store
data in the Upax cluster.

Also, while Upax servers use voting to maintain a consensus on what
data is stored in the cluster, the mirror takes no part in these 
decisions: it simply records the final decision and and commits it to
its log.

Although the mirror must maintain an accurate and up-to-date copy of the 
Upax commit log, the mirror may or may not maintain a complete copy of the 
data store.  If a data file requested is not present in the mirror's 
store, it will fetch a copy from another mirror or a server and then
return that copy to the client.  The effect of this is that data files
will tend to migrate to where they are used.

