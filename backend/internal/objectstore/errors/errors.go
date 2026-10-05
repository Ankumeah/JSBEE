package errors

import stderrors "errors"

var ErrNotFound = stderrors.New("objectstore: no such object")
var ErrDirectUpload = stderrors.New("objectstore: direct browser uploads are unsupported")
