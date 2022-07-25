package main

// gRPC server

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	_ "errors"
	"fmt"
	"net"
	"net/rpc"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

type Return string // 월렛 주소 반환하는 값을 담은 바구니

type Response struct {
	C       int
	Address string
}

type Wallet struct {
	PublicKey  []byte
	PrivateKey ecdsa.PrivateKey
	Address    string
	Alias      string
}
type Wallets struct {
	Wts map[string]*Wallet
}

// ---------------------------------------------------------------- Functions --------------------------------
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	privateKey, _ := ecdsa.GenerateKey(curve, rand.Reader)
	publicKey := privateKey.PublicKey
	bPublicKey := append(publicKey.X.Bytes())
	return *privateKey, bPublicKey
}
func HashPublicKey(publicKey []byte) []byte {
	publicSHA256 := sha256.Sum256(publicKey)
	RIPEMD160Hasher := ripemd160.New()
	_, _ = RIPEMD160Hasher.Write(publicSHA256[:])

	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)
	return publicRIPEMD160
}
func encodeAddress(publicREPEMD160 []byte) string {
	version := byte(0x00)
	s := base58.CheckEncode(publicREPEMD160, version)
	return s
}
func (w *Wallet) MakeWallet(Alias string) string {
	ws := new(Wallets)
	w = new(Wallet)
	PrivateKey, PublicKey := newKeyPair()
	w.PrivateKey = PrivateKey
	w.PublicKey = PublicKey
	// 유효성검사 Address 가 존재하지 않는다면
	if ws.PutWallet(encodeAddress(HashPublicKey(PublicKey))) {
		w.Address = encodeAddress(HashPublicKey(PublicKey))
	} else {
		// 만약 이미 존재하는 Address를 만들었다면 , 다시 키쌍을 만들기
		PrivateKey, PublicKey = newKeyPair()
		w.PrivateKey = PrivateKey
		w.PublicKey = PublicKey
		w.Address = encodeAddress(HashPublicKey(PublicKey))
	}
	w.Alias = Alias
	return w.Address
}
func NewWallets(w *Wallet) *Wallets {
	Ws := &Wallets{}
	Ws.Wts = make(map[string]*Wallet)
	return Ws
}
func (ws *Wallets) PutWallet(encodedAddress string) bool {
	if ws.Wts[encodedAddress] != nil {
		return false
	}
	return true
}

// -- 실제 Response 해주는 Functions
func (r *Return) SendWalletAddress(wallet Wallet, response *Response) error {
	response.Address = wallet.MakeWallet(wallet.Alias)
	return nil
}

// -------------------- main ----------------------------------------------------

func main() {
	rpc.Register(new(Return))
	In, err := net.Listen("tcp", ":9000")
	fmt.Println(In, "In 입니다")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer In.Close()
	for {
		conn, err := In.Accept()
		fmt.Println(conn, err, "In.Accept 입니다")
		if err != nil {
			continue
		}
		defer conn.Close()

		go rpc.ServeConn(conn)
	}
}
