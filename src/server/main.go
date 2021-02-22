package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"

	// "log"

	jwt "github.com/Madhur/Note-it/src/jwt"
	proto "github.com/Madhur/Note-it/src/proto"
	_ "github.com/go-sql-driver/mysql"

	// jwt "github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type server struct{}

var db *sql.DB
var err error

func (s *server) Login(ctx context.Context, request *proto.LoginRequest) (*proto.LoginResponse, error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.DataLoss, "failed to get metadata")
	}
	xrid := md["user"]

	fmt.Println(md)
	fmt.Println(xrid)

	userID := request.GetUser()
	password := request.GetPassword()

	if len(userID) == 0 {
		return nil, status.Errorf(400+codes.InvalidArgument, "missing 'user'")
	}
	if userID == "" {
		return nil, status.Errorf(400+codes.InvalidArgument, "empty 'user")
	}

	if len(password) == 0 {
		return nil, status.Errorf(400+codes.InvalidArgument, "missing 'password")
	}
	if password == "" {
		return nil, status.Errorf(400+codes.InvalidArgument, "empty 'password'")
	}

	var message string

	if userID == "madhurmittal275@gmail.com" && password == "Maddy@1234" {
		token, err := jwt.GenerateJWT(userID)
		if err != nil {
			fmt.Println("Failed to generate token")
			return nil, status.Errorf(400+codes.Unauthenticated, "Failed to generate token.")
		}
		message = token
	} else {
		message = "user and password not matched."
	}
	return &proto.LoginResponse{Message: message}, nil
}

// var mySigningKey = []byte("captainjacksparrowsayshi")

// func homePage(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Hello World")
// 	fmt.Println("Endpoint Hit: homePage")

// }

// func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		if r.Header["Token"] != nil {

// 			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
// 				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 					return nil, fmt.Errorf("There was an error")
// 				}
// 				return mySigningKey, nil
// 			})

// 			fmt.Println(token)

// 			if err != nil {
// 				fmt.Fprintf(w, err.Error())
// 			}

// 			if token.Valid {
// 				endpoint(w, r)
// 			}
// 		} else {

// 			fmt.Fprintf(w, "Not Authorized")
// 		}
// 	})
// }

// func handleRequests() {
// 	http.Handle("/", isAuthorized(homePage))
// 	log.Fatal(http.ListenAndServe(":9000", nil))
// }

// func main() {
// 	handleRequests()
// }

func main() {
	listener, err := net.Listen("tcp", ":9000")

	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer()
	proto.RegisterAddServiceServer(srv, &server{})
	reflection.Register(srv)

	db, err = sql.Open("mysql", "maddy:maddy@tcp(127.0.0.1:3306)/noteit")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Database connected succesfully")
	}

	fmt.Println("Serving gRPC on 0.0.0.0:9000")
	if e := srv.Serve(listener); e != nil {
		panic(err)
	}
}
