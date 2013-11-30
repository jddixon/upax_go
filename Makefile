# by default, first target is only target
all:	c.pb.go s.pb.go s_aes_cnx.go c_aes_cnx.go \
	c_in_handler.go s_in_handler.go

c.pb.go: c.proto
	mkdir -p _pb
	protoc --go_out=_pb $<
	cat _pb/$@\
	|gofmt >$@
	rm -rf _pb

s.pb.go: s.proto
	mkdir -p _pb
	protoc --go_out=_pb $<
	cat _pb/$@\
	|gofmt >$@
	rm -rf _pb

c_aes_cnx.go: c_context aes_cnx.t
	xgoT -c c_context -E .go -p c_ aes_cnx 

s_aes_cnx.go: s_context aes_cnx.t
	xgoT -c s_context -E .go -p s_ aes_cnx 

c_in_handler.go: c_context in_handler.t
	xgoT -c c_context -E .go -p c_ in_handler 

s_in_handler.go: s_context in_handler.t
	xgoT -c s_context -E .go -p s_ in_handler 

