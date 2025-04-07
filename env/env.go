package env

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Env struct {
	SW6_ADMIN_CLIENT_ID          string `env:"SW6_ADMIN_CLIENT_ID,required"`
	SW6_ADMIN_CLIENT_SECRET      string `env:"SW6_ADMIN_CLIENT_SECRET,required"`
	BASE_URL                     string `env:"BASE_URL,required"`
	KOSATEC_URL                  string `env:"KOSATEC_URL,required"`
	WORTMANN_FTP_SERVER          string `env:"WORTMANN_FTP_SERVER,required"`
	WORTMANN_FTP_SERVER_USER     string `env:"WORTMANN_FTP_SERVER_USER,required"`
	WORTMANN_FTP_SERVER_PASSWORD string `env:"WORTMANN_FTP_SERVER_PASSWORD,required"`
	SALES_CHANNEL_ID             string `env:"SALES_CHANNEL_ID,required"`
	MAIN_CATEGORY_ID             string `env:"MAIN_CATEGORY_ID,required"`
	TAX_ID                       string `env:"TAX_ID,required"`
	CURRENCY_ID                  string `env:"CURRENCY_ID,required"`
	MEDIA_FOLDER_ID              string `env:"MEDIA_FOLDER_ID,required"`
	LIEFERZEIT_ID                string `env:"LIEFERZEIT_ID,required"`
	MY_NAMESPACE                 string `env:"MY_NAMESPACE,required"`
	SSH_HOST                     string `env:"SSH_HOST,required"`
	SSH_USER                     string `env:"SSH_USER,required"`
	SSH_PASS                     string `env:"SSH_PASS,required"`
	FTP_PATH                     string `env:"FTP_PATH,required"`
	FTP_HOST                     string `env:"FTP_HOST,required"`
	FTP_USER                     string `env:"FTP_USER,required"`
	FTP_PASSWORD                 string `env:"FTP_PASSWORD,required"`
	FTP_SHOP_PATH                string `env:"FTP_SHOP_PATH,required"`
	LOG_MAIL                     string `env:"LOG_MAIL,required"`
	MAIL_SERVER                  string `env:"MAIL_SERVER,required"`
	MAIL_PORT                    int    `env:"MAIL_PORT,required"`
	MAIL_USER                    string `env:"MAIL_USER,required"`
	MAIL_PASSWORD                string `env:"MAIL_PASSWORD,required"`
}

func Get() (*Env, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	cfg := Env{}

	err = env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
