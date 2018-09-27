package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"nimona.io/go/codec"
	"nimona.io/go/primitives"
	"nimona.io/go/storage"
)

type blockReq struct {
	Type        string                 `json:"type,omitempty"`
	Annotations map[string]interface{} `json:"annotations,omitempty"`
	Payload     map[string]interface{} `json:"payload,omitempty"`
	Recipient   string                 `json:"recipient"`
}

func (api *API) HandleGetBlocks(c *gin.Context) {
	blockIDs, err := api.blockStore.List()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	ms := []map[string]interface{}{}
	for _, blockID := range blockIDs {
		b, err := api.blockStore.Get(blockID)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		m := &primitives.Block{}
		codec.Unmarshal(b, m)
		ms = append(ms, api.mapBlock(m))
	}
	c.JSON(http.StatusOK, ms)
}

func (api *API) HandleGetBlock(c *gin.Context) {
	blockID := c.Param("blockID")
	b, err := api.blockStore.Get(blockID)
	if err != nil {
		if err == storage.ErrNotFound {
			c.AbortWithError(404, err)
			return
		}
		c.AbortWithError(500, err)
		return
	}
	m := &primitives.Block{}
	codec.Unmarshal(b, m)
	c.JSON(http.StatusOK, api.mapBlock(m))
}

func (api *API) HandlePostBlock(c *gin.Context) {
	req := &blockReq{}
	if err := c.BindJSON(req); err != nil {
		c.AbortWithError(400, err)
		return
	}

	if req.Recipient == "" {
		c.AbortWithError(400, errors.New("missing recipient"))
		return
	}

	keyBlock, err := primitives.BlockFromBase58(req.Recipient)
	if err != nil {
		c.AbortWithError(400, errors.New("invalid recipient key"))
		return
	}
	key := &primitives.Key{}
	key.FromBlock(keyBlock)

	block := &primitives.Block{
		Type:    req.Type,
		Payload: req.Payload,
	}

	ctx := context.Background()
	signer := api.addressBook.GetLocalPeerKey()
	if err := api.exchange.Send(ctx, block, key, primitives.SignWith(signer)); err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
