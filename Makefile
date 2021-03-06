compile-calculator:
	protoc calculator/calculatorpb/calculator.proto --go_out=plugins=grpc:.

compile-greet:	
	protoc greet/greetpb/greet.proto --go_out=plugins=grpc:.
 
compile-blog:
	protoc blog/blogpb/blog.proto --go_out=plugins=grpc:.


run-server-calculator:
	go run calculator/calculator_server/server.go
 
 
run-client-calculator:
	go run calculator/calculator_client/client.go


run-server-greet:
	go run greet/greet_server/server.go
 

run-client-greet:
	go run greet/greet_client/client.go


run-server-blog:
	go run blog/blog_server/server.go
 

run-client-blog:
	go run blog/blog_client/client.go