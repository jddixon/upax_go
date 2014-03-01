<h1 class="appTop">Upax Server</h1>

A Upax server is an XLattice Node and so has a unique 20 or 32-byte NodeID;
two RSA keys, one for encryption/decryption and one for digital 
signatures; some number of listening Acceptors, which accept connections;
and some local persistent storage - file space, the local file system.

A number of Upax servers form a cluster which cooperates in maintaining
a distributed store.  In early implementations of Upax, the store will
be identical on all servers, or nearly so.  Although all servers must
have identical commitment logs, some delay will be allowed before all
servers have a copy of a file.

The server's XLattice peers are either clients, other servers, or 
mirrors.  Clients can PUT files to the cluster through the server
and can GET files from the server.  If the server does not have a
requested file, it will get it from another server or a mirror.

The server synchronizes its log and its store with other servers.  The
log is synchronized using a consensus protocol.  All servers must have
identical commit logs.  Servers may exchange files using GET and PUT.


A mirror may request files from the server using GET and will receive
entries for its commit log from the servers.

If there are K <= 3 servers, then all servers must have a copy of a 
file before it is added to the commit log.  If there are K > 3 servers,
then there must be at least 3 copies of the file in the stores of the
servers and mirrors, unless cluster administrators have set a higher
level.  When Upax is in early implementation then K copies will be 
required.  That is, the number of copies must equal or exceed the 
number of servers.

