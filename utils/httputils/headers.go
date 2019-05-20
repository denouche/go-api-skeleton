package httputils

const (
	HeaderNameAccept        = "accept"
	HeaderNameAuthorization = "authorization"
	HeaderNameCacheControl  = "cache-control"
	HeaderNameContentType   = "content-type"
	HeaderNameCorrelationID = "correlationID"
	HeaderNameETag          = "ETag"
	HeaderNameExpires       = "expires"
	HeaderNameIfMatch       = "If-Match"
	HeaderNameLocation      = "location"

	// cors headers
	HeaderNameOrigin                        = "Origin"
	HeaderNameAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderNameAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HeaderNameAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderNameAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HeaderNameAccessControlExposeHeaders    = "access-control-expose-headers"

	HeaderValueApplicationJSONUTF8 = "application/json; charset=UTF-8"
	HeaderValueApplicationYAML     = "application/x-yaml"
)

var AllowedHeaders = []string{
	HeaderNameAuthorization,
	HeaderNameAccept,
	HeaderNameCacheControl,
	HeaderNameContentType,
	HeaderNameCorrelationID,
	HeaderNameExpires,
	HeaderNameIfMatch,
}
