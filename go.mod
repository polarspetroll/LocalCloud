module github.com/polarspetroll/localcloud

go 1.15

require (
	github.com/go-redis/redis v6.15.9+incompatible // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	login v0.0.0
	upload v0.0.0
)

replace upload => ./upload

replace login => ./login
