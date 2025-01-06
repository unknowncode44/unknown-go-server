// solo creamos este archivo para ser usado como una simple interfaz de la base de datos
// en este caso usare postgres, pero eso puede cambiar y este archivo permaneceria igual
// La idea de usar este tipo de interfaces es desacoplar el programa, para poder cambiar partes
// con el minimo impacto posible
package database

// el ORM que usamos es GORM (mas info en // gorm.io //)
import "gorm.io/gorm"

type Database interface {
	GetDb() *gorm.DB
}
