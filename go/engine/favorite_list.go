// Copyright 2015 Keybase, Inc. All rights reserved. Use of
// this source code is governed by the included BSD license.

package engine

import (
	"github.com/keybase/client/go/libkb"
	keybase1 "github.com/keybase/client/go/protocol"
)

// FavoriteList is an engine.
type FavoriteList struct {
	libkb.Contextified
	result FavoritesAPIResult
}

// NewFavoriteList creates a FavoriteList engine.
func NewFavoriteList(g *libkb.GlobalContext) *FavoriteList {
	return &FavoriteList{
		Contextified: libkb.NewContextified(g),
	}
}

// Name is the unique engine name.
func (e *FavoriteList) Name() string {
	return "FavoriteList"
}

// GetPrereqs returns the engine prereqs.
func (e *FavoriteList) Prereqs() Prereqs {
	return Prereqs{
		Device: true,
	}
}

// RequiredUIs returns the required UIs.
func (e *FavoriteList) RequiredUIs() []libkb.UIKind {
	return []libkb.UIKind{}
}

// SubConsumers returns the other UI consumers for this engine.
func (e *FavoriteList) SubConsumers() []libkb.UIConsumer {
	return nil
}

type FavoritesAPIResult struct {
	Status    libkb.AppStatus   `json:"status"`
	Favorites []keybase1.Folder `json:"favorites"`
	Ignored   []keybase1.Folder `json:"ignored"`
	New       []keybase1.Folder `json:"new"`
}

func (f *FavoritesAPIResult) GetAppStatus() *libkb.AppStatus {
	return &f.Status
}

// Run starts the engine.
func (e *FavoriteList) Run(ctx *Context) error {
	return e.G().API.GetDecode(libkb.APIArg{
		Endpoint:    "kbfs/favorite/list",
		NeedSession: true,
		Args:        libkb.HTTPArgs{},
	}, &e.result)
}

// Favorites returns the list of favorites that Run generated.
func (e *FavoriteList) Result() keybase1.FavoritesResult {
	return keybase1.FavoritesResult{
		Favorites: e.result.Favorites,
		Ignored:   e.result.Ignored,
		// The name "new" is illegal in ObjC. "X" is just a filler character.
		Xnew: e.result.New,
	}
}
