// Package router
package router

import (
	"context"
	"encoding/json"
	"net/http"
	"runtime/debug"

	"gitlab.privy.id/go_graphql/internal/appctx"
	"gitlab.privy.id/go_graphql/internal/bootstrap"
	"gitlab.privy.id/go_graphql/internal/consts"
	"gitlab.privy.id/go_graphql/internal/handler"
	"gitlab.privy.id/go_graphql/internal/middleware"
	repositoriesRole "gitlab.privy.id/go_graphql/internal/repositories/role"
	repositoriesProduct "gitlab.privy.id/go_graphql/internal/repositories/product"
	"gitlab.privy.id/go_graphql/internal/ucase"
	"gitlab.privy.id/go_graphql/pkg/logger"
	"gitlab.privy.id/go_graphql/pkg/msgx"
	"gitlab.privy.id/go_graphql/pkg/routerkit"

	//"gitlab.privy.id/go_graphql/pkg/mariadb"
	//"gitlab.privy.id/go_graphql/internal/repositories"
	//"gitlab.privy.id/go_graphql/internal/ucase/example"

	ucaseContract "gitlab.privy.id/go_graphql/internal/ucase/contract"
	"gitlab.privy.id/go_graphql/internal/ucase/role"
	"gitlab.privy.id/go_graphql/internal/ucase/product"
)

type router struct {
	config *appctx.Config
	router *routerkit.Router
}

// NewRouter initialize new router wil return Router Interface
func NewRouter(cfg *appctx.Config) Router {
	bootstrap.RegistryMessage()
	bootstrap.RegistryLogger(cfg)

	return &router{
		config: cfg,
		router: routerkit.NewRouter(routerkit.WithServiceName(cfg.App.AppName)),
	}
}

func (rtr *router) handle(hfn httpHandlerFunc, svc ucaseContract.UseCase, mdws ...middleware.MiddlewareFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lang := r.Header.Get(consts.HeaderLanguageKey)
		if !msgx.HaveLang(consts.RespOK, lang) {
			lang = rtr.config.App.DefaultLang
			r.Header.Set(consts.HeaderLanguageKey, lang)
		}

		defer func() {
			err := recover()
			if err != nil {
				w.Header().Set(consts.HeaderContentTypeKey, consts.HeaderContentTypeJSON)
				w.WriteHeader(http.StatusInternalServerError)
				res := appctx.Response{
					Code: consts.CodeInternalServerError,
				}

				res.WithLang(lang)
				logger.Error(logger.MessageFormat("error %v", string(debug.Stack())))
				json.NewEncoder(w).Encode(res.Byte())

				return
			}
		}()

		ctx := context.WithValue(r.Context(), "access", map[string]interface{}{
			"path":      r.URL.Path,
			"remote_ip": r.RemoteAddr,
			"method":    r.Method,
		})

		req := r.WithContext(ctx)

		// validate middleware
		if !middleware.FilterFunc(w, req, rtr.config, mdws) {
			return
		}

		resp := hfn(req, svc, rtr.config)
		resp.WithLang(lang)
		rtr.response(w, resp)
	}
}

// response prints as a json and formatted string for DGP legacy
func (rtr *router) response(w http.ResponseWriter, resp appctx.Response) {
	w.Header().Set(consts.HeaderContentTypeKey, consts.HeaderContentTypeJSON)
	resp.Generate()
	w.WriteHeader(resp.Code)
	w.Write(resp.Byte())
	return
}

// postgresql
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "intern-privy-simple-orders"
)

// Route preparing http router and will return mux router object
func (rtr *router) Route() *routerkit.Router {

	root := rtr.router.PathPrefix("/").Subrouter()
	//in := root.PathPrefix("/in/").Subrouter()
	liveness := root.PathPrefix("/").Subrouter()
	//inV1 := in.PathPrefix("/v1/").Subrouter()

	// open tracer setup
	bootstrap.RegistryOpenTracing(rtr.config)

	//db := bootstrap.RegistryMariaMasterSlave(rtr.config.WriteDB, rtr.config.ReadDB, rtr.config.App.Timezone))
	//db := bootstrap.RegistryMariaDB(rtr.config.WriteDB, rtr.config.App.Timezone)

	// use case
	healthy := ucase.NewHealthCheck()

	// healthy
	liveness.HandleFunc("/liveness", rtr.handle(
		handler.HttpRequest,
		healthy,
	)).Methods(http.MethodGet)

	// db := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	// result :=

	// db := bootstrap.RegistryPostgreSQLMasterSlave(rtr.config.ReadDB, rtr.config.WriteDB, rtr.config.App.Timezone)

	// this is use case for example purpose, please delete
	//repoExample := repositories.NewExample(db)
	//el := example.NewExampleList(repoExample)
	//ec := example.NewPartnerCreate(repoExample
	//ed := example.NewExampleDelete(repoExample)

	// TODO: create your route here

	// this route for example rest, please delete
	// example list
	//inV1.HandleFunc("/example", rtr.handle(
	//    handler.HttpRequest,
	//    el,
	//)).Methods(http.MethodGet)

	//inV1.HandleFunc("/example", rtr.handle(
	//    handler.HttpRequest,
	//    ec,
	//)).Methods(http.MethodPost)

	//inV1.HandleFunc("/example/{id:[0-9]+}", rtr.handle(
	//    handler.HttpRequest,
	//    ed,
	//)).Methods(http.MethodDelete)

	db := bootstrap.RegistryPostgreSQLMasterSlave(rtr.config.WriteDB, rtr.config.ReadDB, rtr.config.App.Timezone)
	//Repository
	roleRepo := repositoriesRole.NewRoleRepo(db)
	productRepo := repositoriesProduct.NewProductRepo(db)

	//Ucase
	createRoleUseCase := role.NewCreateRole(roleRepo)
	createProductUseCase := product.NewCreateProduct(productRepo)
	getProductUsecase := product.NewGetProducts(productRepo)

	//Route

	
	root.HandleFunc("/api/v1/role/create", rtr.handle(
		handler.HttpRequest,
		createRoleUseCase,
	)).Methods(http.MethodPost)

	root.HandleFunc("/api/v1/product/create", rtr.handle(
		handler.HttpRequest,
		createProductUseCase,
	)).Methods(http.MethodPost)

	root.HandleFunc("/api/v1/products", rtr.handle(
		handler.HttpRequest,
		getProductUsecase,
	)).Methods(http.MethodGet)












	return rtr.router

}
