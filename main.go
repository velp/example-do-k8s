package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type TokenSource struct {
	AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

func main() {
	DOToken := os.Getenv("DO_TOKEN")
	ClusterID := os.Getenv("DO_CLUSTER_ID")
	Namespace := os.Getenv("K8S_NAMESPACE")
	// create DO client
	tokenSource := &TokenSource{
		AccessToken: DOToken,
	}
	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	doClient := godo.NewClient(oauthClient)
	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*15)
	defer cancel()
	// Get credentials
	k8sCredsRequest := &godo.KubernetesClusterCredentialsGetRequest{}
	creds, _, err := doClient.Kubernetes.GetCredentials(ctx, ClusterID, k8sCredsRequest)
	if err != nil {
		log.Fatalf("DO creds failed: %s", err)
	}
	log.Printf("Server: %s", creds.Server)
	log.Printf("CA: %s", creds.CertificateAuthorityData)
	log.Printf("Token: %s", creds.Token)
	// Init k8s client
	clientConfig, err := clientcmd.BuildConfigFromFlags(creds.Server, "")
	if err != nil {
		log.Fatalf("k8s config from flags failed: %s", err)
	}
	clientConfig.CAData = creds.CertificateAuthorityData
	clientConfig.BearerToken = creds.Token
	k8sClient, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		log.Fatalf("new k8s client failed: %s", err)
	}
	ns, err := k8sClient.CoreV1().Namespaces().Get(Namespace, metaV1.GetOptions{})
	if err != nil {
		log.Fatalf("getting k8s namespace failed: %s", err)
	}
	log.Printf("Namespace: %+v", ns)
}
