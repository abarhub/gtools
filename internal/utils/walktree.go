package utils

import (
	"os"
	"path"
)

type FilterFile = func(entry string) (bool, error)
type FileParse = func(entry string, entry2 string) error

type Walktree struct {
	rootDirectory  string
	filter         FilterFile
	directoryParse FileParse
	fileParse      FileParse
	exclude        []string
	include        []string
	directory2     string
	recursive      bool
}

func (receiver *Walktree) SetFilter(filter FilterFile) {
	receiver.filter = filter
}

func (receiver *Walktree) SetDirectoryParse(directoryParse FileParse) {
	receiver.directoryParse = directoryParse
}

func (receiver *Walktree) SetFileParse(fileParse FileParse) {
	receiver.fileParse = fileParse
}

func (receiver *Walktree) SetDir2(dir2 string) {
	receiver.directory2 = dir2
}

func (receiver *Walktree) SetRecursive(recursive bool) {
	receiver.recursive = recursive
}

func CreateWalktree(directory string, exclude, include []string) (*Walktree, error) {
	result := Walktree{rootDirectory: directory, exclude: exclude, include: include, recursive: true}
	return &result, nil
}

func (receiver *Walktree) Parse() error {
	return receiver.parse2(receiver.rootDirectory, receiver.directory2)
}

func (receiver *Walktree) parse2(dir string, dir2 string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, f := range files {
		srcFile := path.Join(dir, f.Name())
		destFile := path.Join(dir2, f.Name())

		skip := false
		if receiver.filter != nil {
			skip, err = receiver.filter(srcFile)
			if err != nil {
				return err
			}
		}

		if !skip {
			toCopy, err := receiver.fileToCopy(srcFile, true)
			if err != nil {
				return err
			} else if toCopy {
				if f.IsDir() {
					if receiver.recursive {
						if receiver.directoryParse != nil {
							err = receiver.directoryParse(srcFile, destFile)
							if err != nil {
								return err
							}
						}
						err = receiver.parse2(srcFile, destFile)
						if err != nil {
							return err
						}
					}
				} else {
					toCopy, err = receiver.fileToCopy(srcFile, false)
					if err != nil {
						return err
					} else if toCopy {
						if receiver.fileParse != nil {
							err = receiver.fileParse(srcFile, destFile)
							if err != nil {
								return err
							}
						}
					}
				}
			}
		}
	}

	return nil
}

func (receiver *Walktree) fileToCopy(file string, exclude bool) (bool, error) {
	if exclude && len(receiver.exclude) > 0 {
		for _, s := range receiver.exclude {
			match, err := matchGlob(file, s)
			if err != nil {
				return false, err
			} else if match {
				return false, nil
			}
		}
	}
	if !exclude && len(receiver.include) > 0 {
		for _, s := range receiver.include {
			match, err := matchGlob(file, s)
			if err != nil {
				return false, err
			} else if match {
				return true, nil
			}
		}
		return false, nil
	}
	return true, nil
}

func matchGlob(file, pattern string) (bool, error) {
	return MatchGlob(file, pattern)
}
