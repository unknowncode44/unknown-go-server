package server

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/unknowncode44/unknown-go-server/api/routes"
	"github.com/unknowncode44/unknown-go-server/config"
	database "github.com/unknowncode44/unknown-go-server/db"
	"github.com/unknowncode44/unknown-go-server/pkg/company"
)

// la estructura fiberServer la utilizaremos para el servidor, la db y la configuracion
// que importaremos desde la libreria fiber, de nuestro archivo de database y nuestro archivo de configuracion
type fiberServer struct {
	app  *fiber.App
	db   database.Database
	conf *config.Config
}

// la funcion NewFiberServer recibe los parametros de configuracion y database
// y sera la encargada de levantar nuestro servidor en el puerto especificado en el archivo de configuracion
// y utilizar nuestra configuracion de base de datos
func NewFiberServer(conf *config.Config, db database.Database) Server {
	// creamos una nueva instancia de fiber app
	fiberApp := fiber.New()
	// configuramos el servidor para que acepte peticiones de cualquier origen
	fiberApp.Use(cors.New())

	fiberApp.Get("/api/health/v1", func(ctx *fiber.Ctx) error {
		return ctx.Send([]byte("Working Cool!"))
	})

	// repos & services
	companyRepo := company.NewRepo(db.GetDb())
	companyService := company.NewService(companyRepo)

	// global api route
	api := fiberApp.Group("/api/v1")

	routes.CompanyRouter(api, companyService)

	return &fiberServer{
		app:  fiberApp,
		db:   db,
		conf: conf,
	}

}

func (s *fiberServer) Start() {

	// definimos la direccion y puerto del servidor
	serveUrl := fmt.Sprintf(":%d", s.conf.Server.Port)
	log.Fatal(s.app.Listen(serveUrl))

}
