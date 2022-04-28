module github.com/Ghytro/go_messenger/message_service/interface

replace github.com/Ghytro/go_messenger/lib => ../../lib

replace github.com/Ghytro/go_messenger/message_service/interface => ./

go 1.17

require (
	github.com/Ghytro/go_messenger/lib v0.0.0-20220322205922-225ffdb5ab12 // indirect
	github.com/lib/pq v1.10.5 // indirect
)
