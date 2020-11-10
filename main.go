/*
 * Parvaeres API
 *
 * Parvaeres magic deployment API
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"

	parvaeres "github.com/parvaeres/parvaeres/pkg/api"
	"github.com/parvaeres/parvaeres/pkg/argocd"
	"github.com/parvaeres/parvaeres/pkg/email"
	"github.com/pkg/errors"
)

func main() {
	var (
		kubeconfig      string
		argoCDNamespace string
		publicURL       string
		mailgunDomain   string
		mailgunAPIKey   string
		mailgunSender   string
		emailEnabled    bool
		adminToken      string
	)

	log.Printf("Server started")
	flag.StringVar(
		&kubeconfig,
		"kubeconfig",
		lookupEnvOrString("PARVAERES_KUBECONFIG", ""),
		"path to Kubernetes config file",
	)
	flag.StringVar(
		&argoCDNamespace,
		"argocd-namespace",
		lookupEnvOrString("PARVAERES_ARGOCD_NAMESPACE", "argocd"),
		"argocd Namespace",
	)
	flag.StringVar(
		&publicURL,
		"public-url",
		lookupEnvOrString("PARVAERES_PUBLIC_URL", "http://poc.parvaeres.io/"),
		"external URL where parvaeres-server is reacheable",
	)
	flag.StringVar(
		&mailgunDomain,
		"mailgun-domain",
		lookupEnvOrString("PARVAERES_MAILGUN_DOMAIN", "poc.parvaeres.io"),
		"mailgun domain",
	)
	flag.StringVar(
		&mailgunAPIKey,
		"mailgun-apikey",
		lookupEnvOrString("PARVAERES_MAILGUN_APIKEY", ""),
		"mailgun API key",
	)
	flag.StringVar(
		&mailgunSender,
		"mailgun-sender",
		lookupEnvOrString("PARVAERES_MAILGUN_SENDER", "Parvaeres Support <support@poc.parvaeres.io>"),
		"mailgun Sender email address",
	)
	flag.BoolVar(
		&emailEnabled,
		"email-enabled",
		lookupEnvOrBool("PARVAERES_EMAIL_ENABLED", false),
		"enable email delivery",
	)
	flag.StringVar(
		&adminToken,
		"admin-token",
		lookupEnvOrString("PARVAERES_ADMIN_TOKEN", "parvaeres"),
		"token to perform admin operations, pass as 'admintoken' header in API requests",
	)
	flag.Parse()

	kubernetesConfig, err := argocd.GetKubernetesConfig(kubeconfig)
	if err != nil {
		log.Fatalf(errors.Wrap(err, "Couldn't get kubernetes config. Bailing out.").Error())
	}
	kubernetesClient, err := argocd.GetKubernetesClientSet(kubernetesConfig)
	if err != nil {
		log.Fatalf(errors.Wrap(err, "Couldn't get kubernetes client. Bailing out.").Error())
	}
	argoCDClient, err := argocd.GetArgoCDClientSet(kubernetesConfig)
	if err != nil {
		log.Fatalf(errors.Wrap(err, "Couldn't get ArgoCD client. Bailing out.").Error())
	}

	DefaultApiService := &parvaeres.DefaultApiService{
		EmailProvider: &email.MailGun{
			Domain: mailgunDomain,
			APIKey: mailgunAPIKey,
			Sender: mailgunSender,
		},
		GitopsProvider: &argocd.ArgoCD{
			ArgoCDNamespace:  argoCDNamespace,
			ArgoCDclient:     argoCDClient,
			KubernetesClient: kubernetesClient,
		},
		FeatureFlagEmail: emailEnabled,
		PublicURL:        publicURL,
		AdminToken:       adminToken,
	}
	DefaultApiController := parvaeres.NewDefaultApiController(DefaultApiService)

	router := parvaeres.NewRouter(DefaultApiController)

	log.Fatal(http.ListenAndServe(":8080", router))
}

// FIXME: this could be improved, maybe cobra
// See: https://www.gmarik.info/blog/2019/12-factor-golang-flag-package/
func lookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func lookupEnvOrInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok {
		v, err := strconv.Atoi(val)
		if err != nil {
			log.Fatalf("LookupEnvOrInt[%s]: %v", key, err)
		}
		return v
	}
	return defaultVal
}

func lookupEnvOrBool(key string, defaultVal bool) bool {
	if val, ok := os.LookupEnv(key); ok {
		v, err := strconv.ParseBool(val)
		if err != nil {
			log.Fatalf("LookupEnvOrBool[%s]: %v", key, err)
		}
		return v
	}
	return defaultVal
}
