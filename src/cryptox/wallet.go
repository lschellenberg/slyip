package cryptox

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"log"
	"os"
	"yip/src/config"
)

type Wallet struct {
	Private *ecdsa.PrivateKey
	Public  common.Address
}

func InitWallet(c *config.Config) (*keystore.Key, error) {
	return nil, nil
}

func (w *Wallet) FromPrivateKey(privateKeyHex string) error {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return err
	}
	w.Private = privateKey

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("unable to create public key from private key")
	}

	w.Public = crypto.PubkeyToAddress(*publicKeyECDSA)

	return nil
}

func WalletFromPrivateKey(privateKey string) (*keystore.Key, error) {
	wallet, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, err
	}
	publicKey := wallet.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting public Private to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	return &keystore.Key{
		Id:         uuid.New(),
		Address:    address,
		PrivateKey: wallet,
	}, nil
}

func FromPrivateKey(privateKeyHex string) (Wallet, error) {
	tk := Wallet{}
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return tk, err
	}
	tk.Private = privateKey

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return tk, fmt.Errorf("unable to create public key from private key")
	}

	tk.Public = crypto.PubkeyToAddress(*publicKeyECDSA)

	return tk, nil
}

func KeyFromWalletAndPasswordFile(keystoreUTCPath string, password string) (*keystore.Key, error) {
	keyJSON, err := os.ReadFile(keystoreUTCPath)
	if err != nil {
		return nil, fmt.Errorf("couldn't read wallet file %v %v", keystoreUTCPath, err.Error())
	}

	key, err := keystore.DecryptKey(keyJSON, password)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func ReadPasswordFile(passwordFile string) (string, error) {
	password, err := os.ReadFile(passwordFile)
	if err != nil {
		return "", fmt.Errorf("couldn't read wallet file %v %v", passwordFile, err.Error())
	}

	return string(password), nil
}

func PrivateKeyFromKey(key *keystore.Key) string {
	privateKeyBytes := crypto.FromECDSA(key.PrivateKey)
	return hexutil.Encode(privateKeyBytes)
}

func PrivateKeyFromECDSA(key *ecdsa.PrivateKey) string {
	privateKeyBytes := crypto.FromECDSA(key)
	return hexutil.Encode(privateKeyBytes)
}

func PublicKeyFromKey(k *keystore.Key) string {
	return k.Address.String()
}

func PublicKeyFromPrivateKey(privateKey *ecdsa.PrivateKey) string {

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)

	if !ok {
		return ""
	}

	return crypto.PubkeyToAddress(*publicKeyECDSA).String()
}

func NewKeyFromECDSA(privateKeyECDSA *ecdsa.PrivateKey) *keystore.Key {
	id, err := uuid.NewRandom()
	if err != nil {
		panic(fmt.Sprintf("Could not create random uuid: %v", err))
	}
	key := &keystore.Key{
		Id:         id,
		Address:    crypto.PubkeyToAddress(privateKeyECDSA.PublicKey),
		PrivateKey: privateKeyECDSA,
	}
	return key
}

func GenerateNewKey() (*keystore.Key, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	return NewKeyFromECDSA(privateKey), nil
}

func test() {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	fmt.Println(hexutil.Encode(privateKeyBytes)[2:]) // 0xfad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println(hexutil.Encode(publicKeyBytes)[4:]) // 0x049a7df67f79246283fdc93af76d4f8cdd62c4886e8cd870944e817dd0b97934fdd7719d0810951e03418205868a5c1b40b192451367f28e0088dd75e15de40c05

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println(address) // 0x96216849c49358B10257cb55b28eA603c874b05E

	//hash := sha3.NewKeccak256()
	//hash.Write(publicKeyBytes[1:])
	//fmt.Println(hexutil.Encode(hash.Sum(nil)[12:])) // 0x96216849c49358b10257cb55b28ea603c874b05e
}
