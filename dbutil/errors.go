package dbutil

import (
	"github.com/jackc/pgx"
)

func IsErrorAboutDublicate(err error) bool {
	if err, ok := err.(pgx.PgError); ok {
		if err.Code == "23505" {
			return true
		}
	}
	return false
}

func IsErrorAboutDublicateReturnConstaint(err error) (bool, string) {
	if err, ok := err.(pgx.PgError); ok {
		if err.Code == "23505" {
			return true, err.ConstraintName
		}
	}
	return false, ""
}

func IsErrorAboutFailedForeignKey(err error) bool {
	if err, ok := err.(pgx.PgError); ok {
		if err.Code == "23503" || err.Code == "23502" {
			return true
		}
	}
	return false
}

func IsErrorAboutFailedForeignKeyReturnConstraint(err error) (bool, string) {
	if err, ok := err.(pgx.PgError); ok {
		if err.Code == "23503" {
			return true, err.ConstraintName
		}
	}
	return false, ""
}

func IsErrorAboutNotFound(err error) bool {
	return err.Error() == "no rows in result set"
}
