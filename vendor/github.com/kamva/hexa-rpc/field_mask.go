package hrpc

import (
	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	"google.golang.org/genproto/protobuf/field_mask"
)

// MaskMapPaths mask all paths in the provided map with the provided depth.
func MaskMapPaths(m hexa.Map, mask *field_mask.FieldMask, depth int) {
	mask.Paths = gutil.MapPathExtractor{Depth: depth, Separator: "."}.Extract(m)
}
