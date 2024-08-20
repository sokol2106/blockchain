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

type Blockverification struct {
	block model.VerificationBlock
	url   string
	nonce string
}

func NewBlockVerification(u, n string) *Blockverification {
	return &Blockverification{
		block: model.VerificationBlock{},
		url:   u,
		nonce: n,
	}
}

func (b *Blockverification) RequestVerificationData() error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", b.url+"/api/blockchain/block/verify", nil)
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

	data := model.VerificationBlock{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}

	b.block = data
	return nil
}

func (b *Blockverification) VerifyData() {
	b.block.Status = model.StatusFailedAuthenticityCheck
	hash := sha256.Sum256([]byte(b.block.Data))
	markleVrf := hex.EncodeToString(hash[:])

	hash = sha256.Sum256([]byte(b.block.Block.Data))
	markleDB := hex.EncodeToString(hash[:])

	if markleDB != markleVrf {
		return
	}

	markleVrfNonce := fmt.Sprintf(markleVrf+"%s", b.block.Block.Head.Nonce)
	markleDBNonce := fmt.Sprintf(markleDB+"%s", b.block.Block.Head.Nonce)

	hash = sha256.Sum256([]byte(markleVrfNonce))
	markleVrfNonceHash := hex.EncodeToString(hash[:])

	hash = sha256.Sum256([]byte(markleDBNonce))
	markleDBNonceHash := hex.EncodeToString(hash[:])

	if !ValidateNonce(markleVrfNonceHash, b.nonce) || !ValidateNonce(markleDBNonceHash, b.nonce) {
		return
	}

	b.block.Status = model.StatusMatched
}

func (b *Blockverification) UpdateStatus() error {
	data := model.QueueIdStatus{
		Status:  b.block.Status,
		QueueId: b.block.QueueId,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := http.Post(b.url+"/api/blockchain/block/verify",
		"application/json", strings.NewReader(string(body)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error bad status: %s", resp.Status)
	}
	return nil
}

func (b *Blockverification) SetBlock(bv model.VerificationBlock) {
	b.block = bv
}

func (b *Blockverification) GetBlock() model.VerificationBlock {
	return b.block
}
