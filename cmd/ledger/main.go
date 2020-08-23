package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/rafaelturon/ledgerd/pb"
	uuid "github.com/satori/go.uuid"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

const (
	event     = "account.created"
	aggregate = "account"
	grpcUri   = "localhost:50051"
)

type rpcClient interface {
	createAccount(account pb.AccountCreateCommand) error
}
type grpcClient struct {
}

// createAccount calls the CreateEvent RPC
func (gc grpcClient) createAccount(account pb.AccountCreateCommand) error {
	conn, err := grpc.Dial(grpcUri, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewEventStoreClient(conn)
	accountJSON, _ := json.Marshal(account)

	event := &pb.Event{
		EventId:       uuid.NewV4().String(),
		EventType:     event,
		AggregateId:   account.AccountId,
		AggregateType: aggregate,
		EventData:     string(accountJSON),
		Stream:        "Accounts",
	}

	resp, err := client.CreateEvent(context.Background(), event)
	if err != nil {
		return fmt.Errorf("error from RPC server: %w", err)
	}
	if resp.IsSuccess {
		return nil
	}
	return errors.New("error from RPC server")
}

type accountHandler struct {
	rpc rpcClient
}

func (h accountHandler) createAccount(w http.ResponseWriter, r *http.Request) {
	var account pb.AccountCreateCommand
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		http.Error(w, "Invalid Account Data", 500)
		return
	}
	aggregateID := uuid.NewV4().String()
	account.AccountId = aggregateID
	account.Status = "Pending"
	account.CreatedOn = time.Now().Unix()
	err = h.rpc.createAccount(account)
	if err != nil {
		log.Print(err)
		http.Error(w, "Failed to create Account", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	j, _ := json.Marshal(account)
	w.Write(j)
}

func initRoutes() *mux.Router {
	router := mux.NewRouter()
	h := accountHandler{
		rpc: grpcClient{},
	}
	router.HandleFunc("/api/account", h.createAccount).Methods("POST")
	return router
}
func main() {
	// Create the Server
	server := &http.Server{
		Addr:    ":3000",
		Handler: initRoutes(),
	}
	log.Println("HTTP Sever listening...")
	// Running the HTTP Server
	server.ListenAndServe()
}
