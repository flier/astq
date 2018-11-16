package query

// Code generated by astgen v1.0 with go1.11.2 DO NOT EDIT

// Filter filters the iterator to only include elements for which filter returns true.
func (c MethodIter) Filter(filter func(item *Method) bool) MethodIter {
	filtered := make(chan *Method)

	go func() {
		defer close(filtered)

		for item := range c {
			if filter(item) {
				filtered <- item
			}
		}
	}()

	return filtered
}

// Find returns the item for which filter returns true.
func (c MethodIter) Find(filter func(item *Method) bool) *Method {
	for item := range c {
		if filter(item) {
			return item
		}
	}

	return nil
}

// Collect returns a new slice including all items from the iterator.
func (c MethodIter) Collect() (items []*Method) {
	for item := range c {
		items = append(items, item)
	}

	return
}
