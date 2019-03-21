package httputils

const (
	HeaderNameAccept        = "accept"
	HeaderNameAuthorization = "authorization"
	HeaderNameContentType   = "content-type"
	HeaderNameCorrelationID = "correlationID"

	// cors headers
	HeaderNameOrigin                        = "Origin"
	HeaderNameAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderNameAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HeaderNameAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderNameAccessControlAllowHeaders     = "Access-Control-Allow-Headers"

	HeaderValueApplicationJSONUTF8 = "application/json; charset=UTF-8"
	HeaderValueApplicationYAML     = "application/x-yaml"
)

var AllowedHeaders = []string{HeaderNameContentType, HeaderNameAuthorization, HeaderNameAccept, HeaderNameCorrelationID}
