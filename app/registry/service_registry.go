package registry

import (
	"dating-apps/app"
	"dating-apps/app/repository"
	"dating-apps/app/service"
)

func RegisterDatingAppService(app *app.Infra) service.DatingAppService {
	return service.NewDatingAppService(
		app,
		repository.NewDatingAppRepository(app.Db),
	)
}
