package interceptor

import (
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect/v2"
	log "github.com/sirupsen/logrus"
)

func Logger() connect.ServerInterceptor {
	return func(next connect.ServerFunc) connect.ServerFunc {
		return func(ctx context.Context, spec connect.Spec, stream connect.ServerStream) error {
			startTime := time.Now()
			logLevel := log.InfoLevel
			logMsg := "OK"

			err := next(ctx, spec, stream)
			if err != nil {
				switch connect.CodeOf(err) {
				case connect.CodeUnauthenticated,
					connect.CodeInvalidArgument,
					connect.CodeNotFound,
					connect.CodeOutOfRange,
					connect.CodePermissionDenied:

					logLevel = log.InfoLevel
					logMsg = "client error"
				case connect.CodeInternal,
					connect.CodeUnavailable,
					connect.CodeDataLoss,
					connect.CodeDeadlineExceeded:
					logLevel = log.ErrorLevel
					logMsg = "server error"

				case connect.CodeUnknown:
					logLevel = log.ErrorLevel
					logMsg = "unknown error"
				}
			}
			logger := log.WithContext(ctx).
				WithField("method", spec.Procedure).
				WithField("elapsed_time", fmt.Sprintf("%.3fms", float64(time.Since(startTime).Microseconds())/1000))
			if err != nil {
				logger = logger.WithField("err", err.Error())
			}
			logger.Log(logLevel, logMsg)

			return err
		}
	}
}
