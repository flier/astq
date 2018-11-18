package query

import (
	"go/ast"
	"reflect"
	"regexp"
	"strings"
)

const (
	tagPrefix = "+tag"
)

var (
	reTag = regexp.MustCompile(`(\w+):`)
)

type Tags []reflect.StructTag

func (tags Tags) Contains(key string) bool {
	_, ok := tags.Lookup(key)

	return ok
}

func (tags Tags) Get(key string) string {
	value, _ := tags.Lookup(key)

	return value
}

func (tags Tags) Lookup(key string) (value string, ok bool) {
	for _, tag := range tags {
		value, ok = tag.Lookup(key)

		if ok {
			return
		}
	}

	return
}

func extractTags(groups ...*ast.CommentGroup) (tags Tags) {
	if groups != nil {
		for _, group := range groups {
			if group == nil {
				continue
			}

			for _, comment := range group.List {
				if comment := strings.TrimSpace(strings.TrimLeft(comment.Text, "/")); strings.HasPrefix(comment, tagPrefix) {
					tags = append(tags, reflect.StructTag(strings.TrimSpace(strings.TrimPrefix(comment, tagPrefix))))
				}
			}
		}
	}

	return
}
