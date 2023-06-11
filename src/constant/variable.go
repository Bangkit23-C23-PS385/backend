package constant

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type TemplateType string
type TemplateSubject string

const (
	Register      TemplateType = "register"
	VerifySuccess TemplateType = "verify success"

	VerifyEmail         TemplateSubject = "Verifikasi Email Anda - Medicare"
	RegistrationSuccess TemplateSubject = "Verifikasi Email Berhasil! - Medicare"
)

var (
	_ = godotenv.Load()

	ApplicationDomain = os.Getenv("APPLICATION_DOMAIN")
	ApplicationName   = os.Getenv("APPLICATION_NAME")
	APIOriginURL      = os.Getenv("API_ORIGIN_URL")
	RedirectURL       = os.Getenv("REDIRECT_URL")

	WasabiBucketRegion   = os.Getenv("WASABI_BUCKET_REGION")
	WasabiBucketEndpoint = os.Getenv("WASABI_BUCKET_ENDPOINT")
	WasabiBucketName     = os.Getenv("WASABI_BUCKET_NAME")
	WasabiAccessKey      = os.Getenv("WASABI_ACCESS_KEY")
	WasabiSecretKey      = os.Getenv("WASABI_SECRET_KEY")

	GCPProjectID       = os.Getenv("GCP_PROJECT_ID")
	GCPTopicPredict    = os.Getenv("GCP_PUBSUB_TOPIC_PREDICT")
	GCPTopicSubmitData = os.Getenv("GCP_PUBSUB_TOPIC_SUBMIT_DATA")

	AccessTokenInterval, _ = strconv.Atoi(os.Getenv("ACCESS_TOKEN_INTERVAL"))
	AccessTokenDuration    = time.Duration(AccessTokenInterval) * time.Second

	MaxEmailLength    = 320
	MinUsernameLength = 5
	MaxUsernameLength = 80
	MinPasswordLength = 8
	MaxPasswordLength = 15
	EmailRegex        = `^[a-zA-Z0-9]+(([\-\._][a-zA-Z0-9]+)?)+\@[a-zA-Z0-9]+([\-][a-zA-Z0-9]+)*(\.[a-zA-Z0-9]{2,})+$`
	UsernameRegex     = `^[a-z0-9\_\-\.]+$`

	PageDefault  = 1
	LimitDefault = 10
	MinLimit     = 1
	MinPage      = 1
	MaxLimit     = 20

	TimeLocation = "Asia/Jakarta"
)
