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
	createBlogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected error: %v\n", err)
	}
	fmt.Printf("Blog has been created: %v\n", createBlogRes)
	blogID := createBlogRes.GetBlog().GetId()

	// Read Blog
	fmt.Println("Reading the blog")

	// read blog with an invalid id
	_, errRead := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId: "SOME_INVALID_ID",
	})
	if errRead != nil {
		fmt.Printf("Error happened while reading blog: %v\n", errRead)
	}

	// read blog with a valid id
	readBlogReq := blogpb.ReadBlogRequest{
		BlogId: blogID,
	}
	readBlogRes, readBlogErr := c.ReadBlog(context.Background(), &readBlogReq)
	if readBlogErr != nil {
		fmt.Printf("Error happened while reading blog: %v\n", readBlogErr)
	}

	fmt.Printf("Blog was read: %v\n", readBlogRes)
}
