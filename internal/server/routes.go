package server

import (
	"GoHTMX/frontend"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/a-h/templ"
	"github.com/gorilla/sessions"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	SESSION_SECRET = os.Getenv("SESSION_SECRET")
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data any, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Static("frontend"))
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(SESSION_SECRET))))

	e.Renderer = &TemplateRenderer{
		templates: template.Must(template.ParseGlob("frontend/*.html")),
	}


	e.GET("/", s.NotFound)
	e.GET("/home", s.IndexTempl())
	e.GET("/about", s.AboutPage)
	e.GET("/login", s.LoginPage)
	e.GET("/signup", s.SignupPage)
	e.GET("/*", s.NotFound)

	e.POST("/api/increment", s.Increment)
	e.POST("/api/signup", s.CreateUser)
	e.POST("/api/login", s.Login)
	e.POST("/api/logout", s.Logout)
	e.POST("/api/get-list", s.GetList)
	e.POST("/api/get-list-length", s.GetListLength)
	return e
}

type IndexPageData struct {
	Username string
}

// func (s *Server) IndexPage(c echo.Context) error {
// 	if !s.IsLoggedIn(c) {
// 		return c.Redirect(http.StatusSeeOther, "/login")
// 	}
//
// 	sess, err := session.Get("login", c)
// 	if err != nil {
// 		return err
// 	}
//
// 	data := IndexPageData{
// 		Username: sess.Values["username"].(string),
// 	}
//
// 	return c.Render(http.StatusOK, "index.html", data)
// }

var indexData frontend.IndexPageData

func (s *Server) IndexTempl() echo.HandlerFunc {
	return func(c echo.Context) error {
		if !s.IsLoggedIn(c) {
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		sess, err := session.Get("login", c)
		if err != nil {
			return err
		}

		indexData = frontend.IndexPageData{
			Username:      sess.Values["username"].(string),
			NumberOfItems: "0",
		}

		return echo.WrapHandler(templ.Handler(frontend.Index(indexData)))(c)
	}
}

type AboutPageData struct {
	Count int
}

var data AboutPageData

func (s *Server) AboutPage(c echo.Context) error {
	if !s.IsLoggedIn(c) {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	session, err := session.Get("login", c)
	if err != nil {
		return err
	}

	username := session.Values["username"].(string)


	count := s.db.GetCount(username)
	if err != nil {
		return err
	}

	data.Count = count

	return c.Render(http.StatusOK, "about.html", data)
}

func (s *Server) Increment(c echo.Context) error {
	if !s.IsLoggedIn(c) {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	session, err := session.Get("login", c)
	if err != nil {
		return err
	}

	username := session.Values["username"].(string)

	err = s.db.Increment(username)
	if err != nil {
		return err
	}

	data.Count = s.db.GetCount(username)

	return c.HTML(http.StatusOK, fmt.Sprintf("%d", data.Count))
}

func (s *Server) LoginPage(c echo.Context) error {
	if s.IsLoggedIn(c) {
		return c.Redirect(http.StatusSeeOther, "/home")
	}
	return c.File("frontend/login.html")
}

func (s *Server) SignupPage(c echo.Context) error {
	if s.IsLoggedIn(c) {
		return c.Redirect(http.StatusSeeOther, "/home")
	}
	return c.File("frontend/signup.html")
}

func (s *Server) NotFound(c echo.Context) error {
	return c.File("frontend/404.html")
}


func (s *Server) IsLoggedIn(c echo.Context) bool {
	sess, err := session.Get("login", c)
	if err != nil {
		return false
	}

	if sess.Values["username"] == nil {
		return false
	}

	return true
}

func (s *Server) CreateUser(c echo.Context) error {
	if s.IsLoggedIn(c) {
		return c.Redirect(http.StatusSeeOther, "/home")
	}

	username := c.FormValue("username")
	password := c.FormValue("password")

	err := s.db.CreateUser(username, password)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, "/login")
}

func (s *Server) Login(c echo.Context) error {
	if s.IsLoggedIn(c) {
		return c.Redirect(http.StatusSeeOther, "/home")
	}

	username := c.FormValue("username")
	password := c.FormValue("password")

	err := s.db.Login(username, password)
	if err != nil {
		return err
	}

	sess, err := session.Get("login", c)
	if err != nil {
		return err
	}

	sess.Values["username"] = username
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusSeeOther, "/home")
}

func (s *Server) Logout(c echo.Context) error {
	if !s.IsLoggedIn(c) {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	sess, err := session.Get("login", c)
	if err != nil {
		return err
	}

	sess.Values["username"] = nil
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusSeeOther, "/login")
}

func (s *Server) GetList(c echo.Context) error {
	if !s.IsLoggedIn(c) {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	return c.HTML(http.StatusOK, `
		<li>Item 1</li>
		<li>Item 2</li>
		<li>Item 3</li>
		<li>Item 4</li>
	`)
}

func (s *Server) GetListLength(c echo.Context) error {
	if !s.IsLoggedIn(c) {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	numberOfItems, err := strconv.Atoi(indexData.NumberOfItems)
	if err != nil {
		return err
	}

	numberOfItems += 4

	indexData.NumberOfItems = strconv.Itoa(numberOfItems)

	return c.HTML(http.StatusOK, indexData.NumberOfItems)
}
