package main

import (
	"context"
	"fmt"
	"log"

	"github.com/weilyuwang/grpc-go/blog/blogpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Blog Client")

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	// Create the blog
	fmt.Println("Creating the blog")
	blog := &blogpb.Blog{
		AuthorId: "some_id",
		Title:    "My First Blog",
		Content:  "Content of my first blog",
	}
	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected error: %v\n", err)
		return
	}

	fmt.Printf("Blog has been created: %v\n", res.Blog)
}
