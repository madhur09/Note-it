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
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type server struct{}

var db *sql.DB
var err error

func (s *server) Login(ctx context.Context, request *proto.LoginRequest) (*proto.LoginResponse, error) {

	// md, ok := metadata.FromIncomingContext(ctx)
	// if !ok {
	// 	return nil, status.Errorf(codes.DataLoss, "failed to get metadata")
	// }
	// xrid := md["user"]
	// fmt.Println(md["User"])
	// fmt.Println(md["user"])
	// fmt.Println(md["grpcgateway-User"])
	// fmt.Println(md["grpcgateway-user"])
	// fmt.Println(md["Grpcgateway-User"])
	// fmt.Println(md.Get("user"))
	// fmt.Println(md.Get("User"))
	// fmt.Println(md.Get("grpcgateway-user"))
	// fmt.Println(md.Get("grpcgateway-User"))
	// fmt.Println(md)
	// fmt.Println(md.Get("grpcgateway-content-type"))
	// fmt.Println(md["grpcgateway-content-type"])

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

func (s *server) Show(ctx context.Context, request *proto.PageLimitRequest) (*proto.PageLimitResponse, error) {
	limit := request.GetLimit()
	pageNo := request.GetPageNo()

	statement := fmt.Sprintf("SELECT * FROM notes where isdeleted=0 limit %d, %d", pageNo*limit, limit)
	rows, _ := db.Query(statement)

	var rowCount int32
	count := db.QueryRow("SELECT count(*) FROM notes where isdeleted=0")
	count.Scan(&rowCount)

	var result []*proto.NoteObject

	for rows.Next() {
		var res proto.NoteObject
		rows.Scan(&res.Id, &res.Title, &res.Description, &res.IsPinned, &res.Color)
		result = append(result, &res)
	}

	return &proto.PageLimitResponse{Notes: result, TotalRows: rowCount}, nil
}

func (s *server) Add(ctx context.Context, request *proto.CreateNoteRequest) (*proto.NoteResponse, error) {
	db.Exec("Insert into notes values(?,?,?,?,?,?)", request.Title, request.Description, request.IsPinned, request.Color, 0)
	rows, _ := db.Query("select * from notes where isdeleted=0")

	var result []*proto.NoteObject

	for rows.Next() {
		var res proto.NoteObject
		rows.Scan(&res.Id, &res.Title, &res.Description, &res.IsPinned, &res.Color)
		result = append(result, &res)
	}

	return &proto.NoteResponse{Notes: result}, nil
}

func (s *server) Update(ctx context.Context, request *proto.NoteIdRequest) (*proto.NoteResponse, error) {
	noteID := request.GetId()

	db.Exec("update notes set title=?,description=?,ispinned=?,color=? where id=?", request.Title, request.Description, request.IsPinned, request.Color, noteID)

	rows, _ := db.Query("select * from notes where isdeleted=0")

	var result []*proto.NoteObject

	for rows.Next() {
		var res proto.NoteObject
		rows.Scan(&res.Id, &res.Title, &res.Description, &res.IsPinned, &res.Color)
		result = append(result, &res)
	}
	return &proto.NoteResponse{Notes: result}, nil
}

func (s *server) Delete(ctx context.Context, request *proto.NoteId) (*proto.NoteResponse, error) {
	noteID := request.GetId()

	sqlStatement := "update notes set isdeleted=1 WHERE id = ?"
	_, err = db.Exec(sqlStatement, noteID)
	if err != nil {
		panic(err)
	}

	rows, _ := db.Query("select * from notes  where isdeleted=0")

	var result []*proto.NoteObject

	for rows.Next() {
		var res proto.NoteObject
		rows.Scan(&res.Id, &res.Title, &res.Description, &res.IsPinned, &res.Color)
		result = append(result, &res)
	}
	return &proto.NoteResponse{Notes: result}, nil
}

func (s *server) Search(ctx context.Context, request *proto.SearchNote) (*proto.NoteResponse, error) {
	item := request.GetItem()

	query := fmt.Sprintf("SELECT * FROM notes WHERE title LIKE '%%%s%%' and isdeleted=0", item)

	rows, _ := db.Query(query)

	var result []*proto.NoteObject

	for rows.Next() {
		var res proto.NoteObject
		rows.Scan(&res.Id, &res.Title, &res.Description, &res.IsPinned, &res.Color)
		result = append(result, &res)
	}
	return &proto.NoteResponse{Notes: result}, nil
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
