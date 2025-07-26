package directory

import "os"

type Directory struct {
	Path    string
	Created bool
}

func New(dir string) (*Directory, error) {
	newDir := &Directory{
		Path:    "",
		Created: false,
	}

	if dir == "" {
		return newDir, os.ErrInvalid
	}

	_, err := os.Stat(dir)
	if err != nil {
		return newDir, err
	}

	if !os.IsNotExist(err) {
		return newDir, os.ErrExist
	}

	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return newDir, err
	}

	newDir.Path = dir
	newDir.Created = true

	return newDir, nil

}

func (d Directory) String() string {
	return string(d.Path)
}
