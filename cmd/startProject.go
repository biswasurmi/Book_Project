package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/biswasurmi/book-cli/api/handler"
	"github.com/biswasurmi/book-cli/infrastructure/persistance/inmemory"
	"github.com/biswasurmi/book-cli/service"
	"github.com/go-chi/jwtauth/v5"
	"github.com/spf13/cobra"
)

var port string
var auth bool

var startProject = &cobra.Command{
	Use:   "startProject",
	Short: "Start the Book Server",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Starting Book Server on port", port)

		bookRepo := inmemory.NewBookRepo()
		bookService := service.NewBookService(bookRepo)
		bookHandler := handler.NewBookHandler(bookService)
		tokenAuth := jwtauth.New("HS256", []byte("supersecretkey123"), nil)

		server := handler.CreateNewServer(bookHandler, auth, tokenAuth)
		server.MountRoutes()

		addr := ":" + port
		if err := http.ListenAndServe(addr, server.Router); err != nil {
			fmt.Printf("Server error: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(startProject)
	startProject.PersistentFlags().StringVarP(&port, "port", "p", "8080", "Port to run server")
	startProject.PersistentFlags().BoolVarP(&auth, "auth", "a", true, "Enable basic auth and JWT")
}