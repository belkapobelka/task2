package main

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	pb "task2/server/proto/consignment"
)

const (
	address = "localhost:50051"
	doc1    = "test1.json"
	doc2    = "test2.json"
	doc3    = "test3.json"
)

var files = []string{doc1, doc2, doc3}

func main() {
	connection, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("cannot cpnnact to port: %v", err)
	}
	defer connection.Close()
	client := pb.NewSolverClient(connection)

	for _, file := range files {
		eqn, err := readJSON(file)
		if err != nil {
			log.Fatalf("cannot read file %v %v", file, err)
		}
		resp, err := client.Solve(context.Background(), &pb.Coefficients{
			A: eqn.A,
			B: eqn.B,
			C: eqn.C,
		})
		if err != nil {
			log.Fatalf("cannot get solution %v", err)
		}
		fmt.Println(resp.NRoots)
	}

	res, err := client.GetAll(context.Background(), &pb.GetRequest{})
	if err != nil {
		log.Fatalf("cannot get allsolutions %v", err)
	}
	for _, sol := range res.Solutions {
		fmt.Printf("A:%v B:%v C:%v Nroots:%v\r\n", sol.Coefs.A, sol.Coefs.B, sol.Coefs.C, sol.NRoots)
	}
}

func readJSON(file string) (*pb.Coefficients, error) {
	var coeff *pb.Coefficients
	fileBody, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(fileBody, &coeff)
	if err != nil {
		return nil, err
	}
	return coeff, nil
}
