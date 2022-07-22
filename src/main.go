package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

type Block struct {
	Hash      [32]byte
	prevHash  [32]byte
	Height    int64
	nonce     int
	Timestamp []byte
	PoW       [32]byte
	MR        [32]byte
	Data      [80]byte
	version   int // blockchain(block) version ---> Wallet Address <--- Key pair
	Nonce     int
	bits      int
	Txid      [32]byte
	// MT []*txID
	sig []byte
	//MT Merkle Tree
	txid []byte
}

type Tx struct {
	TxID      [32]byte
	TimeStamp []byte // 블럭 생성 시간
	Applier   []byte // 신청자
	Company   []byte // 경력회사
	Career    []byte // 경력기간
	Payment   []byte // 결제수단
	Job       []byte // 직종, 업무
	Proof     []byte // 경력증명서 pdf
}

type Wallet struct {
	privateKey ecdsa.PrivateKey
	publicKey  []byte
	Address    string
	Alias      string
}

type Wallets struct {
	Wallets map[ecdsa.PrivateKey]*Wallet
}

func main() {

	// structure

	//embedded function

	/*
		func (block *Block) NewBlock() *Block.hash{
			temp:= &block{}
			temp.timestamp = time.Now().UTC(),
			temp.hash = " ",
			temp.nonce = " ",
			temp.prevhash = " ",
			temp.height = 0
			return temp.hash;
		}
	*/
	// 최조의 블럭 이후 2번쨰 블럭을 만들기
	//JSON

	// t *Tx =NewTx(JSON)

	// ts *Txs = []&Tx{}

	// NewBlock(prebHash,Height, t.TxID)

	privateKey, publicKey := newKeyPair()
	w := &Wallet{}
	wallet := w.newWallet(privateKey, publicKey, "test")
	fmt.Println(wallet.Address, "Address")
	fmt.Println(wallet.privateKey, "privateKey")
	fmt.Println(wallet.publicKey, "publicKey")

	/*
		block := NewBlock(GenesisBlock().getBlockID(), GenesisBlock().getHeight(), [32]byte{})
		t := newTx("abc", "def", "ghi", "jkl", "mno", "pqr")
		bs := newBlockChain(&Block{})
		for i := 0; i < 10; i++ {
			block = NewBlock(block.getBlockID(), block.getHeight(), t.TxID)
			bs.addBlock(block)
			block.printBlock()
			fmt.Println(" ")
			fmt.Println(bs.blockchain, "bs.blockchain 입니다!!")
		}
	*/
}

// to make new block we should know prevHase and Tx informations at least.
func NewBlock(prevhash [32]byte, height int64 /*Data string*/, Txid [32]byte) *Block {
	temp := &Block{}
	loc, _ := time.LoadLocation("Asia/Seoul")
	now := time.Now()
	t := now.In(loc)
	temp.Timestamp = []byte(t.String())
	// temp.nonce = int64(1234) + 3
	temp.prevHash = prevhash
	temp.Height = height + 1
	// temp.Hash = Run()
	pow := newProofOfWork(temp)
	temp.nonce, temp.Hash = pow.Run()
	temp.Txid = Txid
	fmt.Printf("this block has just been generated")
	return temp
}
func GenesisBlock() *Block {
	temp := &Block{}
	loc, _ := time.LoadLocation("Asia/Seoul")
	now := time.Now()
	t := now.In(loc)
	temp.Timestamp = []byte(t.String())
	temp.nonce = int(0001)
	temp.prevHash = [32]byte{}
	temp.Height = 0
	temp.Hash = [32]byte{1, 2}
	fmt.Println("this is genesis block")
	return temp
}

// 				------------------- Wallet --------------------
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	privateKey, _ := ecdsa.GenerateKey(curve, rand.Reader)
	publicKey := privateKey.PublicKey
	bPublicKey := append(publicKey.X.Bytes())
	return *privateKey, bPublicKey
}

func (w *Wallet) newWallet(privateKey ecdsa.PrivateKey, publicKey []byte, Alias string) *Wallet {
	w = &Wallet{}
	w.privateKey = privateKey
	w.publicKey = publicKey
	w.Address = base58.Encode(publicKey)
	w.Alias = "Test"
	return w
}

func NewWallets(w *Wallet) *Wallets {
	ws := &Wallets{}
	ws.Wallets = make(map[ecdsa.PrivateKey]*Wallet)
	return ws
}

func HashPublicKey(publicKey [32]byte) []byte {
	publicSHA256 := sha256.Sum256(publicKey[:])
	RIPEMD160Hasher := ripemd160.New()

	_, _ = RIPEMD160Hasher.Write(publicSHA256[:])
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160
}

func (w *Wallet) Sign() {}

//-------------------------------------------------
func (t *Tx) printTx() {
	fmt.Println("==========Transaction Info=============")
	fmt.Printf("TxId: %x\n Applier: %s\n Company: %s\n Career: %s\n Payment: %s\n TimeStamp: %d\n Job: %s\n Proof: %s\n", t.TxID, t.Applier, t.Company, t.Career, t.Payment, t.TimeStamp, t.Job, t.Proof)

}

//func NewTx(){}

// func findTx(TxID [32]byte) *Tx {

// }

func (b *Block) printBlock() {
	fmt.Println("==========블록체인 정보============")
	fmt.Printf("Hash: %x\nHeight: %d\nPrev Hash: %x\nNonce: %d\nPoW: %d\nTimeStamp: %d\nData: %s\nSign: %b\n", b.Hash, b.Height, b.prevHash, b.nonce, b.PoW, b.Timestamp, b.Data, b.sig)
}

