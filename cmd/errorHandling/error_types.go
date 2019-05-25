package errorHandling

/*ErrorCode is a type alias for the valid error codes in this
application*/
type ErrorCode uint

const (
	/*Memory represents that something has failed while attempting to read or write from memory*/
	Memory ErrorCode = iota
)
