module github.com/Ghytro/go_messenger/user_service/worker

replace github.com/Ghytro/go_messenger/user_service/worker/config => ./

replace github.com/Ghytro/go_messenger/lib => ../../lib

go 1.17

require (
	github.com/Ghytro/go_messenger/lib v0.0.0-20220304115416-aa74160b67b3 // indirect
	github.com/go-redis/redis v6.15.9+incompatible // indirect
	github.com/lib/pq v1.10.4 // indirect
)
