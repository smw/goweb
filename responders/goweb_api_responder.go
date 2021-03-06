package responders

import (
	"github.com/stretchrcom/codecs/constants"
	codecservices "github.com/stretchrcom/codecs/services"
	"github.com/stretchrcom/goweb/context"
	"net/http"
)

const (
	// DefaultStandardFieldDataKey is the default response object field for the data.
	DefaultStandardFieldDataKey string = "d"
	// DefaultStandardFieldStatusKey is the default response object field for the status.
	DefaultStandardFieldStatusKey string = "s"
	// DefaultStandardFieldErrorsKey is the default response object field for the errors.
	DefaultStandardFieldErrorsKey string = "e"
)

type GowebAPIResponder struct {
	httpResponder HTTPResponder
	codecService  codecservices.CodecService

	transformer func(ctx context.Context, object map[string]interface{}) (map[string]interface{}, error)

	// field names

	// StandardFieldDataKey is the response object field name for the data.
	StandardFieldDataKey string

	// StandardFieldStatusKey is the response object field name for the status.
	StandardFieldStatusKey string

	// StandardFieldErrorsKey is the response object field name for the errors.
	StandardFieldErrorsKey string
}

func NewGowebAPIResponder(codecService codecservices.CodecService, httpResponder HTTPResponder) *GowebAPIResponder {
	api := new(GowebAPIResponder)
	api.SetCodecService(codecService)
	api.httpResponder = httpResponder
	api.StandardFieldDataKey = DefaultStandardFieldDataKey
	api.StandardFieldStatusKey = DefaultStandardFieldStatusKey
	api.StandardFieldErrorsKey = DefaultStandardFieldErrorsKey
	return api
}

// SetCodecService sets the codec service to use.
func (a *GowebAPIResponder) SetCodecService(service codecservices.CodecService) {
	a.codecService = service
}

// GetCodecService gets the codec service that will be used by this object.
func (a *GowebAPIResponder) GetCodecService() codecservices.CodecService {

	if a.codecService == nil {
		a.codecService = new(codecservices.WebCodecService)
	}

	return a.codecService
}

// TransformStandardResponseObject transforms the standard response object before it is written to the response if a
// transformer func has been set via SetStandardResponseObjectTransformer.
func (a *GowebAPIResponder) TransformStandardResponseObject(ctx context.Context, object map[string]interface{}) (map[string]interface{}, error) {
	if a.transformer != nil {
		return a.transformer(ctx, object)
	}
	return object, nil
}

// SetStandardResponseObjectTransformer sets the function to use to transform the standard response object beore it is
// written to the response.
func (a *GowebAPIResponder) SetStandardResponseObjectTransformer(transformer func(ctx context.Context, object map[string]interface{}) (map[string]interface{}, error)) {
	a.transformer = transformer
}

// WriteResponseObject writes the status code and response object to the HttpResponseWriter in
// the specified context, in the format best suited based on the request.
//
// Goweb uses the WebCodecService to decide which codec to use when responding
// see http://godoc.org/github.com/stretchrcom/codecs/services#WebCodecService for more information.
//
// This method should be used when the Goweb Standard Response Object does not satisfy the needs of
// the API, but other Respond* methods are recommended.
func (a *GowebAPIResponder) WriteResponseObject(ctx context.Context, status int, responseObject interface{}) error {

	service := a.GetCodecService()

	acceptHeader := ctx.HttpRequest().Header.Get("Accept")
	extension := ctx.FileExtension()
	hasCallback := len(ctx.QueryValue("callback")) > 0

	codec, codecError := service.GetCodecForResponding(acceptHeader, extension, hasCallback)

	if codecError != nil {
		return codecError
	}

	var options map[string]interface{}

	// do we need to add some options?
	if hasCallback {
		options = map[string]interface{}{constants.OptionKeyClientCallback: ctx.QueryValue("callback")}
	}

	output, marshalErr := codec.Marshal(responseObject, options)

	if marshalErr != nil {
		return marshalErr
	}

	// use the HTTP responder to respond
	ctx.HttpResponseWriter().Header().Set("Content-Type", codec.ContentType()) // TODO: test me
	a.httpResponder.With(ctx, status, output)

	return nil

}

// Responds to the Context with the specified status, data and errors.
func (a *GowebAPIResponder) Respond(ctx context.Context, status int, data interface{}, errors []string) error {

	// make the standard response object
	sro := map[string]interface{}{a.StandardFieldStatusKey: status}

	if data != nil {
		sro[a.StandardFieldDataKey] = data
	}
	if len(errors) > 0 {
		sro[a.StandardFieldErrorsKey] = errors
	}

	// transofm the object
	var transformErr error
	sro, transformErr = a.TransformStandardResponseObject(ctx, sro)

	if transformErr != nil {
		return transformErr
	}

	return a.WriteResponseObject(ctx, status, sro)

}

// RespondWithData responds with the specified data, no errors and a 200 StatusOK response.
func (a *GowebAPIResponder) RespondWithData(ctx context.Context, data interface{}) error {
	return a.Respond(ctx, http.StatusOK, data, nil)
}

// RespondWithError responds with the specified error and status code.
func (a *GowebAPIResponder) RespondWithError(ctx context.Context, status int, err string) error {
	return a.Respond(ctx, status, nil, []string{err})
}
