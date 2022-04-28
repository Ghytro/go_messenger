module github.com/Ghytro/go_messenger/message_service/worker

replace github.com/Ghytro/go_messenger/message_service/worker/config => ./

replace github.com/Ghytro/go_messenger/lib => ../../lib

go 1.17

require (
	github.com/Ghytro/go_messenger/lib v0.0.0-20220322205922-225ffdb5ab12 // indirect
	github.com/go-redis/redis v6.15.9+incompatible // indirect
	github.com/lib/pq v1.10.4 // indirect
)
