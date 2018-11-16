package query

func (items NamedFieldMap) WithTagValue(key, value string) NamedFieldMap {
	return items.Filter(func(name string, f *NamedField) bool {
		v, found := f.Tag().Lookup(key)

		return found && v == value
	})
}

func (items NamedFieldMap) WithTag(key string) NamedFieldMap {
	return items.Filter(func(name string, f *NamedField) bool {
		_, found := f.Tag().Lookup(key)

		return found
	})
}

func (items NamedFieldMap) WithoutTag(key string) NamedFieldMap {
	return items.Filter(func(name string, f *NamedField) bool {
		_, found := f.Tag().Lookup(key)

		return !found
	})
}
