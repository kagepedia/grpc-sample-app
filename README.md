# grpc-sample-app

## Description
GoでgRPC通信のサンプルアプリケーションを作成
以下の4つの通信について実装・検証を行う

- Unary RPC
- Server Streaming RPC
- Client Streaming RPC
- Bidirectional Streaming RPC

### コマンド
```shell
protoc --go_out=./ --go_opt=paths=source_relative \
	--go-grpc_out=./ --go-grpc_opt=paths=source_relative \
	hoge.proto
```