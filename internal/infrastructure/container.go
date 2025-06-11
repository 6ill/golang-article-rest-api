package infrastructure

import (
	"database/sql"
	"fmt"

	"github.com/6ill/go-article-rest-api/internal/helper"
	"github.com/6ill/go-article-rest-api/internal/pkg/repository"
	"github.com/6ill/go-article-rest-api/internal/pkg/service"
	"github.com/spf13/viper"
)

var v *viper.Viper

type (
	App struct {
		ServerHost string
		ServerPort int
	}

	Container struct {
		Db             *sql.DB
		App            *App
		ArticleService service.ArticleService
	}
)

func AppInit(v *viper.Viper) (app App) {
	err := v.Unmarshal(&app)
	if err != nil {
		helper.Logger(helper.LoggerLevelPanic, fmt.Sprint("Error when unmarshal configuration file : ", err.Error()), err)
	}
	helper.Logger(helper.LoggerLevelInfo, "Succeed when unmarshal configuration file", err)
	return
}

func InitContainer() *Container {
	app := App{
		ServerHost: v.GetString("SERVER_HOST"),
		ServerPort: v.GetInt("SERVER_PORT"),
	}

	db := DatabaseInit(v)
	articleRepo := repository.NewArticleRepo(db)

	articleService := service.NewArticleService(articleRepo)

	return &Container{
		Db:             db,
		App:            &app,
		ArticleService: articleService,
	}
}

func InitMockViper() *viper.Viper {
	vMock := viper.New()
	vMock.SetEnvPrefix("TEST")
	vMock.AutomaticEnv()
	return vMock
}

func InitMockContainer(v *viper.Viper) *Container {
	app := App{
		ServerHost: v.GetString("SERVER_HOST"),
		ServerPort: v.GetInt("SERVER_PORT"),
	}

	db := DatabaseInit(v)
	articleRepo := repository.NewArticleRepo(db)

	articleService := service.NewArticleService(articleRepo)

	return &Container{
		Db:             db,
		App:            &app,
		ArticleService: articleService,
	}
}

func init() {
	v = viper.New()
	v.AutomaticEnv()

	helper.Logger(helper.LoggerLevelInfo, "Succeed read configuration file: ", nil)
}
