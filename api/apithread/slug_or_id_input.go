package apithread

import "strconv"

type slugOrIdInput struct {
	Id    int64  `json:"-"`
	Slug  string `json:"-"`
	HasId bool   `json:"-"`
}

func resolveSlugOrIdInput(slug string, t *slugOrIdInput) {
	t.Slug = slug

	if id, err := strconv.ParseInt(slug, 10, 64); err == nil {
		t.Id = id
		t.HasId = true
	} else {
		t.HasId = false
	}

	return
}
