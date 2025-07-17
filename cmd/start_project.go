package cmd

import (
	"log"
	"net/http"

	"github.com/biswasurmi/book-cli/api/handler"
	"github.com/biswasurmi/book-cli/infrastructure/persistance/inmemory"
	"github.com/biswasurmi/book-cli/service"
	"github.com/spf13/cobra"
)

var port string
var auth bool

var startProject = &cobra.Command{
	Use:   "startProject",
	Short: "Start the Book Server",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Starting Book Server on port", port)

		repos := inmemory.GetRepositories()
		services := service.GetServices(repos)
		h := &handler.Handler{
			BookHandler: handler.NewBookHandler(services.BookService),
			UserHandler: handler.NewUserHandler(services.UserService),
		}

		server := handler.CreateNewServer(h, services, auth)
		server.MountRoutes()

		addr := ":" + port
		log.Printf("Server listening on %s", addr)
		if err := http.ListenAndServe(addr, server.Router); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(startProject)
	startProject.PersistentFlags().StringVarP(&port, "port", "p", "8080", "Port to run server")
	startProject.PersistentFlags().BoolVarP(&auth, "auth", "a", true, "Enable basic auth and JWT")
}