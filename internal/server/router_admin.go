package server

import (
	"net/http"
	"os"
	"path/filepath"

	"connectrpc.com/connect/v2"
	"connectrpc.com/connect/v2/connecthttp"
	"github.com/samber/do/v2"

	admv1 "github.com/modelgate/modelgate/internal/app/admin/v1"
	"github.com/modelgate/modelgate/internal/server/interceptor"
	"github.com/modelgate/modelgate/internal/server/middleware"
	v1pb "github.com/modelgate/modelgate/pkg/proto/admin/v1"
)

func registerHandlers(container do.Injector) http.Handler {
	apiMux := http.NewServeMux()
	rpcServer := connect.NewServer(interceptor.Logger())

	v1pb.RegisterAuthServiceHandler(rpcServer, do.MustInvoke[*admv1.AuthService](container))
	v1pb.RegisterSystemServiceHandler(rpcServer, do.MustInvoke[*admv1.SystemService](container))
	v1pb.RegisterRelayServiceHandler(rpcServer, do.MustInvoke[*admv1.RelayService](container))
	connecthttp.Mount(apiMux, rpcServer, connecthttp.WithReadMaxBytes(0))

	apiHandler := middleware.Auth(container)(apiMux)

	rootMux := http.NewServeMux()
	// webHandler
	webHandler(rootMux, http.Dir("./web/dist"))
	// apiHandler
	rootMux.Handle("/api/", http.StripPrefix("/api", apiHandler))

	handler := wrapHandler(rootMux, middleware.CORS())
	return handler
}

// 先执行的中间件在最里层
func wrapHandler(handler http.Handler, chains ...func(http.Handler) http.Handler) http.Handler {
	for i := len(chains) - 1; i >= 0; i-- {
		handler = chains[i](handler)
	}
	return handler
}

func webHandler(mux *http.ServeMux, dir http.Dir) {
	fs := http.FileServer(dir)

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(string(dir), filepath.Clean(r.URL.Path))
		if _, err := os.Stat(path); os.IsNotExist(err) {
			r.URL.Path = "/"
		}
		fs.ServeHTTP(w, r)
	}))
}
