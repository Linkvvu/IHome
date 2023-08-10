package dao

type notExistErr struct {
	detail string
}

func (err notExistErr) Error() string {
	return err.detail
}

var NotExistError notExistErr = notExistErr{"target not exist"}
