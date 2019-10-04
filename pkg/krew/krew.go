package krew

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"text/template"

	"github.com/google/go-github/github"
	"github.com/rajatjindal/krew-plugin-release/pkg/actions"
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/krew/pkg/constants"
)

// //UpdatePluginManifest2 is 2
// func UpdatePluginManifest2(baseDir, pluginName string, release *github.RepositoryRelease) error {
// 	pluginManifest, err := indexscanner.ReadPluginFile(filepath.Join(actions.GetWorkspace(), ".krew.yaml"))
// 	if err != nil {
// 		return err
// 	}

// 	for _, platform := range pluginManifest.Spec.Platforms {
// 		t := template.New("uri").Parse(platform.URI)
// 		buf := new(bytes.Buffer)
// 		t.Execute(buf, release)

// 		sha256, err := getSha256ForAsset(buf.String())
// 		if err != nil {
// 			return err
// 		}

// 	}

// 	return nil
// }

//UpdatePluginManifest updates the manifest with latest release info
func UpdatePluginManifest(baseDir, pluginName string, release *github.RepositoryRelease) error {
	processedPluginBytes, err := processPluginTemplate(release)
	if err != nil {
		return err
	}

	// pluginFileWithSha256, err := addSha256ToPluginFile(pluginName, processedPluginBytes)
	// if err != nil {
	// 	return err
	// }

	pluginsFile := filepath.Join(baseDir, "plugins", pluginFileName(pluginName))
	return ioutil.WriteFile(pluginsFile, processedPluginBytes, 0644)
}

func processPluginTemplate(releaseInfo *github.RepositoryRelease) ([]byte, error) {
	t := template.New(".krew.yaml").Funcs(map[string]interface{}{
		"addURIAndSha": func(url, tag string) string {
			t := struct {
				TagName string
			}{
				TagName: tag,
			}
			buf := new(bytes.Buffer)
			temp, err := template.New("url").Parse(url)
			if err != nil {
				panic(err)
			}

			err = temp.Execute(buf, t)
			if err != nil {
				panic(err)
			}

			logrus.Infof("getting sha256 for %s", buf.String())
			sha256, err := getSha256ForAsset(buf.String())
			if err != nil {
				panic(err)
			}

			return fmt.Sprintf(`uri: %s
    sha256: %s`, buf.String(), sha256)
		},
	})

	templateObject, err := t.ParseFiles(filepath.Join(actions.GetWorkspace(), ".krew.yaml"))
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

func pluginFileName(name string) string {
	return fmt.Sprintf("%s%s", name, constants.ManifestExtension)
}
