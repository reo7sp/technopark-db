package dbutil

import (
	"github.com/lib/pq"
)

func IsErrorAboutDublicate(err error) bool {
	if err, ok := err.(*pq.Error); ok {
		if err.Code == "23505" {
			return true
		}
	}
	return false
}

func IsErrorAboutDublicateReturnConstaint(err error) (bool, string) {
	if err, ok := err.(*pq.Error); ok {
		if err.Code == "23505" {
			return true, err.Constraint
		}
	}
	return false, ""
}

func IsErrorAboutFailedForeignKey(err error) bool {
	if err, ok := err.(*pq.Error); ok {
		if err.Code == "23503" {
			return true
		}
	}
	return false
}

func IsErrorAboutFailedForeignKeyReturnConstaint(err error) (bool, string) {
	if err, ok := err.(*pq.Error); ok {
		if err.Code == "23503" {
			return true, err.Constraint
		}
	}
	return false, ""
}

func IsErrorAboutNotFound(err error) bool {
	return err.Error() == "sql: no rows in result set"
}
