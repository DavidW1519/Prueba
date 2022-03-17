package main

import (
	"context"
	"flag"
	"strings"

	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/gitextractor/models"
	"github.com/merico-dev/lake/plugins/gitextractor/parser"
	"github.com/merico-dev/lake/plugins/gitextractor/store"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	url := flag.String("url", "", "-url")
	proxy := flag.String("proxy", "", "-proxy")
	id := flag.String("id", "", "-id")
	user := flag.String("user", "", "-user")
	password := flag.String("password", "", "-password")
	output := flag.String("output", "", "-output")
	db := flag.String("db", "", "-db")
	flag.Parse()
	var storage models.Store
	var err error
	if *url == "" {
		panic("url is missing")
	}
	if *id == "" {
		panic("id is missing")
	}
	if *output != "" {
		storage, err = store.NewCsvStore(*output)
		if err != nil {
			panic(err)
		}
	} else if *db != "" {
		database, err := gorm.Open(mysql.Open(*db))
		if err != nil {
			panic(err)
		}
		storage = store.NewDatabase(database)
	} else {
		panic("either specify `-output` or `-db` argument as destination")
	}
	defer storage.Close()
	ctx := context.Background()
	subTaskCtx := helper.NewStandaloneSubTaskContext(
		"git extractor",
		ctx,
		helper.NewDefaultTaskLogger(nil, "git extractor"),
		nil,
	)
	p := parser.NewLibGit2(storage, subTaskCtx)
	if strings.HasPrefix(*url, "http") {
		err = p.CloneOverHTTP(ctx, *url, *id, *user, *password, *proxy)
		if err != nil {
			panic(err)
		}
	}
	if strings.HasPrefix(*url, "/") {
		err = p.LocalRepo(ctx, *url, *id)
		if err != nil {
			panic(err)
		}
	}
}
