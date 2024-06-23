package config

import (
	"os"

	"github.com/joho/godotenv"
)

var (
	JWT_SECRET          = ""
	PORT                = ""
	RSAPUBKEY           = []byte("")
	err                 error
	DB_HOST             = "localhost"
	DB_USER             = "postgres"
	DB_PASS             = "pgadmin123"
	DB_NAME             = "billing"
	DB_PORT             = "5432"
	DB_SSL              = "disable"
	IAM_HOST            = "https://api-iam.34d.me"
	IAM_BEARER          = "secret"
	MIDTRANS_SERVER_KEY = "SB-Mid-server-4Q9Q6Q1Q9Q1Q"
)

func init() {
	godotenv.Load()

	JWT_SECRET = os.Getenv("JWT_SECRET")
	if JWT_SECRET == "" {
		panic("JWT_SECRET is not set")
	}

	PORT = os.Getenv("PORT")
	if PORT == "" {
		panic("PORT is not set")
	}

	// RSAPUBKEY, err = os.ReadFile("saasiam.pem")
	// if err != nil {
	// 	panic(err)
	// }

	// if RSAPUBKEY == nil {
	// 	panic("RSA public key is not set")
	// }

	if os.Getenv("DB_HOST") != "" {
		DB_HOST = os.Getenv("DB_HOST")
	}

	if os.Getenv("DB_USER") != "" {
		DB_USER = os.Getenv("DB_USER")
	}

	if os.Getenv("DB_PASS") != "" {
		DB_PASS = os.Getenv("DB_PASS")
	}

	if os.Getenv("DB_NAME") != "" {
		DB_NAME = os.Getenv("DB_NAME")
	}

	if os.Getenv("DB_PORT") != "" {
		DB_PORT = os.Getenv("DB_PORT")
	}

	if os.Getenv("DB_SSL") != "" {
		DB_SSL = os.Getenv("DB_SSL")
	}

	if os.Getenv("IAM_HOST") != "" {
		IAM_HOST = os.Getenv("IAM_HOST")
	}

	if os.Getenv("IAM_BEARER") != "" {
		IAM_BEARER = os.Getenv("IAM_BEARER")
	}

	if os.Getenv("MIDTRANS_SERVER_KEY") != "" {
		MIDTRANS_SERVER_KEY = os.Getenv("MIDTRANS_SERVER_KEY")
	}
}
