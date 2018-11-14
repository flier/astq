package query

import "strings"

type TypeFilter func(ty *Type) bool

func (items TypeMap) Filter(filter TypeFilter) TypeMap {
	filtered := make(TypeMap)

	for name, intf := range items {
		if filter(intf) {
			filtered[name] = intf
		}
	}

	return filtered
}

func (items TypeMap) WithPrefix(prefix string) TypeMap {
	return items.Filter(func(ty *Type) bool {
		return strings.HasPrefix(ty.Name(), prefix)
	})
}

func (items TypeMap) WithSuffix(prefix string) TypeMap {
	return items.Filter(func(ty *Type) bool {
		return strings.HasSuffix(ty.Name(), prefix)
	})
}

func (items TypeMap) WithTagValue(name, value string) TypeMap {
	return items.Filter(func(ty *Type) bool {
		tags := ty.Tags()
		v, found := tags[name]

		return found && v == value
	})
}

func (items TypeMap) WithTag(name string) TypeMap {
	return items.Filter(func(ty *Type) bool {
		tags := ty.Tags()
		_, found := tags[name]

		return found
	})
}

func (items TypeMap) WithoutTag(name string) TypeMap {
	return items.Filter(func(ty *Type) bool {
		tags := ty.Tags()
		_, found := tags[name]

		return !found
	})
}

type InterfaceFilter func(intf *Interface) bool

func (items InterfaceMap) Filter(filter InterfaceFilter) InterfaceMap {
	filtered := make(InterfaceMap)

	for name, intf := range items {
		if filter(intf) {
			filtered[name] = intf
		}
	}

	return filtered
}

func (items InterfaceMap) WithPrefix(prefix string) InterfaceMap {
	return items.Filter(func(intf *Interface) bool {
		return strings.HasPrefix(intf.Name(), prefix)
	})
}

func (items InterfaceMap) WithSuffix(prefix string) InterfaceMap {
	return items.Filter(func(intf *Interface) bool {
		return strings.HasSuffix(intf.Name(), prefix)
	})
}

type StructFilter func(s *Struct) bool

func (items StructMap) Filter(filter StructFilter) StructMap {
	filtered := make(StructMap)

	for name, intf := range items {
		if filter(intf) {
			filtered[name] = intf
		}
	}

	return filtered
}

func (items StructMap) WithPrefix(prefix string) StructMap {
	return items.Filter(func(s *Struct) bool {
		return strings.HasPrefix(s.Name(), prefix)
	})
}

func (items StructMap) WithSuffix(suffix string) StructMap {
	return items.Filter(func(s *Struct) bool {
		return strings.HasSuffix(s.Name(), suffix)
	})
}

type MethodFilter func(m *Method) bool

func (items MethodMap) Filter(filter MethodFilter) MethodMap {
	filtered := make(MethodMap)

	for name, intf := range items {
		if filter(intf) {
			filtered[name] = intf
		}
	}

	return filtered
}

func (items MethodMap) WithPrefix(prefix string) MethodMap {
	return items.Filter(func(m *Method) bool {
		return strings.HasPrefix(m.Name(), prefix)
	})
}

func (items MethodMap) WithSuffix(suffix string) MethodMap {
	return items.Filter(func(m *Method) bool {
		return strings.HasSuffix(m.Name(), suffix)
	})
}

type FieldFilter func(f *Field) bool

func (items FieldMap) Filter(filter FieldFilter) FieldMap {
	filtered := make(FieldMap)

	for name, intf := range items {
		if filter(intf) {
			filtered[name] = intf
		}
	}

	return filtered
}

func (items FieldMap) WithPrefix(prefix string) FieldMap {
	return items.Filter(func(f *Field) bool {
		return strings.HasPrefix(f.Name(), prefix)
	})
}

func (items FieldMap) WithSuffix(suffix string) FieldMap {
	return items.Filter(func(f *Field) bool {
		return strings.HasSuffix(f.Name(), suffix)
	})
}

func (items FieldMap) WithTagValue(name, value string) FieldMap {
	return items.Filter(func(f *Field) bool {
		v, found := f.Tag().Lookup(name)

		return found && v == value
	})
}

func (items FieldMap) WithTag(name string) FieldMap {
	return items.Filter(func(f *Field) bool {
		_, found := f.Tag().Lookup(name)

		return found
	})
}

func (items FieldMap) WithoutTag(name string) FieldMap {
	return items.Filter(func(f *Field) bool {
		_, found := f.Tag().Lookup(name)

		return !found
	})
}
