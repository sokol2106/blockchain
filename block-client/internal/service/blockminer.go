package service

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ivan/blockchain/block-client/internal/model"
	"io"
	"net/http"
	"strings"
)

type Blockminer struct {
	url   string
	nonce string
	data  model.MiningData
	block model.Block
}

func NewBlockMiner(url string, nonce string) *Blockminer {
	return &Blockminer{
		url:   url,
		nonce: nonce,
		block: model.Block{},
	}
}

func (b *Blockminer) RequestMiningData() error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", b.url+"/api/blockchain/data", nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error bad status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	data := model.MiningData{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}

	b.data = data
	return nil
}

func (b *Blockminer) MineData() {
	hash := sha256.Sum256([]byte(b.data.Data))
	markle := hex.EncodeToString(hash[:])

	nonce := 0

	for {
		markleNonce := fmt.Sprintf(markle+"%d", nonce)

		hash = sha256.Sum256([]byte(markleNonce))
		markleNonceHash := hex.EncodeToString(hash[:])

		if ValidateNonce(markleNonceHash, b.nonce) {
			b.block.Head.Hash = markleNonceHash
			b.block.Head.Merkley = markle
			b.block.Head.Nonce = fmt.Sprintf("%d", nonce)
			b.block.Data = b.data.Data
			b.block.Head.Key = b.data.Key
			return
		}
		nonce++
	}
}

func (b *Blockminer) SendMiningBlock() error {
	jsonBlock, err := json.Marshal(b.block)
	resp, err := http.Post(b.url+"/api/blockchain/block", "application/json", strings.NewReader(string(jsonBlock)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("error bad status: %s", resp.Status)
	}
	return nil
}

func (b *Blockminer) SetData(data model.MiningData) {
	b.data = data
}

func (b *Blockminer) GetBlock() model.Block {
	return b.block
}

func ValidateNonce(data string, nonce string) bool {
	if len(data) < len(nonce) {
		return false
	}
	for i := 0; i < len(nonce); i++ {
		if data[i] != nonce[i] {
			return false
		}
	}
	return true
}
