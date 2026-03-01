// en este archivo:
// definimos una estructura para nuestra instancia de conexion a db
// creamos una funcion que devuelve una instancia de la clase Database archivo ./database.go
package database

import (
	"fmt"
	"log"
	"sync"

	"github.com/unknowncode44/unknown-go-server/config"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// creamos un tipo para almacenar la instancia de nuestra db
type postgresDatabase struct {
	Db *gorm.DB
}

// variables de uso en la funcion NewPostgresDatabase
var (
	once       sync.Once
	dbInstance *postgresDatabase
)

// la funcion NewPostgresDatabase toma como parametro nuestra configuracion y devuelve una instancia de Database,
// que es la interfaz que creamos para el manejo de la base de datos
// si en el futuro quisieramos cambiar a otra base de datos, deberiamos crear una funcion similar
// que devuelva un obj del tipo Database
func NewPostgresDatabase(conf *config.Config) Database {
	once.Do(func() {
		// pasamos los argumentos necesarios para conexion de nuestro proyecto (../config/config.go)
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d TimeZone=%s",
			conf.Db.Host,
			conf.Db.User,
			conf.Db.Password,
			conf.Db.DBName,
			conf.Db.Port,
			conf.Db.TimeZone,
		)

		// intentamos conectar
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

		// panic si errores
		if err != nil {
			panic("Fallo La Conexion a la Base de Datos")
		} else {
			fmt.Println("Conexion a DB Exitosa!")
		}

		// migra las entidades existentes
		db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
		err = db.AutoMigrate(&entities.Material{}, &entities.Vendor{}, &entities.VendorMaterial{}, &entities.MaterialCost{}, &entities.Currency{}, &entities.Asset{}, &entities.Location{}, &entities.AssetMovement{}, &entities.AssetSerialCounter{}, &entities.PurchaseOrderSync{})
		if err != nil {
			log.Fatalf("Falla migrando la bd: %v", err)
		} else {
			fmt.Println("Migracion exitosa")
		}

		// caso contrario asignamos la instancia a nuestra variable
		dbInstance = &postgresDatabase{Db: db}
	})
	// y la retornamos
	return dbInstance

}

// para cumplir con lo requerido por GORM, creamos el metodo GetDB() que devuelve el valor Db de la instancia dbInstance
func (p *postgresDatabase) GetDb() *gorm.DB {
	return dbInstance.Db
}
