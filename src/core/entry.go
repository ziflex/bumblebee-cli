package core

type Entry struct {
	Id   int
	Name string
}

func NewEntry(name string) *Entry {
	return &Entry{Name: name}
}
