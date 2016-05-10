package blog

import (
    "fmt"
	"github.com/ThatsMrTalbot/scaffold/errors"
)

func ErrorReponse(status int, message string, err error) error {
    msg := fmt.Sprintf("%s: %s", message, err.Error())
    return errors.NewErrorStatus(500, msg)
}