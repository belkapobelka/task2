package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	pb "task2/server/proto/consignment"
)

const (
	port = ":50051"
)

type Service struct{}

var solutions pb.Solutions

func (s *Service) Solve(ctx context.Context, coeffs *pb.Coefficients) (*pb.Solution, error) {
	var solution pb.Solution
	solution = countNRoots(coeffs)
	solutions.Solutions = append(solutions.Solutions, &solution)
	return &solution, nil
}

func (s *Service) GetAll(context.Context, *pb.GetRequest) (*pb.Solutions, error) {
	return &solutions, nil
}

func countNRoots(coeffs *pb.Coefficients) pb.Solution {
	var s = pb.Solution{
		Coefs: &pb.Coefficients{
			A: coeffs.A,
			B: coeffs.B,
			C: coeffs.C,
		},
		NRoots: 0,
	}
	if (coeffs.A == 0 && coeffs.B != 0) || (coeffs.A != 0 && coeffs.C == 0 && coeffs.B == 0) || (coeffs.A == coeffs.B && coeffs.C == 0) {
		s.NRoots = 1
		return s
	}
	if coeffs.A == 0 && coeffs.B == 0 {
		s.NRoots = 0
		return s
	}

	disc := coeffs.B*coeffs.B - 4*coeffs.A*coeffs.C
	if disc == 0 {
		s.NRoots = 1
	} else if disc > 0 {
		s.NRoots = 2
	} else {
		s.NRoots = 0
	}
	return s
}

func main() {
	s := &Service{}

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("cannot listen the port %v", err)
	}
	server := grpc.NewServer()
	pb.RegisterSolverServer(server, s)

	reflection.Register(server)

	log.Printf("gRPC server running on port: %v", port)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to server from port: %v", err)
	}
}
