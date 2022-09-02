package main

import (
	"fatalisa-public-api/config"
	"fatalisa-public-api/router"
	"github.com/subchen/go-log"
)

func init() {
	config.Init()
}

func main() {
	// any code after routerInit.Run() won't be executed, place it in init_task.go inside config
	log.Info("Starting service")
	routerInit := &router.Config{}
	routerInit.Run()

	//srv := &http.Server{
	//	Addr:    ":8080",
	//	Handler: routerInit.Gin,
	//}
	//
	//run(srv)
}

// this function referenced from Gin documentation
//func run(server *http.Server) {
//	// Initializing the server in a goroutine so that
//	// it won't block the graceful shutdown handling below
//	go func() {
//		if err := server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
//			log.Printf("listen: %s\n", err)
//		}
//	}()
//
//	// Wait for interrupt signal to gracefully shut down the server with
//	// a timeout of 5 seconds.
//	quit := make(chan os.Signal)
//	// kill (no param) default send syscall.SIGTERM
//	// kill -2 is syscall.SIGINT
//	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
//	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
//	<-quit
//	log.Println("Shutting down server...")
//
//	// The context is used to inform the server it has 5 seconds to finish
//	// the request it is currently handling
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	if err := server.Shutdown(ctx); err != nil {
//		log.Fatal("Server forced to shutdown:", err)
//	}
//
//	log.Println("Server exiting")
//}
