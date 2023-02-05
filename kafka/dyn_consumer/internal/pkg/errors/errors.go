package errors

/*
	Taken from https://git.sipp-now.com/spid/direct-debit/core-v2/-/blob/master/pkg/core/service/error.go
*/

type Error interface {
	error
	Code() int
}
