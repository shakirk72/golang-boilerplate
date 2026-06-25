package config

import (
	"context"
)

// LoadSecrets injects secret values as environment variables before Load is called.
// In dev/qa this is a no-op — values come from .env via godotenv.
//
// In prod, replace the stub below with AWS SecretManager calls, e.g.:
//
//	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion("ap-southeast-1"))
//	svc := secretsmanager.NewFromConfig(awsCfg)
//	out, err := svc.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
//	    SecretId: aws.String("go-mvc-app/prod"),
//	})
//	var kv map[string]string
//	json.Unmarshal([]byte(*out.SecretString), &kv)
//	for k, v := range kv { os.Setenv(k, v) }
//
// All injected keys become visible to the subsequent config.Load() call.
func LoadSecrets(ctx context.Context, env string) error {
	if env != EnvProd {
		return nil
	}
	// prod: fetch from AWS SecretManager and inject via os.Setenv
	_ = ctx
	return nil
}
