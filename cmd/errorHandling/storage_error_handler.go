package errorHandling

import (
	"fmt"
)

/*StorageErrorHandler handles errors of any interface and simply stores what was passed to it
in an error. There is a maximum number of errors it can store before it simply panics and
prints out all the errors to stdout*/
type StorageErrorHandler struct {
	errors     [256]errorObject
	maxError   uint8
	errorCount uint8
}

type errorObject struct {
	code    uint
	message string
}

func (e errorObject) String() string {
	return fmt.Sprintf("code %d, %s", e.code, e.message)
}

/*MakeStorageErrorHandler is a constructor for StorageErrorHandler*/
func MakeStorageErrorHandler(maxErrorNumber uint8) StorageErrorHandler {
	handler := StorageErrorHandler{
		maxError:   maxErrorNumber,
		errorCount: 0,
	}

	return handler
}

/*Handle is designed to handle the storage of error objects in the form of
a Uint, representing an error code, and a String, as a descriptor of the source, if available, as well as an accompanying string to describe
the source of the error object, if available.
*/
func (h *StorageErrorHandler) Handle(code uint, message string) {
	if h.GetErrorCount() < h.GetErrorCapacity() {
		h.addError(code, message)
	} else {
		// we cannot handle another error. Time to dumpa all the errors.
		for i := uint(0); i < h.GetErrorCount(); i++ {
			fmt.Printf("%d: %s", i, h.errors[i])
		}
		fmt.Printf("%d, %s", h.GetErrorCount(), errorObject{code: code, message: message})
	}
}

func (h *StorageErrorHandler) addError(code uint, message string) {
	h.errors[h.GetErrorCount()] = errorObject{
		code:    code,
		message: message,
	}
	h.incrementErrorCount()
}

func (h *StorageErrorHandler) incrementErrorCount() {
	h.errorCount++
}

/*GetErrorCount returns the number of errors currently stored in the handler
 */
func (h *StorageErrorHandler) GetErrorCount() uint {
	return uint(h.errorCount)
}

/*GetErrorCapacity returns the maximum number of errors the handler can hold*/
func (h *StorageErrorHandler) GetErrorCapacity() uint {
	return uint(h.maxError) + 1
}
