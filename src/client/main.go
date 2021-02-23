package main

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"time"

// 	jwt "github.com/dgrijalva/jwt-go"
// )

// var mySigningKey = []byte("captainjacksparrowsayshi")

// func homePage(w http.ResponseWriter, r *http.Request) {
// 	validToken, err := GenerateJWT()
// 	if err != nil {
// 		fmt.Println("Failed to generate token")
// 	}

// 	client := &http.Client{}
// 	req, _ := http.NewRequest("GET", "http://localhost:9000/", nil)
// 	req.Header.Set("Token", validToken)
// 	res, err := client.Do(req)
// 	if err != nil {
// 		fmt.Fprintf(w, "Error: %s", err.Error())
// 	}

// 	body, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Fprintf(w, string(body))
// }

// //GenerateJWT is to create access token
// func GenerateJWT() (string, error) {
// 	token := jwt.New(jwt.SigningMethodHS256)

// 	claims := token.Claims.(jwt.MapClaims)

// 	claims["authorized"] = true
// 	claims["client"] = "Elliot Forbes"
// 	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

// 	tokenString, err := token.SignedString(mySigningKey)

// 	if err != nil {
// 		fmt.Errorf("Something Went Wrong: %s", err.Error())
// 		return "", err
// 	}

// 	return tokenString, nil
// }

// func handleRequests() {
// 	http.HandleFunc("/", homePage)

// 	log.Fatal(http.ListenAndServe(":9002", nil))
// }

// func main() {
// 	handleRequests()
// }

import (
	"context"
	"fmt"
	"log"
	"net/http"

	proto "github.com/Madhur/Note-it/src/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"
)

func serveSwagger(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Header["Token"])
	http.ServeFile(w, r, "../swaggerui/swagger.json")
}

func main() {
	conn, err := grpc.DialContext(
		context.Background(),
		"localhost:9000",
		grpc.WithBlock(),
		grpc.WithInsecure(),
	)

	if err != nil {
		panic(err)
	}

	// Register grpc-gateway
	gwmux := runtime.NewServeMux()
	// Register AddService
	err = proto.RegisterAddServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	// Serve the swagger-ui and swagger file
	mux := http.NewServeMux()
	mux.Handle("/", gwmux)
	mux.HandleFunc("../swaggerui/swagger.json", serveSwagger)
	fs := http.FileServer(http.Dir("../swaggerui"))
	mux.Handle("/swaggerui/", http.StripPrefix("/swaggerui/", fs))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	})

	gwServer := &http.Server{
		Addr:    ":9009",
		Handler: c.Handler(mux),
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:9009")
	log.Fatalln(gwServer.ListenAndServe())

}
