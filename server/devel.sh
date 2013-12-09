export MUSTACHE_REDIS="127.0.0.1:6379"
export MUSTACHE_LISTENON=":8080"
export MUSTACHE_AUTHKEY="im an auth key lololol"
export MUSTACHE_CRYPTKEY="im a crypto key lololollllllllll"
go get ./...
cd bootstrap
go run *.go
cd ..
go clean
go run *.go
