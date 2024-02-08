package phaul

import (
	"github.com/checkpoint-restore/go-criu/v7"
)

// Config is the configuration which is passed around
//
// Pid is what we migrate
// Memfd is the file descriptor via which criu can transfer memory pages.
// Wdir is the directory where phaul can put images and other stuff
type Config struct {
	Pid   int32
	Memfd int
	Wdir  string
}

// Local interface
// Interface to local classes. Client calls them when it needs something on the source node.
//
// Methods:
//
//   - DumpCopyRestore() is called on client side when the
//     pre-iterations are over and it's time to do full dump,
//     copy images and restore them on the server side.
//     All the time this method is executed victim tree is
//     frozen on client. Returning nil kills the tree, error
//     unfreezes it and resumes. The criu argument is the
//     pointer on created criu.Criu object on which client
//     may call Dump(). The requirement on opts passed are:
//     set Ps.Fd to comm.Memfd
//     set ParentImg to lastClientImagesPath
//     set TrackMem to true
type Local struct {
}

func (li *Local) DumpCopyRestore(criu *criu.Criu, c Config, lastClientImagesPath string) error {
	return nil
}
