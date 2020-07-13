package main

import (
	"context"
	"fmt"
	"io"
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

	// -------------- Create the blog --------------
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

	// -------------- Read Blog --------------
	fmt.Println("Reading the blog")

	// read blog with an invalid id
	_, errRead := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId: "5f0bc73ceef9a6fa8c71d476",
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

	// -------------- Update Blog --------------
	fmt.Println("Updating the blog")

	newBlog := &blogpb.Blog{
		Id:       blogID,
		AuthorId: "changed_author_id",
		Title:    "My First Blog (edited)",
		Content:  "Content of my first blog with some awesome editions",
	}

	updateBlogReq := blogpb.UpdateBlogRequest{
		Blog: newBlog,
	}
	updateBlogRes, updateBlogErr := c.UpdateBlog(context.Background(), &updateBlogReq)
	if updateBlogErr != nil {
		fmt.Printf("Error happened while updateing blog: %v\n", updateBlogErr)
	}
	fmt.Printf("Blog was updated: %v\n", updateBlogRes)

	// -------------- Delete Blog --------------
	fmt.Println("Deleting the blog")

	deleteBlogReq := blogpb.DeleteBlogRequest{
		BlogId: blogID,
	}
	deleteBlogRes, deleteBlogErr := c.DeleteBlog(context.Background(), &deleteBlogReq)
	if deleteBlogErr != nil {
		fmt.Printf("Error happened while deleting blog: %v\n", deleteBlogErr)
	}
	fmt.Printf("Blog was deleted: %v\n", deleteBlogRes)

	// -------------- List All Blogs --------------
	fmt.Println("Listing all the blogs")

	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("Error while calling ListBlog RPC: %v\n", err)
		return
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while receiving blog from server stream: %v\n", err)
		}
		fmt.Println(res.GetBlog())
	}

}
