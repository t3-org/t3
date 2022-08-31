package arranger

import (
	"context"

	"github.com/kamva/hexa/hexatranslator"
)

// ApplicationErr tries to convert hexa to application error.
func ApplicationErr(ctx context.Context, err error) error {
	return HexaToApplicationErr(err, hexatranslator.CtxTranslator(ctx))
}
