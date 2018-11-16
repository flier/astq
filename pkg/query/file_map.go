package query

// Code generated by astgen v1.0 with go1.11.2 DO NOT EDIT

import (
	"strings"
)

// Keys returns a new slice containing the set of map keys. The order is unspecified.
func (m ConstDeclMap) Keys() (keys []string) {
	for name := range m {
		keys = append(keys, name)
	}

	return
}

// Values returns a new slice containing the set of map values. The order is unspecified.
func (m ConstDeclMap) Values() (values []*ConstDecl) {
	for _, value := range m {
		values = append(values, value)
	}

	return
}

// Contains reports whether key is within map.
func (m ConstDeclMap) Contains(key string) bool {
	_, found := m[key]

	return found
}

// Clone returns a shadow copy of map.
func (m ConstDeclMap) Clone() ConstDeclMap {
	cloned := make(ConstDeclMap)

	for key, value := range m {
		cloned[key] = value
	}

	return cloned
}

// Filter filters the map to only include elements for which filter returns true.
func (m ConstDeclMap) Filter(filter func(key string, value *ConstDecl) bool) ConstDeclMap {
	filtered := make(ConstDeclMap)

	for key, value := range m {
		if filter(key, value) {
			filtered[key] = value
		}
	}

	return filtered
}

// WithPrefix filters the map to only include elements for which contains prefix.
func (m ConstDeclMap) WithPrefix(prefix string) ConstDeclMap {
	return m.Filter(func(key string, value *ConstDecl) bool {
		return strings.HasPrefix(key, prefix)
	})
}

// WithSuffix filters the map to only include elements for which contains suffix.
func (m ConstDeclMap) WithSuffix(suffix string) ConstDeclMap {
	return m.Filter(func(key string, value *ConstDecl) bool {
		return strings.HasSuffix(key, suffix)
	})
}

// Keys returns a new slice containing the set of map keys. The order is unspecified.
func (m FileMap) Keys() (keys []string) {
	for name := range m {
		keys = append(keys, name)
	}

	return
}

// Values returns a new slice containing the set of map values. The order is unspecified.
func (m FileMap) Values() (values []*File) {
	for _, value := range m {
		values = append(values, value)
	}

	return
}

// Contains reports whether key is within map.
func (m FileMap) Contains(key string) bool {
	_, found := m[key]

	return found
}

// Clone returns a shadow copy of map.
func (m FileMap) Clone() FileMap {
	cloned := make(FileMap)

	for key, value := range m {
		cloned[key] = value
	}

	return cloned
}

// Filter filters the map to only include elements for which filter returns true.
func (m FileMap) Filter(filter func(key string, value *File) bool) FileMap {
	filtered := make(FileMap)

	for key, value := range m {
		if filter(key, value) {
			filtered[key] = value
		}
	}

	return filtered
}

// WithPrefix filters the map to only include elements for which contains prefix.
func (m FileMap) WithPrefix(prefix string) FileMap {
	return m.Filter(func(key string, value *File) bool {
		return strings.HasPrefix(key, prefix)
	})
}

// WithSuffix filters the map to only include elements for which contains suffix.
func (m FileMap) WithSuffix(suffix string) FileMap {
	return m.Filter(func(key string, value *File) bool {
		return strings.HasSuffix(key, suffix)
	})
}

// Keys returns a new slice containing the set of map keys. The order is unspecified.
func (m FuncDeclMap) Keys() (keys []string) {
	for name := range m {
		keys = append(keys, name)
	}

	return
}

// Values returns a new slice containing the set of map values. The order is unspecified.
func (m FuncDeclMap) Values() (values []*FuncDecl) {
	for _, value := range m {
		values = append(values, value)
	}

	return
}

// Contains reports whether key is within map.
func (m FuncDeclMap) Contains(key string) bool {
	_, found := m[key]

	return found
}

// Clone returns a shadow copy of map.
func (m FuncDeclMap) Clone() FuncDeclMap {
	cloned := make(FuncDeclMap)

	for key, value := range m {
		cloned[key] = value
	}

	return cloned
}

// Filter filters the map to only include elements for which filter returns true.
func (m FuncDeclMap) Filter(filter func(key string, value *FuncDecl) bool) FuncDeclMap {
	filtered := make(FuncDeclMap)

	for key, value := range m {
		if filter(key, value) {
			filtered[key] = value
		}
	}

	return filtered
}

