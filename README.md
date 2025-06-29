
# gRPC Bidirectional Communication in Go

## Struttura del progetto

- `proto/comm.proto`: definizione del servizio gRPC.
- `server/main.go`: server gRPC che invia eventi periodici e riceve dati dal client.
- `client/main.go`: client che invia messaggi e riceve eventi dal server.

## Requisiti

- Go 1.18+
- `protoc` (Protocol Buffers compiler)
- Plugin `protoc-gen-go` e `protoc-gen-go-grpc`

## Setup

1. **Installazione plugin protoc:**

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

2. **Generazione codice:**

```
protoc --go_out=. --go-grpc_out=. proto/comm.proto
```

3. **Avvio server:**

```
go run server/main.go
```

4. **Avvio client:**

```
go run client/main.go
```

## Note

Il server invia ogni 5 secondi un messaggio al client.
Il client invia 5 messaggi, uno ogni 3 secondi.
