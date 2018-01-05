package dbutil

import "github.com/lib/pq"

func IsErrorAboutDublicate(err error) bool {
	if err, ok := err.(*pq.Error); ok {
		if err.Code.Class() == "23505" {
			return true
		}
	}
	return false
}