// WithPrefix filters the map to only include elements for which contains prefix.
func (m FuncDeclMap) WithPrefix(prefix string) FuncDeclMap {
	return m.Filter(func(key string, value *FuncDecl) bool {
		return strings.HasPrefix(key, prefix)
	})
}

// WithSuffix filters the map to only include elements for which contains suffix.
func (m FuncDeclMap) WithSuffix(suffix string) FuncDeclMap {
	return m.Filter(func(key string, value *FuncDecl) bool {
		return strings.HasSuffix(key, suffix)
	})
}

// Keys returns a new slice containing the set of map keys. The order is unspecified.
func (m ImportDeclMap) Keys() (keys []string) {
	for name := range m {
		keys = append(keys, name)
	}

	return
}

// Values returns a new slice containing the set of map values. The order is unspecified.
func (m ImportDeclMap) Values() (values []*ImportDecl) {
	for _, value := range m {
		values = append(values, value)
	}

	return
}

// Contains reports whether key is within map.
func (m ImportDeclMap) Contains(key string) bool {
	_, found := m[key]

	return found
}

// Clone returns a shadow copy of map.
func (m ImportDeclMap) Clone() ImportDeclMap {
	cloned := make(ImportDeclMap)

	for key, value := range m {
		cloned[key] = value
	}

	return cloned
}

// Filter filters the map to only include elements for which filter returns true.
func (m ImportDeclMap) Filter(filter func(key string, value *ImportDecl) bool) ImportDeclMap {
	filtered := make(ImportDeclMap)

	for key, value := range m {
		if filter(key, value) {
			filtered[key] = value
		}
	}

	return filtered
}

// WithPrefix filters the map to only include elements for which contains prefix.
func (m ImportDeclMap) WithPrefix(prefix string) ImportDeclMap {
	return m.Filter(func(key string, value *ImportDecl) bool {
		return strings.HasPrefix(key, prefix)
	})
}

// WithSuffix filters the map to only include elements for which contains suffix.
func (m ImportDeclMap) WithSuffix(suffix string) ImportDeclMap {
	return m.Filter(func(key string, value *ImportDecl) bool {
		return strings.HasSuffix(key, suffix)
	})
}

// Keys returns a new slice containing the set of map keys. The order is unspecified.
func (m InterfaceMap) Keys() (keys []string) {
	for name := range m {
		keys = append(keys, name)
	}

	return
}

// Values returns a new slice containing the set of map values. The order is unspecified.
func (m InterfaceMap) Values() (values []*InterfaceDef) {
	for _, value := range m {
		values = append(values, value)
	}

	return
}

// Contains reports whether key is within map.
func (m InterfaceMap) Contains(key string) bool {
	_, found := m[key]

	return found
}

// Clone returns a shadow copy of map.
func (m InterfaceMap) Clone() InterfaceMap {
	cloned := make(InterfaceMap)

	for key, value := range m {
		cloned[key] = value
	}

	return cloned
}

// Filter filters the map to only include elements for which filter returns true.
func (m InterfaceMap) Filter(filter func(key string, value *InterfaceDef) bool) InterfaceMap {
	filtered := make(InterfaceMap)

	for key, value := range m {
		if filter(key, value) {
			filtered[key] = value
		}
	}

	return filtered
}

// WithPrefix filters the map to only include elements for which contains prefix.
func (m InterfaceMap) WithPrefix(prefix string) InterfaceMap {
	return m.Filter(func(key string, value *InterfaceDef) bool {
		return strings.HasPrefix(key, prefix)
	})
}

// WithSuffix filters the map to only include elements for which contains suffix.
func (m InterfaceMap) WithSuffix(suffix string) InterfaceMap {
	return m.Filter(func(key string, value *InterfaceDef) bool {
		return strings.HasSuffix(key, suffix)
	})
}

// Keys returns a new slice containing the set of map keys. The order is unspecified.
func (m StructMap) Keys() (keys []string) {
	for name := range m {
		keys = append(keys, name)
	}

	return
}

// Values returns a new slice containing the set of map values. The order is unspecified.
func (m StructMap) Values() (values []*StructDef) {
	for _, value := range m {
		values = append(values, value)
	}

	return
}

