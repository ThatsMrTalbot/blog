package blog

// Index is an article index
type Index []Article

func (i Index) Len() int {
	return len(i)
}

func (i Index) Less(e, j int) bool {
	return i[e].Mod.After(i[j].Mod)
}

func (i Index) Swap(e, j int) {
	i[e], i[j] = i[j], i[e]
}

// Pages is the number of pages given a page length
func (i Index) Pages(length int) int {
	return len(i) / length
}

// Page gets the articles on a given page given a page length
func (i Index) Page(page int, length int) []Article {
	offset := page * length
	if len(i) < offset {
		return nil
	}

	p := i[offset:]
	if len(p) > length {
		p = p[:length]
	}

	return p
}
