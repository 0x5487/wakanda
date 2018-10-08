package middleware

import (
	"fmt"

	"gitlab.paradise-soft.com.tw/bbs/bbs/types"
	"gitlab.paradise-soft.com.tw/rd/log"
	"gitlab.paradise-soft.com.tw/rd/napnap"
)

type ErrorHandingMiddleware struct {
}

func NewErrorHandingMiddleware() *ErrorHandingMiddleware {
	return &ErrorHandingMiddleware{}
}

func (m *ErrorHandingMiddleware) Invoke(c *napnap.Context, next napnap.HandlerFunc) {
	defer func() {
		// we only handle error for bifrost application and don't handle can't error from upstream.
		if r := recover(); r != nil {
			// bad request.  hanlder status code is 400 series.
			appError, ok := r.(types.AppError)
			if ok {
				if appError.ErrorCode == "not_found" {
					c.JSON(404, appError)
					return
				}
				c.JSON(400, appError)
				return
			}

			// unknown error.  hanlder status code is 500 series.
			logger := log.StackTrace()
			customFields := log.Fields{
				"url": c.Request.RequestURI,
			}
			err, ok := r.(error)
			if !ok {
				if err == nil {
					err = fmt.Errorf("%v", r)
				} else {
					err = fmt.Errorf("%v", err)
				}
			}
			logger = logger.WithFields(customFields)
			logger.Errorf("unknown error: %v", err)

			appError = types.AppError{
				ErrorCode: "unknown_error",
				Message:   err.Error(),
			}
			c.JSON(500, appError)
		}
	}()
	next(c)
}
