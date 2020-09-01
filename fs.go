package styx

import "github.com/jecoz/flexi/fs"

type FS struct {
	fs.RWFS
}

func (fys *FS) HandleT(t Request) {
	switch msg := t.(type) {
	case Topen:
		msg.Ropen(fys.Open(t.Path()))
	case Tstat:
		file, err := fys.Open(t.Path())
		if err != nil {
			msg.Rstat(nil, err)
			return
		}
		msg.Rstat(file.Stat())
	case Twalk:
		file, err := fys.Open(t.Path())
		if err != nil {
			msg.Rwalk(nil, err)
			return
		}
		msg.Rwalk(file.Stat())
	case Tcreate:
		msg.Rcreate(fys.Create(msg.Path(), msg.Mode))
	case Tremove:
		msg.Rremove(fys.Remove(msg.Path()))
	case Ttruncate:
		file, err := fys.Open(msg.Path())
		if err != nil {
			msg.Rtruncate(err)
			return
		}
		msg.Rtruncate(fs.Truncate(file, msg.Size))
	case Tutimes:
		// Each file can handle this information without
		// requiring the user telling when the file has
		// been modified.
		msg.Rutimes(nil)
	default:
		// Default responses will take
		// care of the remaining/new messages.
	}
}

func (fys *FS) Serve9P(s *Session) {
	for s.Next() {
		fys.HandleT(s.Request())
	}
}

func NewFS(p fs.RWFS) *FS { return &FS{p} }
