package app

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/kamva/hexa/hexatranslator"
	vcomplement "github.com/kamva/ozzo-vcomplement"
	"github.com/kamva/tracer"
)

func v(ctx context.Context, in validation.Validatable) error {
	return tracer.Trace(vcomplement.TranslateByHexa(hexatranslator.CtxTranslator(ctx), in.Validate()))
}
