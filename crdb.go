package crdb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nicelogic/config"
)

type Client struct {
	Pool *pgxpool.Pool
}

type clientConfig struct {
	User_name              string
	Pwd                   string
	Host                  string
	Port                  int
	Ssl_mode               string
	Ssl_root_cert_path       string
	Max_connection_idle_time int
}

func (client *Client) Init(ctx context.Context,
	 basicConfigFilePath string,
	 dbName string,
	 maxConnections int32) (err error) {

	clientConfig := clientConfig{}
	err = config.Init(basicConfigFilePath, &clientConfig)
	if err != nil {
		log.Println("config init fail, err: ", err)
		return err
	}
	log.Printf("%#v\n", clientConfig)

	url := fmt.Sprintf("postgresql://%s:%s@%s:%v/%s?sslmode=%s&sslrootcert=%s",
		clientConfig.User_name,
		clientConfig.Pwd,
		clientConfig.Host,
		clientConfig.Port,
		dbName,
		clientConfig.Ssl_mode,
		clientConfig.Ssl_root_cert_path)
	config, err := pgxpool.ParseConfig(url)
	config.MaxConns = maxConnections
	config.MaxConnIdleTime = time.Duration(clientConfig.Max_connection_idle_time) * time.Second
	if err != nil {
		log.Println("error configuring the database: ", err)
		return err
	}
	client.Pool, err = pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		log.Println("error connecting to the database: ", err)
		return err
	}
	log.Printf("connect crdb success\n")
	return nil
}
