package blog

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"net/http"
	"unsafe"

	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
)

// GetLog gets the log
func GetLog(ctx context.Context) logrus.FieldLogger {
	if l, ok := ctx.Value("logger").(logrus.FieldLogger); ok {
		return l
	}
	return logrus.StandardLogger()
}

// StoreLog stores a log in the context
func StoreLog(ctx context.Context, log logrus.FieldLogger) context.Context {
	return context.WithValue(ctx, "logger", log)
}

// RequestID generates a request id from a request
func RequestID(r *http.Request) string {
	ptr := uint64(uintptr(unsafe.Pointer(r)))
	w := md5.New()
	binary.Write(w, binary.BigEndian, ptr)
	binary.Write(w, binary.BigEndian, r.URL.String())
	return hex.EncodeToString(w.Sum(nil))
}
