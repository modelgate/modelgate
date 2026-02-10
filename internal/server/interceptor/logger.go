package interceptor

import (
	"context"
	"fmt"
	"time"

	connect "connectrpc.com/connect"
	log "github.com/sirupsen/logrus"
)

func Logger() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			startTime := time.Now()
			logLevel := log.InfoLevel
			logMsg := "OK"

			resp, err := next(ctx, req)
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
				WithField("method", req.Spec().Procedure).
				WithField("elapsed_time", fmt.Sprintf("%.3fms", float64(time.Since(startTime).Microseconds())/1000))
			if err != nil {
				logger = logger.WithField("err", err.Error())
			}
			logger.Log(logLevel, logMsg)

			return resp, err
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}
