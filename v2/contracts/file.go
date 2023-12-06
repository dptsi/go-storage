package contracts

import "io"

type File interface {
	io.Reader
	io.Seeker
}
