package errorHandling

import (
	"fmt"
)

/*The interesting thing about error handling with panics is that every application must handle errors independently of the
libraries that are used by the application. The library has no way of knowing how the application intends to interpret any
errors that occur, so in that sense, some information is lost when errors are propagated from generic library code to application code,
which is very interesting.
*/

/*MemoryErrorHandler handles memory errors.
 */
type MemoryErrorHandler struct {
	storage StorageErrorHandler
}

/*MakeMemoryErrorHandler constructs a MemoryErrorHandler*/
func MakeMemoryErrorHandler(maxErrorNumber uint8) MemoryErrorHandler {
	handler := MemoryErrorHandler{
		storage: MakeStorageErrorHandler(maxErrorNumber),
	}
	return handler
}

/*Handle attempts to handle a message sent by some error-producing code. In reality,
it expects a string describing the error, and only a string*/
func (h *MemoryErrorHandler) Handle(message interface{}) {
	switch m := message.(type) {
	case string:
		h.storage.Handle(uint(Memory), m)
	default:
		panic(fmt.Sprintf("%s", message))
	}
}

/*HasErrors returns whether this memory error stores any errors.
 */
func (h *MemoryErrorHandler) HasErrors() bool {
	return h.storage.GetErrorCount() > 0
}

/*GetErrorCount returns the count of memory errors
 */
func (h *MemoryErrorHandler) GetErrorCount() uint {
	return h.storage.GetErrorCount()
}

/*ClearErrors clears the errors -- the error count will go to 0, and HasErrors will go back to false
 */
func (h *MemoryErrorHandler) ClearErrors() {
	h.storage = MakeStorageErrorHandler(uint8(h.storage.GetMaxErrorNumber()))
}