func (b *Block) setHash() [32]byte {
	// b.Hash = //SHA-256
	b.Hash = sha256.Sum256(b.Timestamp)
	return b.Hash
}
func (b *Block) getBlockID() [32]byte {
	if b != nil {
		return b.Hash
	} else {
		return [32]byte{}
	}
}
func (b *Block) getHeight() int64 {

	if b == nil {
		return 0
	} else {
		return b.Height
	}
}

func (b *Block) findTx(Txid [32]byte) [32]byte {
	if b.isExisted(Txid) {
		return b.Hash
	} else {
		return [32]byte{}
	}
}

func (b *Block) isExisted(Txid [32]byte) bool {
	return reflect.DeepEqual(Txid, b.Txid)
}

type blocks struct {
	blockchain map[[32]byte]*Block
}

// 블럭을 줄테니 블록체인으로 만들어라
func newBlockChain(b *Block) *blocks {
	bs := &blocks{}
	bs.blockchain = make(map[[32]byte]*Block)
	return bs

}

type Txs struct {
	Txs map[[32]byte]*Tx
}

func (tx *Tx) prepareData() []byte {
	data := bytes.Join([][]byte{
		tx.TimeStamp,
		tx.Payment,
		tx.Applier,
		tx.Company,
		tx.Career,
		tx.Job,
		tx.Proof,
	}, []byte{})
	return data
}

//새로운 트랜잭션 생성
func newTx(applier, company, career, payment, job, proof string) *Tx {
	newTx := &Tx{}
	newTx.Applier = []byte(applier)
	newTx.Company = []byte(company)
	newTx.Career = []byte(career)
	newTx.Payment = []byte(payment)
	newTx.Job = []byte(job)
	newTx.Proof = []byte(proof)
	loc, _ := time.LoadLocation("Asia/Seoul")
	now := time.Now()
	t := now.In(loc)
	newTx.TimeStamp = []byte(t.String())
	data := newTx.prepareData()
	newTx.TxID = sha256.Sum256(data)
	return newTx
}

func newTxs(t *Tx) *Txs {
	txs := &Txs{}
	txs.Txs = make(map[[32]byte]*Tx)
	return txs

}

// Txs(트랜잭션 DB 대용) 생성(최초 한번만 실행)
func createTxDB(tx *Tx) *Txs {
	txs := &Txs{}
	txs.Txs = make(map[[32]byte]*Tx)
	txs.Txs[tx.TxID] = tx
	return txs
}

func (bs *blocks) addBlock(o *Block) error {
	// prevHash를 받아올 방법이 필요함
	if o != nil {
		bs.blockchain[(o.Hash)] = o
	}
	return nil
}

// 트랜잭션 ID를 이용해 Block 조회
func findBlockByTx(txID [32]byte, bs *blocks) *Block {
	// 최신부터 돌려보자
	//최신 블록체인의 높이를 구한다
	current_height := int64(len(bs.blockchain))
	// 최신 블록ID를 찾는다
	curBlockID := [32]byte{}
	for _, v := range bs.blockchain {
		if v.Height == current_height {
			curBlockID = v.Hash
			break
		}
	}
	for {
		blk := bs.blockchain[curBlockID]
		if blk.isExisted(txID) {
			return blk
		} else {
			if reflect.DeepEqual(blk.prevHash, [32]byte{}) {
				return nil
			}
			curBlockID = blk.prevHash
		}
	}
}

// Hash 를 줄테니 블록 주소를 달라
func (bs *blocks) getBlock(blkID []byte) *Block {
	temp := &Block{}
	for _, temp := range bs.blockchain {
		if blkID == nil {
			return nil
		} else if reflect.DeepEqual(temp.Hash, blkID) {
			return temp
		}
	}
	return temp
}

func (bs *blocks) findBlock(height int64) *Block {
	// 최신부터 돌려보자
	//최신 블록체인의 높이를 구한다
	current_height := int64(len(bs.blockchain) - 1)
	if height == 0 {
	}
	// 최신 블록ID를 찾는다
	curBlockID := [32]byte{}
	for _, v := range bs.blockchain {
		if v.Height == current_height {
			curBlockID = v.Hash
		}
		if v.prevHash == [32]byte{} {
			return nil
		}
	}
	//최신 블록을 받아오기
	for {
		blk := bs.blockchain[curBlockID]
		if blk.Height == height {
			return blk
		} else {
			if reflect.DeepEqual(blk.prevHash, [32]byte{}) {
				return nil
			}
			curBlockID = blk.prevHash
		}
	}
	//해당 블록을 검사한다
	//if 있다 --> 종료
	// if 없다 --> next (prevHash)를 찾는다
	block1 := &Block{}
	for _, v := range bs.blockchain {
		if v.Height == height {
			block1 = v
		}
		if v.prevHash == [32]byte{} {
			return nil
		}

	}
	return block1
}

var (
	maxNonce = math.MaxInt64
)

var targetBites = 20

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func IntToHex(n int64) []byte {

	s := fmt.Sprint(n)
	return []byte(s)
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join([][]byte{
		pow.block.prevHash[:],
		pow.block.Data[:],
		//IntToHex(pow.block.Timestamp),
		IntToHex(int64(targetBites)),
		IntToHex(int64(nonce)),
	}, []byte{})
	return data
}
func (pow *ProofOfWork) Run() (int, [32]byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0
	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	return nonce, hash
}
func newProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBites))
	pow := &ProofOfWork{block, target}
	return pow
}

// nonce 만드는 방법
// hash 만드는 방법
// height 와 hash를 바꾸기

//nosql db를

// 블록db
// tx db
// wallet db

// raw data 를 miner가 요청하면 줄 수 있어야함.

// block reading ratency
