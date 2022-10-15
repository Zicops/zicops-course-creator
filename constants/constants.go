package constants

const (
	// ENV_DEBUG ..... for debugging serverity
	ENV_DEBUG             string = "DEBUG"
	MAX_READ_TIMEOUT             = 50 // secs
	MAX_WRITE_TIMEOUT            = 50 // secs
	RESPONSE_JSON_DATA    string = "data"
	RESPONSDE_JSON_ERROR  string = "error"
	COURSES_BUCKET               = "courses-zicops-one"
	COURSES_PUBLIC_BUCKET        = "courses-public-zicops-one"
	ZICOPS_ASSETS_BUCKET         = "zicops-assets"
)

var (
	StaticTypeMap = map[string]string{
		"scorm":  "story.html",
		"cmi5":   "story.html",
		"tincan": "story.html",
		"html5":  "story.html",
	}
)
