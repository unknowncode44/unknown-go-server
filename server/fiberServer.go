package server

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/unknowncode44/unknown-go-server/api/handlers"
	"github.com/unknowncode44/unknown-go-server/api/routes"
	"github.com/unknowncode44/unknown-go-server/config"
	database "github.com/unknowncode44/unknown-go-server/db"
	"github.com/unknowncode44/unknown-go-server/pkg/asset"
	"github.com/unknowncode44/unknown-go-server/pkg/asset_movement"
	"github.com/unknowncode44/unknown-go-server/pkg/auth"
	"github.com/unknowncode44/unknown-go-server/pkg/currency"
	"github.com/unknowncode44/unknown-go-server/pkg/delivery_record"
	"github.com/unknowncode44/unknown-go-server/pkg/location"
	"github.com/unknowncode44/unknown-go-server/pkg/material"
	"github.com/unknowncode44/unknown-go-server/pkg/material_cost"
	"github.com/unknowncode44/unknown-go-server/pkg/material_inventory"
	"github.com/unknowncode44/unknown-go-server/pkg/need"
	"github.com/unknowncode44/unknown-go-server/pkg/sync_purchase_order"
	"github.com/unknowncode44/unknown-go-server/pkg/user"
	"github.com/unknowncode44/unknown-go-server/pkg/vendor"
	"github.com/unknowncode44/unknown-go-server/pkg/vendor_material"
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

	fiberApp.Get("/api/v1/health", func(ctx *fiber.Ctx) error {
		return ctx.Send([]byte("Working Cool!"))
	})

	// repos & services

	// Material (handler se construye más abajo, necesita mi/dr services)
	materialRepo := material.NewRepo(db.GetDb())
	materialService := material.NewService(materialRepo)

	// Vendor
	vendorRepo := vendor.NewRepo(db.GetDb())
	vendorService := vendor.NewService(vendorRepo)
	vendorHandler := handlers.NewVendorHandler(vendorService)

	// VendorMaterial
	vmRepo := vendor_material.NewRepo(db.GetDb())
	vmService := vendor_material.NewService(vmRepo)
	vmHandler := handlers.NewVendorMaterialHandler(vmService)

	// Currency
	currencyRepo := currency.NewRepo(db.GetDb())
	currencyService := currency.NewService(currencyRepo)
	currencyHandler := handlers.NewCurrencyHandler(currencyService)

	// Location (handler se construye más abajo, necesita assetService)
	locationRepo := location.NewRepo(db.GetDb())
	locationService := location.NewService(locationRepo)

	// AssetMovement
	amRepo := asset_movement.NewRepo(db.GetDb())
	amService := asset_movement.NewService(amRepo)
	amHandler := handlers.NewAssetMovementHandler(amService)

	// Asset
	assetRepo := asset.NewRepo(db.GetDb())
	assetService := asset.NewService(assetRepo, materialRepo)
	assetHandler := handlers.NewAssetHandler(assetService, amService)

	// Location handler (necesita assetService para la ruta pública)
	locationHandler := handlers.NewLocationHandler(locationService, assetService)

	// MaterialCost
	mcRepo := material_cost.NewRepo(db.GetDb())
	mcService := material_cost.NewService(mcRepo, vmRepo, vendorRepo, currencyRepo)
	mcHandler := handlers.NewMaterialCostHandler(mcService)

	// MaterialInventory
	miRepo := material_inventory.NewRepo(db.GetDb())
	miService := material_inventory.NewService(miRepo, materialRepo)

	// DeliveryRecord
	drRepo := delivery_record.NewRepo(db.GetDb())
	drService := delivery_record.NewService(drRepo, miRepo)

	miHandler := handlers.NewMaterialInventoryHandler(miService, drService)
	drHandler := handlers.NewDeliveryRecordHandler(drService)

	// Material handler (necesita mi/dr services para la ruta pública)
	materialHandler := handlers.NewMaterialHandler(materialService, miService, drService)

	// Need
	needRepo := need.NewRepo(db.GetDb())
	needService := need.NewService(needRepo, materialRepo)
	needHandler := handlers.NewNeedHandler(needService)

	// SyncPurchaseOrder
	spoRepo := sync_purchase_order.NewRepo(db.GetDb())
	spoService := sync_purchase_order.NewService(spoRepo, materialRepo, vendorRepo)
	spoHandler := handlers.NewSyncOrderHandler(spoService)

	// User / Auth
	userRepo := user.NewRepo(db.GetDb())
	userService := user.NewService(userRepo)
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(userService)

	// global api route
	api := fiberApp.Group("/api/v1")

	// rutas de autenticación: /auth/login es pública (puerta de entrada),
	// /auth/me exige token propio (ver auth_routes.go). Se registran ANTES
	// del middleware global para que el login no exija token.
	routes.AuthRoutes(api, authHandler)

	// a partir de acá, todas las rutas /api/v1/* registradas debajo exigen
	// un Bearer token válido (las /public/* quedan afuera del grupo)
	api.Use(auth.RequireAuth())

	// routes.CompanyRouter(api, companyService)
	routes.MaterialRoutes(api, materialHandler)
	routes.CurrencyRoutes(api, currencyHandler)
	routes.VendorRoutes(api, vendorHandler)
	routes.VendorMaterialRoutes(api, vmHandler)
	routes.MaterialCostRoutes(api, mcHandler)
	routes.LocationRoutes(api, locationHandler)
	routes.AssetRoutes(api, assetHandler)
	routes.AssetMovementRoutes(api, amHandler)
	// register asset-scoped movement endpoints
	routes.AssetMovementByAssetRoutes(api, amHandler)
	routes.MaterialInventoryRoutes(api, miHandler)
	routes.MaterialInventoryByMaterialRoutes(api, miHandler)
	routes.DeliveryRecordRoutes(api, drHandler)
	routes.DeliveryRecordByInventoryRoutes(api, miHandler)
	routes.NeedRoutes(api, needHandler)
	// gestión de usuarios: solo ADMIN (ver user_routes.go)
	routes.UserRoutes(api, userHandler)

	// public routes
	public := fiberApp.Group("/public")
	routes.PublicAssetRoute(public, assetHandler)
	routes.PublicMaterialRoute(public, materialHandler)
	routes.PublicLocationRoute(public, locationHandler)
	routes.PublicSyncPurchaseOrdersRoute(public, spoHandler)

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
