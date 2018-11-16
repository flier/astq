package query

func (items TypeDeclMap) WithTagValue(key, value string) TypeDeclMap {
	return items.Filter(func(name string, ty *TypeDecl) bool {
		tags := ty.Tags()
		v, found := tags[key]

		return found && v == value
	})
}

func (items TypeDeclMap) WithTag(key string) TypeDeclMap {
	return items.Filter(func(name string, ty *TypeDecl) bool {
		tags := ty.Tags()
		_, found := tags[key]

		return found
	})
}

func (items TypeDeclMap) WithoutTag(key string) TypeDeclMap {
	return items.Filter(func(name string, ty *TypeDecl) bool {
		tags := ty.Tags()
		_, found := tags[key]

		return !found
	})
}

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
