package main

import(
	
    "google.golang.org/grpc"
	userpb "github.com/dojinkimm/go-grpc-example/protos/v1/user"
    "log"
    "time"
)

func main(){
 conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure(), grpc.WithBlock())
 if err != nil {
	log.Fatalf("did not connect: %v", err)
}
defer conn.Close()
c := userpb.NewUserClient(conn)

ctx, cancel := context.WithTimeout(context.Background(), time.Second)
defer cancel()

var in = userpb.ListUsersRequest{} 
r, err := c.ListUsers(ctx,&in)
if err != nil {
	log.Fatalf("could not request: %v", err)
}

log.Printf("Config: %v",r)
}