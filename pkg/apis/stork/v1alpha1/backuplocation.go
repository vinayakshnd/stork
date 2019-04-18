package v1alpha1

import (
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	// BackupLocationResourceName is name for "backuplocation" resource
	BackupLocationResourceName = "backuplocation"
	// BackupLocationResourcePlural is plural for "backuplocation" resource
	BackupLocationResourcePlural = "backuplocations"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BackupLocation represents a backuplocation object
type BackupLocation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Location          BackupLocationItem `json:"location"`
}

// BackupLocationItem is the spec used to store a backup location
// Only one of S3Config, AzureConfig or GoogleConfig should be specified and
// should match the Type field. Members of the config can be specified inline or
// through the SecretConfig
type BackupLocationItem struct {
	Type BackupLocationType `json:"type"`
	// Path is either the bucket or any other path for the backup location
	Path         string        `json:"path"`
	S3Config     *S3Config     `json:"s3Config"`
	AzureConfig  *AzureConfig  `json:"azureConfig"`
	GoogleConfig *GoogleConfig `json:"googleConfig"`
	SecretConfig string        `json:"secretConfig"`
}

// BackupLocationType is the type of the backup location
type BackupLocationType string

const (
	// BackupLocationS3 stores the backup in an S3-compliant objectstore
	BackupLocationS3 BackupLocationType = "s3"
	// BackupLocationAzure stores the backup in Azure Blob Storage
	BackupLocationAzure BackupLocationType = "azure"
	// BackupLocationGoogle stores the backup in Google Cloud Storage
	BackupLocationGoogle BackupLocationType = "google"
)

// S3Config speficies the config required to connect to an S3-compliant
// objectstore
type S3Config struct {
	// Endpoint will be defaulted to s3.amazonaws.com by the controller if not provided
	Endpoint        string `json:"endpoint"`
	AccessKeyID     string `json:"accessKeyID"`
	SecretAccessKey string `json:"secretAccessKey"`
	// Region will be defaulted to us-east-1 by the controller if not provided
	Region string `json:"region"`
}

// AzureConfig specifies the config required to connect to Azure Blob Storage
type AzureConfig struct {
	StorageAccountName string `json:"storageAccountName"`
	StorageAccountKey  string `json:"storageAccountKey"`
}

// GoogleConfig specifies the config required to connect to Google Cloud Storage
type GoogleConfig struct {
	ProjectID  string `json:"projectID"`
	AccountKey string `json:"accountKey"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BackupLocationList is a list of ApplicationBackups
type BackupLocationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []BackupLocation `json:"items"`
}

// UpadteFromSecret updated the config information from the secret if not provided inline
func (bl *BackupLocation) UpadteFromSecret(client kubernetes.Interface) error {
	if bl.Location.SecretConfig == "" {
		return nil
	}
	switch bl.Location.Type {
	case BackupLocationS3:
		return bl.getMergedS3Config(client)
	case BackupLocationAzure:
		return bl.getMergedAzureConfig(client)
	case BackupLocationGoogle:
		return bl.getMergedGoogleConfig(client)
	default:
		return fmt.Errorf("Invalid BackupLocation type %v", bl.Location.Type)
	}
}

func (bl *BackupLocation) getMergedS3Config(client kubernetes.Interface) error {
	if bl.Location.S3Config == nil {
		bl.Location.S3Config = &S3Config{}
		bl.Location.S3Config.Endpoint = "s3.amazonaws.com"
		bl.Location.S3Config.Region = "us-east-1"
	}
	if bl.Location.SecretConfig != "" {
		secretConfig, err := client.CoreV1().Secrets(bl.Namespace).Get(bl.Location.SecretConfig, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("error getting secretConfig for backupLocation: %v", err)
		}
		if _, ok := secretConfig.Data["endpoint"]; ok {
			bl.Location.S3Config.Endpoint = strings.TrimSuffix(string(secretConfig.Data["endpoint"]), "\n")
		}
		if _, ok := secretConfig.Data["accessKeyID"]; ok {
			bl.Location.S3Config.AccessKeyID = strings.TrimSuffix(string(secretConfig.Data["accessKeyID"]), "\n")
		}
		if _, ok := secretConfig.Data["secretAccessKey"]; ok {
			bl.Location.S3Config.SecretAccessKey = strings.TrimSuffix(string(secretConfig.Data["secretAccessKey"]), "\n")
		}
		if _, ok := secretConfig.Data["region"]; ok {
			bl.Location.S3Config.SecretAccessKey = strings.TrimSuffix(string(secretConfig.Data["region"]), "\n")
		}
	}
	return nil
}

func (bl *BackupLocation) getMergedAzureConfig(client kubernetes.Interface) error {
	if bl.Location.AzureConfig == nil {
		bl.Location.AzureConfig = &AzureConfig{}
	}
	if bl.Location.SecretConfig != "" {
		secretConfig, err := client.CoreV1().Secrets(bl.Namespace).Get(bl.Location.SecretConfig, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("error getting secretConfig for backupLocation: %v", err)
		}
		if _, ok := secretConfig.Data["storageAccountName"]; ok {
			bl.Location.AzureConfig.StorageAccountName = strings.TrimSuffix(string(secretConfig.Data["storageAccountName"]), "\n")
		}
		if _, ok := secretConfig.Data["storageAccountKey"]; ok {
			bl.Location.AzureConfig.StorageAccountKey = strings.TrimSuffix(string(secretConfig.Data["storageAccountKey"]), "\n")
		}
	}
	return nil
}

func (bl *BackupLocation) getMergedGoogleConfig(client kubernetes.Interface) error {
	if bl.Location.GoogleConfig == nil {
		bl.Location.GoogleConfig = &GoogleConfig{}
	}
	if bl.Location.SecretConfig != "" {
		secretConfig, err := client.CoreV1().Secrets(bl.Namespace).Get(bl.Location.SecretConfig, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("error getting secretConfig for backupLocation: %v", err)
		}
		if _, ok := secretConfig.Data["projectID"]; ok {
			bl.Location.GoogleConfig.ProjectID = strings.TrimSuffix(string(secretConfig.Data["projectID"]), "\n")
		}
		if _, ok := secretConfig.Data["accountKey"]; ok {
			bl.Location.GoogleConfig.AccountKey = strings.TrimSuffix(string(secretConfig.Data["accountKey"]), "\n")
		}
	}
	return nil
}
