package commands

import (
	"context"
	"testing"
	"time"

	"github.com/keybase/client/go/chat/globals"
	"github.com/keybase/client/go/chat/types"
	"github.com/keybase/client/go/externalstest"
	"github.com/keybase/client/go/kbtest"
	"github.com/keybase/client/go/protocol/chat1"
	"github.com/keybase/client/go/protocol/gregor1"
	"github.com/stretchr/testify/require"
)

func TestFlipPreview(t *testing.T) {
	tc := externalstest.SetupTest(t, "flip", 0)
	ui := kbtest.NewChatUI()
	g := globals.NewContext(tc.G, &globals.ChatContext{})
	g.CoinFlipManager = types.DummyCoinFlipManager{}
	g.UIRouter = &fakeUIRouter{ui: ui}
	flip := NewFlip(g)
	ctx := context.TODO()
	uid := gregor1.UID{1, 2, 3, 4}
	convID := chat1.ConversationID{1, 2, 3, 4}
	timeout := 2 * time.Second

	flip.Preview(ctx, uid, convID, "/flip")
	select {
	case s := <-ui.CommandMarkdown:
		require.NotZero(t, len(s))
	case <-time.After(timeout):
		require.Fail(t, "no text")
	}

	flip.Preview(ctx, uid, convID, "/flip 6")
	select {
	case s := <-ui.CommandMarkdown:
		require.NotZero(t, len(s))
	case <-time.After(timeout):
		require.Fail(t, "no text")
	}

	flip.Preview(ctx, uid, convID, "/flip ")
	select {
	case s := <-ui.CommandMarkdown:
		require.NotZero(t, len(s))
	case <-time.After(timeout):
		require.Fail(t, "no text")
	}

	flip.Preview(ctx, uid, convID, "/fli")
	select {
	case s := <-ui.CommandMarkdown:
		require.Zero(t, len(s))
	case <-time.After(timeout):
		require.Fail(t, "no text")
	}

	flip.Preview(ctx, uid, convID, "/fli")
	select {
	case <-ui.CommandMarkdown:
		require.Fail(t, "no text expected")
	default:
	}

}
