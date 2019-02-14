package manifests

import (
	"fmt"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	ingCfgFilename = filepath.Join(manifestDir, "cluster-ingress-02-config.yml")
)

// Ingress generates the cluster-ingress-*.yml files.
type Ingress struct {
	FileList []*asset.File
}

var _ asset.WritableAsset = (*Ingress)(nil)

// Name returns a human friendly name for the asset.
func (*Ingress) Name() string {
	return "Ingress Config"
}

// Dependencies returns all of the dependencies directly needed to generate
// the asset.
func (*Ingress) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.InstallConfig{},
	}
}

// Generate generates the ingress config and its CRD.
func (ing *Ingress) Generate(dependencies asset.Parents) error {
	installConfig := &installconfig.InstallConfig{}
	dependencies.Get(installConfig)

	config := &configv1.Ingress{
		TypeMeta: metav1.TypeMeta{
			APIVersion: configv1.SchemeGroupVersion.String(),
			Kind:       "Ingress",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "cluster",
			// not namespaced
		},
		Spec: configv1.IngressSpec{
			Domain: fmt.Sprintf("apps.%s.%s", installConfig.Config.ObjectMeta.Name, installConfig.Config.BaseDomain),
		},
	}

	configData, err := yaml.Marshal(config)
	if err != nil {
		return errors.Wrapf(err, "failed to create %s manifests from InstallConfig", ing.Name())
	}

	ing.FileList = []*asset.File{
		{
			Filename: ingCfgFilename,
			Data:     configData,
		},
	}

	return nil
}

// Files returns the files generated by the asset.
func (ing *Ingress) Files() []*asset.File {
	return ing.FileList
}

// Load returns false since this asset is not written to disk by the installer.
func (ing *Ingress) Load(f asset.FileFetcher) (bool, error) {
	return false, nil
}
