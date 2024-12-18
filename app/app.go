package app

import (
	"log"
	"net/http"

	"github.com/josuebrunel/sportdropin/account"
	"github.com/josuebrunel/sportdropin/app/config"
	"github.com/josuebrunel/sportdropin/group"
	"github.com/josuebrunel/sportdropin/pkg/view"
	"github.com/josuebrunel/sportdropin/pkg/view/base"
	"github.com/josuebrunel/sportdropin/pkg/xsession"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

type App struct {
	Opts config.Config
}

func NewApp() App {
	opts := config.NewConfig()
	return App{Opts: opts}
}

func (a App) Run() {
	// pocket base app
	app := pocketbase.New()
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		ctx := app.RootCmd.Context()

		e.Router.Use(middleware.Logger())
		e.Router.Use(middleware.CORS())
		e.Router.Use(middleware.Recover())
		e.Router.Use(xsession.LoadAndSave(xsession.SessionManager))

		e.Router.Static("/static", "public")
		e.Router.GET("/", func(c echo.Context) error { return view.Render(c, http.StatusOK, base.Index(), nil) })
		groupHandler := group.NewGroupHandler(app.Dao(), app.Settings().Meta.AppUrl)
		g := e.Router.Group("/group")
		g.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
			TokenLookup: "form:csrf,header:csrf",
		}))
		// GROUPS
		g.AddRoute(echo.Route{Method: http.MethodGet, Path: "", Handler: groupHandler.List(ctx), Name: "group.list"})
		g.AddRoute(echo.Route{Method: http.MethodGet, Path: "/:groupid", Handler: groupHandler.Get(ctx), Name: "group.get"})
		g.AddRoute(echo.Route{Method: http.MethodGet, Path: "/create", Handler: groupHandler.Create(ctx), Name: "group.create"})
		g.AddRoute(echo.Route{Method: http.MethodPost, Path: "/create", Handler: groupHandler.Create(ctx), Name: "group.created"})
		g.AddRoute(echo.Route{Method: http.MethodGet, Path: "/:groupid/edit", Handler: groupHandler.Update(ctx), Name: "group.update"})
		g.AddRoute(echo.Route{Method: http.MethodPatch, Path: "/:groupid/edit", Handler: groupHandler.Update(ctx), Name: "group.update"})
		g.AddRoute(echo.Route{Method: http.MethodDelete, Path: "/:groupid", Handler: groupHandler.Delete(ctx), Name: "group.delete"})
		// SEASONS
		g.AddRoute(echo.Route{Method: http.MethodGet, Path: "/:groupid/season/create", Handler: groupHandler.SeasonCreate(ctx), Name: "season.create"})
		g.AddRoute(echo.Route{Method: http.MethodPost, Path: "/:groupid/season/create", Handler: groupHandler.SeasonCreate(ctx), Name: "season.create"})
		g.AddRoute(echo.Route{Method: http.MethodGet, Path: "/:groupid/seasons", Handler: groupHandler.SeasonList(ctx), Name: "season.list"})
		g.AddRoute(echo.Route{Method: http.MethodGet, Path: "/:groupid/season/:seasonid/edit", Handler: groupHandler.SeasonEdit(ctx), Name: "season.edit"})
		g.AddRoute(echo.Route{Method: http.MethodPatch, Path: "/:groupid/season/:seasonid/edit", Handler: groupHandler.SeasonEdit(ctx), Name: "season.edit"})
		g.AddRoute(echo.Route{Method: http.MethodDelete, Path: "/:groupid/season/:seasonid", Handler: groupHandler.SeasonDelete(ctx), Name: "season.delete"})
		// MEMBERS
		g.AddRoute(echo.Route{Method: http.MethodGet, Path: "/:groupid/member/create", Handler: groupHandler.MemberCreate(ctx), Name: "member.create"})
		g.AddRoute(echo.Route{Method: http.MethodPost, Path: "/:groupid/member/create", Handler: groupHandler.MemberCreate(ctx), Name: "member.created"})
		g.AddRoute(echo.Route{Method: http.MethodGet, Path: "/:groupid/members", Handler: groupHandler.MemberList(ctx), Name: "member.list"})
		g.AddRoute(echo.Route{Method: http.MethodGet, Path: "/:groupid/member/:memberid/edit", Handler: groupHandler.MemberEdit(ctx), Name: "member.edit"})
		g.AddRoute(echo.Route{Method: http.MethodPatch, Path: "/:groupid/member/:memberid/edit", Handler: groupHandler.MemberEdit(ctx), Name: "member.edit"})
		g.AddRoute(echo.Route{Method: http.MethodDelete, Path: "/:groupid/member/:memberid", Handler: groupHandler.MemberDelete(ctx), Name: "member.delete"})
		// STATS
		g.AddRoute(echo.Route{Method: http.MethodGet, Path: "/:groupid/stat/create", Handler: groupHandler.StatCreate(ctx), Name: "stat.create"})
		g.AddRoute(echo.Route{Method: http.MethodPost, Path: "/:groupid/stat/create", Handler: groupHandler.StatCreate(ctx), Name: "stat.create"})
		g.AddRoute(echo.Route{Method: http.MethodGet, Path: "/:groupid/stat/", Handler: groupHandler.StatList(ctx), Name: "stat.list"})
		// ACCOUNTS
		accountHandler := account.NewAccountHandler(app.App.Settings().Meta.AppUrl)
		a := e.Router.Group("/account")
		a.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
			TokenLookup: "form:csrf,header:csrf",
		}))
		a.AddRoute(echo.Route{Method: http.MethodGet, Path: "/login", Handler: accountHandler.Login(ctx), Name: "account.login"})
		a.AddRoute(echo.Route{Method: http.MethodPost, Path: "/login", Handler: accountHandler.Login(ctx), Name: "account.login"})
		a.AddRoute(echo.Route{Method: http.MethodGet, Path: "/logout", Handler: accountHandler.Logout(ctx), Name: "account.logout"})
		a.AddRoute(echo.Route{Method: http.MethodGet, Path: "/:accountid", Handler: accountHandler.Get(ctx), Name: "account.get",
			Middlewares: []echo.MiddlewareFunc{xsession.LoginRequired}})
		a.AddRoute(echo.Route{Method: http.MethodGet, Path: "/register", Handler: accountHandler.Create(ctx), Name: "account.register"})
		a.AddRoute(echo.Route{Method: http.MethodPost, Path: "/register", Handler: accountHandler.Create(ctx), Name: "account.register"})
		a.AddRoute(echo.Route{Method: http.MethodGet, Path: "/:accountid/edit", Handler: accountHandler.Update(ctx), Name: "account.update",
			Middlewares: []echo.MiddlewareFunc{xsession.LoginRequired}})
		a.AddRoute(echo.Route{Method: http.MethodPatch, Path: "/:accountid/edit", Handler: accountHandler.Update(ctx), Name: "account.update",
			Middlewares: []echo.MiddlewareFunc{xsession.LoginRequired}})
		a.AddRoute(echo.Route{Method: http.MethodGet, Path: "/:accountid/groups", Handler: accountHandler.Groups(ctx), Name: "account.groups",
			Middlewares: []echo.MiddlewareFunc{xsession.LoginRequired}})
		// a.AddRoute(echo.Route{Method: http.MethodDelete, Path: "/:accountid", Handler: accountHandler.Delete(ctx), Name: "account.delete"})
		return nil
	})

	if err := app.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatal("shutting down the server")
	}
}