// Contains reports whether key is within map.
func (m StructMap) Contains(key string) bool {
	_, found := m[key]

	return found
}

// Clone returns a shadow copy of map.
func (m StructMap) Clone() StructMap {
	cloned := make(StructMap)

	for key, value := range m {
		cloned[key] = value
	}

	return cloned
}

// Filter filters the map to only include elements for which filter returns true.
func (m StructMap) Filter(filter func(key string, value *StructDef) bool) StructMap {
	filtered := make(StructMap)

	for key, value := range m {
		if filter(key, value) {
			filtered[key] = value
		}
	}

	return filtered
}

// WithPrefix filters the map to only include elements for which contains prefix.
func (m StructMap) WithPrefix(prefix string) StructMap {
	return m.Filter(func(key string, value *StructDef) bool {
		return strings.HasPrefix(key, prefix)
	})
}

// WithSuffix filters the map to only include elements for which contains suffix.
func (m StructMap) WithSuffix(suffix string) StructMap {
	return m.Filter(func(key string, value *StructDef) bool {
		return strings.HasSuffix(key, suffix)
	})
}

// Keys returns a new slice containing the set of map keys. The order is unspecified.
func (m TypeDeclMap) Keys() (keys []string) {
	for name := range m {
		keys = append(keys, name)
	}

	return
}

// Values returns a new slice containing the set of map values. The order is unspecified.
func (m TypeDeclMap) Values() (values []*TypeDecl) {
	for _, value := range m {
		values = append(values, value)
	}

	return
}

// Contains reports whether key is within map.
func (m TypeDeclMap) Contains(key string) bool {
	_, found := m[key]

	return found
}

// Clone returns a shadow copy of map.
func (m TypeDeclMap) Clone() TypeDeclMap {
	cloned := make(TypeDeclMap)

	for key, value := range m {
		cloned[key] = value
	}

	return cloned
}

// Filter filters the map to only include elements for which filter returns true.
func (m TypeDeclMap) Filter(filter func(key string, value *TypeDecl) bool) TypeDeclMap {
	filtered := make(TypeDeclMap)

	for key, value := range m {
		if filter(key, value) {
			filtered[key] = value
		}
	}

	return filtered
}

// WithPrefix filters the map to only include elements for which contains prefix.
func (m TypeDeclMap) WithPrefix(prefix string) TypeDeclMap {
	return m.Filter(func(key string, value *TypeDecl) bool {
		return strings.HasPrefix(key, prefix)
	})
}

// WithSuffix filters the map to only include elements for which contains suffix.
func (m TypeDeclMap) WithSuffix(suffix string) TypeDeclMap {
	return m.Filter(func(key string, value *TypeDecl) bool {
		return strings.HasSuffix(key, suffix)
	})
}

// Keys returns a new slice containing the set of map keys. The order is unspecified.
func (m VarDeclMap) Keys() (keys []string) {
	for name := range m {
		keys = append(keys, name)
	}

	return
}

// Values returns a new slice containing the set of map values. The order is unspecified.
func (m VarDeclMap) Values() (values []*VarDecl) {
	for _, value := range m {
		values = append(values, value)
	}

	return
}

// Contains reports whether key is within map.
func (m VarDeclMap) Contains(key string) bool {
	_, found := m[key]

	return found
}

// Clone returns a shadow copy of map.
func (m VarDeclMap) Clone() VarDeclMap {
	cloned := make(VarDeclMap)

	for key, value := range m {
		cloned[key] = value
	}

	return cloned
}

// Filter filters the map to only include elements for which filter returns true.
func (m VarDeclMap) Filter(filter func(key string, value *VarDecl) bool) VarDeclMap {
	filtered := make(VarDeclMap)

	for key, value := range m {
		if filter(key, value) {
			filtered[key] = value
		}
	}

	return filtered
}

// WithPrefix filters the map to only include elements for which contains prefix.
func (m VarDeclMap) WithPrefix(prefix string) VarDeclMap {
	return m.Filter(func(key string, value *VarDecl) bool {
		return strings.HasPrefix(key, prefix)
	})
}

// WithSuffix filters the map to only include elements for which contains suffix.
func (m VarDeclMap) WithSuffix(suffix string) VarDeclMap {
	return m.Filter(func(key string, value *VarDecl) bool {
		return strings.HasSuffix(key, suffix)
	})
}
