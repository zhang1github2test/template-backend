package main

import "template-backend/cmd/server"

// @title template-backend API
// @version 1.0
// @description template-backend管理系统的API文档
// @host localhost:8080
// @BasePath /api
func main() {
	server.ServerMain()
}
