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
	pool *pgxpool.Pool
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

func (client *Client) Init(ctx context.Context, basicConfigFilePath string, dbName string, maxConnections int32) (err error) {

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
	client.pool, err = pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		log.Println("error connecting to the database: ", err)
		return err
	}
	log.Printf("connect crdb success\n")
	return nil
}

func (client *Client) Query(ctx context.Context, sql string, args ...interface{})(result []any, err error){

	connection, err := client.pool.Acquire(ctx)
	defer connection.Release()
	if err != nil {
		log.Printf("dbpool acquire fail: %v", err)
		return 
	}

	log.Printf("begin query: %s, args: %v\n", sql, args)
	rows, err := connection.Query(ctx, sql, args...)
	if err != nil {
		log.Printf("err: %v\n", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		values, valuesErr := rows.Values()
		if valuesErr != nil{
			err = valuesErr
			log.Printf("values error: %v", err)
			return
		}
		result = append(result, values)
	}
	err = rows.Err()
	if err != nil {
		log.Printf("rows error: %v", rows.Err())
		return
	}
	log.Printf("result: %v\n", result)
	return
}

func (client *Client) Exec(ctx context.Context, sql string, args ...interface{})(err error){

	connection, err := client.pool.Acquire(ctx)
	defer connection.Release()
	if err != nil {
		log.Printf("dbpool acquire fail: %v", err)
		return 
	}

	log.Printf("begin exec: %s, args: %v\n", sql, args)
	commondTag, err := connection.Exec(ctx, sql, args...)
	if err != nil{
		log.Printf("err: %v\n", err)
		return
	}
	log.Printf("commondTag: %s\n", commondTag.String())
	return
}