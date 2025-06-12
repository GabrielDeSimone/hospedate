package models

import (
    "crypto/ecdsa"
    "encoding/hex"

    "github.com/ethereum/go-ethereum/common/hexutil"
    "github.com/ethereum/go-ethereum/crypto"

    "crypto/sha256"

    "github.com/mr-tron/base58"
)

type PrivateKey struct {
    ecdsaPk ecdsa.PrivateKey
}
type PublicKey struct {
    ecdsaPublicKey ecdsa.PublicKey
}
type Address struct {
    Address string
}

func NewPrivateKey() (*PrivateKey, error) {
    privateKey, err := crypto.GenerateKey()
    if err != nil {
        return nil, err
    }
    return &PrivateKey{*privateKey}, nil
}

func (pk *PrivateKey) String() string {
    privateKeyBytes := crypto.FromECDSA(&pk.ecdsaPk)
    return hexutil.Encode(privateKeyBytes)
}

func (pk *PrivateKey) GetPublicKey() (*PublicKey, error) {
    publicKey := (&pk.ecdsaPk).Public()
    publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
    if !ok {
        return nil, ErrCastPublicKey
    }
    return &PublicKey{*publicKeyECDSA}, nil
}

func (pub *PublicKey) GetAddress() (*Address, error) {
    address := crypto.PubkeyToAddress(pub.ecdsaPublicKey).Hex()
    address = "41" + address[2:]
    addressBytes, err := hex.DecodeString(address)
    if err != nil {
        return nil, err
    }
    hash1 := s256(s256(addressBytes))
    secret := hash1[:4]
    for _, v := range secret {
        addressBytes = append(addressBytes, v)
    }
    addressBase58 := base58.Encode(addressBytes)
    return &Address{addressBase58}, nil

}

func (pub *PublicKey) String() string {
    publicKeyBytes := crypto.FromECDSAPub(&pub.ecdsaPublicKey)
    return hexutil.Encode(publicKeyBytes)
}

func s256(hash []byte) []byte {
    h := sha256.New()
    h.Write(hash)
    return h.Sum(nil)
}
