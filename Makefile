# This makefile is needed to build all the executables and run the server on single machine.
# More advanced script needed to wrap all the executables in docker and send them to remote machines, for example.
run_backend: all
	

all: web_interface user_service message_service file_storage_service notification_service tools


# Building web interface
web_interface: web_interface_dir
	cd web_interface\cmd && go build -o ../../build/web_interface/cmd/
	copy web_interface\config\config.json build\web_interface\config\config.json

web_interface_dir: build_dir
	-echo y | rmdir /s build\web_interface
	mkdir build\web_interface
	mkdir build\web_interface\cmd build\web_interface\config

# Building user service
user_service: uservice_interface uservice_worker

uservice_interface: uservice_interface_dir
	cd user_service\interface\cmd && go build -o ../../../build/user_service/interface/cmd/
	copy user_service\interface\config\config.json build\user_service\interface\config\config.json

uservice_worker: uservice_worker_dir
	cd user_service\worker\cmd && go build -o ../../../build/user_service/worker/cmd/
	copy user_service\worker\config\config.json build\user_service\worker\config\config.json

uservice_interface_dir: user_service_dir
	-echo y | rmdir /s build\user_service\interface
	mkdir build\user_service\interface
	mkdir build\user_service\interface\config build\user_service\interface\cmd

uservice_worker_dir: user_service_dir
	-echo y | rmdir /s build\user_service\worker
	mkdir build\user_service\worker
	mkdir build\user_service\worker\config build\user_service\worker\cmd

user_service_dir: build_dir
	-echo y | rmdir /s build\user_service
	mkdir build\user_service

# Building message service
message_service: mservice_interface mservice_worker

mservice_interface: mservice_interface_dir
	cd message_service\interface\cmd && go build -o ../../../build/message_service/interface/cmd/
	copy message_service\interface\config\config.json build\message_service\interface\config\config.json

mservice_worker: mservice_worker_dir
	cd message_service\worker\cmd && go build -o ../../../build/message_service/worker/cmd/
	copy message_service\worker\config\config.json build\message_service\worker\config\config.json

mservice_interface_dir: message_service_dir
	-echo y | rmdir /s build\message_service\interface
	mkdir build\message_service\interface
	mkdir build\message_service\interface\config build\message_service\interface\cmd

mservice_worker_dir: message_service_dir
	-echo y | rmdir /s build\message_service\worker
	mkdir build\message_service\worker
	mkdir build\message_service\worker\config build\message_service\worker\cmd

message_service_dir: build_dir
	-echo y | rmdir /s build\message_service
	mkdir build\message_service

# Building file storage service
file_storage_service: fservice_interface fservice_worker

fservice_interface: fservice_interface_dir
	cd file_storage_service\interface\cmd && go build -o ../../../build/file_storage_service/interface/cmd/
	copy file_storage_service\interface\config\config.json build\file_storage_service\interface\config\config.json

fservice_worker: fservice_worker_dir
	cd file_storage_service\worker\cmd && go build -o ../../../build/file_storage_service/worker/cmd/
	copy file_storage_service\worker\config\config.json build\file_storage_service\worker\config\config.json

fservice_interface_dir: file_storage_service_dir
	-echo y | rmdir /s build\file_storage_service\interface
	mkdir build\file_storage_service\interface
	mkdir build\file_storage_service\interface\config build\file_storage_service\interface\cmd

fservice_worker_dir: file_storage_service_dir
	-echo y | rmdir /s build\file_storage_service\worker
	mkdir build\file_storage_service\worker
	mkdir build\file_storage_service\worker\config build\file_storage_service\worker\cmd

file_storage_service_dir: build_dir
	-echo y | rmdir /s build\file_storage_service
	mkdir build\file_storage_service

# Building notification service
notificaton_service: notification_service_dir
	cd notification_service\cmd && go build -o ../../../build/notification_service/cmd/
	copy notification_service\config\config.json build\notification_service\config\config.json

notification_service_dir: build_dir
	-echo y | rmdir /s build\notification_service
	mkdir build\notification_service
	mkdir build\notification_service\cmd build\notification_service\config

tools: redis_filler test_redis

redis_filler: tools_dir
	-echo y | rmdir /s build\tools\redis_filler
	mkdir build\tools\redis_filler
	cd tools\redis_filler && go build ../../build/tools/redis_filler/

test_redis: tools_dir
	-echo y | rmdir /s build\tools\test_redis
	mkdir build\tools\test_redis
	cd tools\test_redis && go build ../../build/tools/test_redis/

tools_dir: build_dir
	-echo y | rmdir /s build\tools
	mkdir build\tools
	

build_dir:
	-echo y | rmdir /s build
	mkdir build
