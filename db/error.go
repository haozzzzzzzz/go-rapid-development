package db

import (
	"github.com/go-sql-driver/mysql"
)

const ErrorDuplicateEntryForKey = 1062
const ErrorTableAlreadyExists = 1050

func IsErrorDuplicateEntryForKey(err error) (result bool) {
	mysqlError, ok := err.(*mysql.MySQLError)
	if !ok {
		return
	} else {
		if mysqlError.Number == ErrorDuplicateEntryForKey {
			result = true
		}
	}
	return
}

func IsErrorTableAlreadyExists(err error) (result bool) {
	mysqlError, ok := err.(*mysql.MySQLError)
	if !ok {
		return
	} else {
		if mysqlError.Number == ErrorTableAlreadyExists {
			result = true
		}
	}
	return
}
