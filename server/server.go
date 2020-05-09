package server

// Init initializes the server routes.
func Init() {
	router := UtsuruRouter()
	router.Run()
}
