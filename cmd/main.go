package main

import (
	initiation "gin-base/init"
	"gin-base/internal/web/router"
)

func main() {
	initiation.Initialize()
	router.Router.Run() // listen and serve on 0.0.0.0:8080*/

	/*rpcListener, err := net.Listen("tcp", "localhost:9090") //开启监听
	if err != nil {
		panic(err)
	}
	rpcServer := grpc.NewServer()
	pb.RegisterFileMetaServiceServer(rpcServer, &server.Server{})

	go func() {
		conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure()) //连接到你的服务地址
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		c := pb.NewFileMetaServiceClient(conn) //返回一个client连接，通过这个连接就可以访问到对应的服务资源，就像一个对象

		ctx, cancel := context.WithTimeout(context.Background(), time.Second) //返回一个client，并设置超时时间
		defer cancel()
		r, err := c.Get(ctx, &pb.InBound{Ids: []int64{2, 3, 4}}) //访问对应的服务器上面的服务方法
		if err != nil {
			panic(err)
		}
		fmt.Println(r.List)
	}()

	if err := rpcServer.Serve(rpcListener); err != nil {
		panic(err)
	}*/

}
