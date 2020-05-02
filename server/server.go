package server

func Init() {
	router := UtsuruRouter()
	router.Run()
}
