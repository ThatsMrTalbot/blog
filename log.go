package blog

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"math/rand"
	"net/http"
	"unsafe"

	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
)

func GetLog(ctx context.Context) logrus.FieldLogger {
	if l, ok := ctx.Value("logger").(logrus.FieldLogger); ok {
		return l
	}
	return logrus.StandardLogger()
}

func StoreLog(ctx context.Context, log logrus.FieldLogger) context.Context {
	return context.WithValue(ctx, "logger", log)
}

func RequestID(r *http.Request) string {
	ptr := uint64(uintptr(unsafe.Pointer(r)))
	w := md5.New()
	binary.Write(w, binary.BigEndian, ptr)
	binary.Write(w, binary.BigEndian, rand.Int63())
	return hex.EncodeToString(w.Sum(nil))
}
