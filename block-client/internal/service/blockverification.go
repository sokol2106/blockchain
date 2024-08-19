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
	block   model.VerificationBlock
	address string
	nonce   string
}

func NewBlockVerification(a, n string) *Blockverification {
	return &Blockverification{
		block:   model.VerificationBlock{},
		address: a,
		nonce:   n,
	}
}

func (b *Blockverification) RequestVerificationData() error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", b.address+"/api/blockchain/block/verify", nil)
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

func (b *Blockverification) VerifyData() error {

	hash := sha256.Sum256([]byte(b.block.Data))
	markleVrf := hex.EncodeToString(hash[:])

	hash = sha256.Sum256([]byte(b.block.Block.Data))
	markleDB := hex.EncodeToString(hash[:])

	if markleDB != markleVrf {
		return b.UpdateStatus(model.StatusFailedAuthenticityCheck, b.block.QueueId)
	}

	markleVrfNonce := fmt.Sprintf(markleVrf+"%s", b.block.Block.Head.Nonce)
	markleDBNonce := fmt.Sprintf(markleDB+"%s", b.block.Block.Head.Nonce)

	hash = sha256.Sum256([]byte(markleVrfNonce))
	markleVrfNonceHash := hex.EncodeToString(hash[:])

	hash = sha256.Sum256([]byte(markleDBNonce))
	markleDBNonceHash := hex.EncodeToString(hash[:])

	if !ValidateNonce(markleVrfNonceHash, b.nonce) || !ValidateNonce(markleDBNonceHash, b.nonce) {
		return b.UpdateStatus(model.StatusFailedAuthenticityCheck, b.block.QueueId)
	}

	return b.UpdateStatus(model.StatusMatched, b.block.QueueId)
}

func (b *Blockverification) UpdateStatus(s model.Status, id string) error {
	data := model.QueueIdStatus{
		Status:  s,
		QueueId: id,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := http.Post(b.address+"/api/blockchain/block/verify",
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
