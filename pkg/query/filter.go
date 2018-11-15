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

func (items FieldMap) WithTagValue(key, value string) FieldMap {
	return items.Filter(func(name string, f *Field) bool {
		v, found := f.Tag().Lookup(key)

		return found && v == value
	})
}

func (items FieldMap) WithTag(key string) FieldMap {
	return items.Filter(func(name string, f *Field) bool {
		_, found := f.Tag().Lookup(key)

		return found
	})
}

func (items FieldMap) WithoutTag(key string) FieldMap {
	return items.Filter(func(name string, f *Field) bool {
		_, found := f.Tag().Lookup(key)

		return !found
	})
}
