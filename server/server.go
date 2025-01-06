// solo creamos este archivo para ser usado como una simple interfaz del servidor web
// en este caso usare Echo, pero eso puede cambiar y este archivo permaneceria igual
// La idea de usar este tipo de interfaces es desacoplar el programa, para poder cambiar partes
// con el minimo impacto posible
package server

type Server interface {
	Start()
}
