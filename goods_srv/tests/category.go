package main

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
)

func TestGetCategoryList() {
	rsp, err := brandClient.GetAllCategoryList(context.Background(), &empty.Empty{})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	fmt.Println(rsp.JsonData)
}
