package api

import (
	db "danielsxiong/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server serve http requests for banking service
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// NewServer creates a new http server and setup routing
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)

	server.router = router
	return server
}

// Start runs the http server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
