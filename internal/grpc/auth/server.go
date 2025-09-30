package auth

import ssov1 "github.com/Kosk0l/Protos/gen/go/sso"

type serverAPI struct {
	ssov1.UnimplementedAuthServer
}