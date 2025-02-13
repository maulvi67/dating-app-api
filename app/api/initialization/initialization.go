package initialization

import (
	"dating-apps/app"
	"dating-apps/app/api/handler"
	"dating-apps/app/api/middleware"
	"dating-apps/app/model/entity"
	"dating-apps/app/registry"
	"dating-apps/helper/config"
	db "dating-apps/helper/database"
	"net/http"
)

func InitDatabase(cfg *config.DBConfig) (db.Database, error) {
	return db.NewDatabaseConnect(cfg, &db.Option{
		Migrate: []interface{}{
			&entity.Swipe{},
			&entity.User{},
		},
	})
}

func InitRouting(app *app.Infra) *http.ServeMux {

	//authenticate
	auth := middleware.Authenticate(&app.Config.SecurityConfig)

	//Service registry
	datingAppSvc := registry.RegisterDatingAppService(app)

	// Handler initialization
	swagHttp := handler.SwaggerHttpHandler(app.Config.Url) //swagger
	authHtpp := handler.AuthHttpHandler(datingAppSvc, app)
	userHttp := handler.DatingAppHttpHandler(datingAppSvc, app)

	// Routing path
	mux := http.NewServeMux()
	mux.Handle("/", swagHttp) //swagger
	mux.HandleFunc("/__health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(http.StatusText(http.StatusOK)))
	})
	mux.Handle(app.UrlWithPrefix("auth/"), authHtpp)
	mux.Handle(app.UrlWithPrefix("user/"), middleware.Adapt(userHttp, auth))

	return mux
}
