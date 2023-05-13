package main

import (
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	os.Exit(m.Run())
}


type myHandler struct {

}

func (mh *myHandler) ServeHTTP( w http.ResponseWriter, r *http.Request){
	// this function would give as type to myHandler struct , we use it where we need type as ServeHTTP function
}