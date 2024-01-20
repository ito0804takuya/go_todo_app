package auth

import (
	_ "embed"
)

// NOTE: go:embedによって実行バイナリにpemファイルを埋め込むことができる。そうすることでシングルバイナリで実行可能になる

//go:embed cert/secret.pem
var rawPrivKey []byte

//go:embed cert/public.pem
var rawPubKey []byte
