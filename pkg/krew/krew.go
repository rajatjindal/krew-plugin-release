package krew

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"text/template"

	"github.com/google/go-github/github"
	"github.com/kubernetes-sigs/krew/pkg/index/indexscanner"
	"github.com/rajatjindal/krew-plugin-release/pkg/actions"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"sigs.k8s.io/krew/pkg/constants"
)

//UpdatePluginManifest updates the manifest with latest release info
func UpdatePluginManifest(baseDir, pluginName string, release *github.RepositoryRelease) error {
	processedPluginBytes, err := processPluginTemplate(release)
	if err != nil {
		return err
	}

	pluginFileWithSha256, err := addSha256ToPluginFile(pluginName, processedPluginBytes)
	if err != nil {
		return err
	}

	pluginsFile := filepath.Join(baseDir, "plugins", pluginFileName(pluginName))
	return ioutil.WriteFile(pluginsFile, pluginFileWithSha256, 0644)
}

func processPluginTemplate(releaseInfo *github.RepositoryRelease) ([]byte, error) {
	templateObject, err := template.ParseFiles(filepath.Join(actions.GetWorkspace(), ".krew.yaml"))
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = templateObject.Execute(buf, releaseInfo)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func addSha256ToPluginFile(pluginName string, pluginFileBytes []byte) ([]byte, error) {
	d, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(filepath.Join(d, pluginFileName(pluginName)), pluginFileBytes, 0644)
	if err != nil {
		return nil, err
	}

	pluginManifest, err := indexscanner.ReadPluginFile(filepath.Join(d, pluginFileName(pluginName)))
	if err != nil {
		return nil, err
	}

	for _, platform := range pluginManifest.Spec.Platforms {
		logrus.Infof("getting sha for %s", platform.URI)
		sha256, err := getSha256ForAsset(platform.URI)
		if err != nil {
			return nil, err
		}

		platform.Sha256 = sha256
	}

	return yaml.Marshal(pluginManifest)
}

func pluginFileName(name string) string {
	return fmt.Sprintf("%s%s", name, constants.ManifestExtension)
}
