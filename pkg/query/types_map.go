package query

// Code generated by astgen v1.0 with go1.11.2 DO NOT EDIT

import (
	"strings"
)

// Keys returns a new slice containing the set of map keys. The order is unspecified.
func (m MethodMap) Keys() (keys []string) {
	for name := range m {
		keys = append(keys, name)
	}

	return
}

// Values returns a new slice containing the set of map values. The order is unspecified.
func (m MethodMap) Values() (values []*Method) {
	for _, value := range m {
		values = append(values, value)
	}

	return
}

// Contains reports whether key is within map.
func (m MethodMap) Contains(key string) bool {
	_, found := m[key]

	return found
}

// Clone returns a shadow copy of map.
func (m MethodMap) Clone() MethodMap {
	cloned := make(MethodMap)

	for key, value := range m {
		cloned[key] = value
	}

	return cloned
}

// Filter filters the map to only include elements for which filter returns true.
func (m MethodMap) Filter(filter func(key string, value *Method) bool) MethodMap {
	filtered := make(MethodMap)

	for key, value := range m {
		if filter(key, value) {
			filtered[key] = value
		}
	}

	return filtered
}

// WithPrefix filters the map to only include elements for which contains prefix.
func (m MethodMap) WithPrefix(prefix string) MethodMap {
	return m.Filter(func(key string, value *Method) bool {
		return strings.HasPrefix(key, prefix)
	})
}

// WithSuffix filters the map to only include elements for which contains suffix.
func (m MethodMap) WithSuffix(suffix string) MethodMap {
	return m.Filter(func(key string, value *Method) bool {
		return strings.HasSuffix(key, suffix)
	})
}

// Keys returns a new slice containing the set of map keys. The order is unspecified.
func (m NamedFieldMap) Keys() (keys []string) {
	for name := range m {
		keys = append(keys, name)
	}

	return
}

// Values returns a new slice containing the set of map values. The order is unspecified.
func (m NamedFieldMap) Values() (values []*NamedField) {
	for _, value := range m {
		values = append(values, value)
	}

	return
}

// Contains reports whether key is within map.
func (m NamedFieldMap) Contains(key string) bool {
	_, found := m[key]

	return found
}

// Clone returns a shadow copy of map.
func (m NamedFieldMap) Clone() NamedFieldMap {
	cloned := make(NamedFieldMap)

	for key, value := range m {
		cloned[key] = value
	}

	return cloned
}

// Filter filters the map to only include elements for which filter returns true.
func (m NamedFieldMap) Filter(filter func(key string, value *NamedField) bool) NamedFieldMap {
	filtered := make(NamedFieldMap)

	for key, value := range m {
		if filter(key, value) {
			filtered[key] = value
		}
	}

	return filtered
}

// WithPrefix filters the map to only include elements for which contains prefix.
func (m NamedFieldMap) WithPrefix(prefix string) NamedFieldMap {
	return m.Filter(func(key string, value *NamedField) bool {
		return strings.HasPrefix(key, prefix)
	})
}

// WithSuffix filters the map to only include elements for which contains suffix.
func (m NamedFieldMap) WithSuffix(suffix string) NamedFieldMap {
	return m.Filter(func(key string, value *NamedField) bool {
		return strings.HasSuffix(key, suffix)
	})
}

// Keys returns a new slice containing the set of map keys. The order is unspecified.
func (m Tags) Keys() (keys []string) {
	for name := range m {
		keys = append(keys, name)
	}

	return
}

// Values returns a new slice containing the set of map values. The order is unspecified.
func (m Tags) Values() (values []string) {
	for _, value := range m {
		values = append(values, value)
	}

	return
}

// Contains reports whether key is within map.
func (m Tags) Contains(key string) bool {
	_, found := m[key]

	return found
}

// Clone returns a shadow copy of map.
func (m Tags) Clone() Tags {
	cloned := make(Tags)

	for key, value := range m {
		cloned[key] = value
	}

	return cloned
}

// Filter filters the map to only include elements for which filter returns true.
func (m Tags) Filter(filter func(key string, value string) bool) Tags {
	filtered := make(Tags)

	for key, value := range m {
		if filter(key, value) {
			filtered[key] = value
		}
	}

	return filtered
}

// WithPrefix filters the map to only include elements for which contains prefix.
func (m Tags) WithPrefix(prefix string) Tags {
	return m.Filter(func(key string, value string) bool {
		return strings.HasPrefix(key, prefix)
	})
}

// WithSuffix filters the map to only include elements for which contains suffix.
func (m Tags) WithSuffix(suffix string) Tags {
	return m.Filter(func(key string, value string) bool {
		return strings.HasSuffix(key, suffix)
	})
}

// Keys returns a new slice containing the set of map keys. The order is unspecified.
func (m ValueSpecMap) Keys() (keys []string) {
	for name := range m {
		keys = append(keys, name)
	}

	return
}

// Values returns a new slice containing the set of map values. The order is unspecified.
func (m ValueSpecMap) Values() (values []*ValueSpec) {
	for _, value := range m {
		values = append(values, value)
	}

	return
}

// Contains reports whether key is within map.
func (m ValueSpecMap) Contains(key string) bool {
	_, found := m[key]

	return found
}

// Clone returns a shadow copy of map.
func (m ValueSpecMap) Clone() ValueSpecMap {
	cloned := make(ValueSpecMap)

	for key, value := range m {
		cloned[key] = value
	}

	return cloned
}

// Filter filters the map to only include elements for which filter returns true.
func (m ValueSpecMap) Filter(filter func(key string, value *ValueSpec) bool) ValueSpecMap {
	filtered := make(ValueSpecMap)

	for key, value := range m {
		if filter(key, value) {
			filtered[key] = value
		}
	}

	return filtered
}

// WithPrefix filters the map to only include elements for which contains prefix.
func (m ValueSpecMap) WithPrefix(prefix string) ValueSpecMap {
	return m.Filter(func(key string, value *ValueSpec) bool {
		return strings.HasPrefix(key, prefix)
	})
}

// WithSuffix filters the map to only include elements for which contains suffix.
func (m ValueSpecMap) WithSuffix(suffix string) ValueSpecMap {
	return m.Filter(func(key string, value *ValueSpec) bool {
		return strings.HasSuffix(key, suffix)
	})
}
