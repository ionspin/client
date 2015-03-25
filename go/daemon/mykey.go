package main

import (
	"github.com/keybase/client/go/engine"
	"github.com/keybase/client/go/libkb"
	keybase_1 "github.com/keybase/client/protocol/go"
	"github.com/maxtaco/go-framed-msgpack-rpc/rpc2"
)

type MykeyHandler struct {
	BaseHandler
}

func NewMykeyHandler(xp *rpc2.Transport) *MykeyHandler {
	return &MykeyHandler{BaseHandler{xp: xp}}
}

func (h *MykeyHandler) PgpKeyGen(arg keybase_1.PgpKeyGenArg) (err error) {
	earg := engine.ImportPGPKeyImportEngineArg(arg)
	return h.keygen(earg, true)
}

func (h *MykeyHandler) keygen(earg engine.PGPKeyImportEngineArg, doInteractive bool) (err error) {
	sessionId := nextSessionId()
	ctx := &engine.Context{LogUI: h.getLogUI(sessionId), SecretUI: h.getSecretUI(sessionId)}
	earg.Gen.AddDefaultUid()
	eng := engine.NewPGPKeyImportEngine(earg)
	err = engine.RunEngine(eng, ctx)
	return err
}

func (h *MykeyHandler) PgpKeyGenDefault(arg keybase_1.PgpCreateUids) (err error) {
	earg := engine.PGPKeyImportEngineArg{
		Gen: &libkb.PGPGenArg{
			Ids:         libkb.ImportPgpIdentities(arg.Ids),
			NoDefPGPUid: !arg.UseDefault,
		},
	}
	return h.keygen(earg, false)
}

func (h *MykeyHandler) DeletePrimary() (err error) {
	return libkb.DeletePrimary()
}

func (h *MykeyHandler) Select(sarg keybase_1.SelectArg) error {
	sessionID := nextSessionId()
	gpgui := NewRemoteGPGUI(sessionID, h.getRpcClient())
	secretui := h.getSecretUI(sessionID)
	arg := engine.GPGImportKeyArg{Query: sarg.Query, AllowMulti: sarg.AllowMulti, SkipImport: sarg.SkipImport}
	gpg := engine.NewGPGImportKeyEngine(&arg)
	ctx := &engine.Context{
		GPGUI:    gpgui,
		SecretUI: secretui,
		LogUI:    h.getLogUI(sessionID),
		LoginUI:  h.getLoginUI(sessionID),
	}
	return engine.RunEngine(gpg, ctx)
}

func (h *MykeyHandler) Update(arg keybase_1.UpdateArg) error {
	ctx := engine.Context{
		LogUI: h.getLogUI(arg.SessionID),
	}
	eng := engine.NewPGPUpdateEngine(arg.Fingerprints, arg.All)
	return engine.RunEngine(eng, &ctx)
}
