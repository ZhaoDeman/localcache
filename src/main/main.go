package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"localcache/src/config"
	"localcache/src/model"
	"localcache/src/service"
	"os"

	"time"
)

var ConfigFilePath = flag.String("config_file", "../config.json", "config配置地址")

func InitFile() {
	if *ConfigFilePath == "" {
		fmt.Println("请指定配置文件地址")
	}
	fmt.Println("configFilePath:", *ConfigFilePath)
	data, err := os.ReadFile(*ConfigFilePath)
	if err != nil {
		_ = fmt.Errorf("%s:%s", "config.json读取失败：", err.Error())
		return
	}
	err = json.Unmarshal(data, &config.ConfInstance)
	if err != nil {
		_ = fmt.Errorf("%s:%s", "config.json文件解析失败：", err.Error())
	}
	fmt.Printf("dataSource:用户名:%s,密码：%s,数据库：%s\n", config.ConfInstance.DataSource.User, config.ConfInstance.DataSource.Password, config.ConfInstance.DataSource.Database)
}

func Init() {
	InitFile()
	config.InitDb()
}


func main() {
	flag.Parse()
	Init()

	ctx := context.Background()
	service.InitOrganizationCache(ctx,
		service.ICacheConfig{
		    LoadFunc: model.GetMap,
			CacheName: "公司",
			LoadInterval: 20 * time.Minute,
		})
	service.InitAppCache(ctx, service.ICacheConfig{
		LoadFunc: model.GetAppMap,
		CacheName: "应用",
		LoadInterval: 20 * time.Minute,
	})
	service.ICacheReLoadFunc(ctx, service.OrganizationCache)
	service.ICacheReLoadFunc(ctx, service.AppCache)

	tenant := model.Tenant{}
	list, _ := tenant.QueryTenantNameMap()

	app := model.App{}
	_, data, _ := app.List(map[string]interface{}{})

	go func() {
		for {
			select {
			case <-time.After(1 * time.Second):
				for i := 0; i < 100; i++ {
					go func() {
						for _, v := range list {
							b, _ := service.ICacheGet(service.OrganizationCache, v.Organization)
							str := ""
							err := json.Unmarshal(b, &str)
							if str == "" || err != nil {
								fmt.Println(v.Organization)
							}
						}
					}()
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case <-time.After(1 * time.Second):
				for i := 0; i < 100; i++ {
					go func() {
						for _, v := range data {
							b, _ := service.ICacheGet(service.AppCache, v.Organization+v.AppId)
							str := ""
							err := json.Unmarshal(b, &str)
							if str == "" || err != nil {
								fmt.Println("app:", v.Organization)
							}
						}
					}()
				}
			}
		}
	}()

	select {}
}
