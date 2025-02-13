package app

import (
	"dating-apps/helper/config"
	"dating-apps/helper/logger"
	"fmt"

	db "dating-apps/helper/database"
)

type Infra struct {
	Db     *db.Database
	Log    logger.Logger
	Config *config.Config
}

func (app *Infra) UrlWithPrefix(url string) string {
	return fmt.Sprintf("%s%s", app.Config.BasePrefix(), url)
}
