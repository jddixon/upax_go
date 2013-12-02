package upax_go

// upax_go/ihave_obj.go

// Each ID is a content key for an entry that the peer claims to have
// in its store.
//
type IHaveObj struct {
	IDs [][]byte
}
