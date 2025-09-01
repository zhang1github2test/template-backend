package main

import "template-backend/cmd/server"

func main() {
	//password, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Print(string(password))
	server.ServerMain()
}
