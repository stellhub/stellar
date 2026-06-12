package boot

import (
	goredis "github.com/redis/go-redis/v9"
	mysqlclient "github.com/stellhub/stellar/clients/mysql"
	postgresqlclient "github.com/stellhub/stellar/clients/postgresql"
	redisclient "github.com/stellhub/stellar/clients/redis"
)

func (a *App) RedisClient() (*goredis.Client, bool) {
	return GetAs[*goredis.Client](a.Registry(), redisclient.DefaultClientName)
}

func (a *App) MySQLDB() (*mysqlclient.DB, bool) {
	return GetAs[*mysqlclient.DB](a.Registry(), mysqlclient.DefaultDBName)
}

func (a *App) PostgreSQLDB() (*postgresqlclient.DB, bool) {
	return GetAs[*postgresqlclient.DB](a.Registry(), postgresqlclient.DefaultDBName)
}
