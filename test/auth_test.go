package test

import (
	"brank/core"
	"brank/core/auth"
	"brank/core/storage"
	"brank/repository"
	"fmt"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Setenv("ENV", "test")
	os.Exit(m.Run())
}

func Test__GenerateAppKey(t *testing.T) {
	c := core.NewConfig()
	pg, err := storage.NewPostgres(c)
	if err != nil {
		log.Fatal("postgres conn failed", err)
	}

	r := repository.New(pg)

	a := auth.New(r)

	fmt.Println(a.GenerateAppAccessToken(1, "4RH93HFFUYFBY384V9F3Fbhbe3y34g36w7v37273v2tv32v73v2"))

}
