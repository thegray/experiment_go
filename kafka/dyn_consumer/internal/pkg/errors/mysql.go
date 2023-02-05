package errors

const (
	systemDatabaseError          = "system:mysql:database"
	systemDatabaseNoRows         = "system:mysql:no_rows"
	systemDatabaseDuplicateEntry = "system:mysql:duplicate_entry"
)

var (
	ErrDatabase = func(err error) ServiceError {
		return serviceError{
			code:    systemDatabaseError,
			message: "general database",
			inner:   err,
		}
	}

	ErrDatabaseNoRows = func(err error) ServiceError {
		return serviceError{
			code:    systemDatabaseNoRows,
			message: "no rows found",
			inner:   err,
		}
	}

	ErrDatabaseDuplicateEntry = func(err error) ServiceError {
		return serviceError{
			code:    systemDatabaseDuplicateEntry,
			message: "duplicate entry",
			inner:   err,
		}
	}
)
