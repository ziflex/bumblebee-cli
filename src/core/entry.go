package core

type Entry struct {
	Id   int
	Name string
}

func NewEntry(name string) *Entry {
	return &Entry{Name: name}
}

func EntryMapToList(entries map[string]*Entry) []*Entry {
	results := make([]*Entry, 0, len(entries))

	for _, entry := range entries {
		results = append(results, entry)
	}

	return results
}

func EntryListToMap(entries map[string]*Entry) []*Entry {
	results := make([]*Entry, 0, len(entries))

	for _, entry := range entries {
		results = append(results, entry)
	}

	return results
}
