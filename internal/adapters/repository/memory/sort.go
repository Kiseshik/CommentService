package memory

import (
	"sort"
	"time"
)

type sortable interface {
	GetCreatedAt() time.Time
	GetID() string
}

func sortByCreatedAt[T sortable](items []T) {
	sort.Slice(items, func(i, j int) bool {
		if items[i].GetCreatedAt().Equal(items[j].GetCreatedAt()) {
			return items[i].GetID() < items[j].GetID()
		}
		return items[i].GetCreatedAt().Before(items[j].GetCreatedAt())
	})
}
