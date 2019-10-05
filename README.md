# krew-plugin-release

`krew-plugin-release` is a [github action](https://github.com/features/actions) that opens a PR with upstream [krew-index](https://github.com/kubernetes-sigs/krew-index) on new release of your plugin.

# Installing
#### create token
- Navigate to [tokens](https://github.com/settings/tokens)
- Click on 'Generate new token'
- Choose a name e.g. `krew-plugin-release-token`
- Select `repo.public_repo` checkbox
- Generate and copy token
#### add secret to plugin repo
- Navigate to your github repo that has plugin code
- Go to Settings
- Click on `Secrets` on left nav bar
- Click on 'Add new secret'
- Give name of the secret to be `KREW_PLUGIN_RELEASE_TOKEN` (Note: this has to exactly same like this)
- Add the token you created above and save the secret

#### create a .krew.yaml file in your plugin git repository. 
Below is the file created for [modify-secret](https://github.com/rajatjindal/kubectl-modify-secret). Note that this is a template file, and you can templatize your download uri. The action will process template, fetch the `.tar.gz` file from processed uri, calculate sha, writes the processed plugin manifest and then opens PR in the `krew-index` repo with these changes.

Look at this [sample PR](https://github.com/kubernetes-sigs/krew-index/pull/266) created using this action, reviewed using bot from [@ahmetb](https://github.com/ahmetb) and merged using [@k8s-ci-robot](https://github.com/k8s-ci-robot)

```yaml
apiVersion: krew.googlecontainertools.github.com/v1alpha2
kind: Plugin
metadata:
  name: modify-secret
spec:
  version: "{{ .TagName }}"
  platforms:
  - selector:
      matchLabels:
        os: darwin
        arch: amd64
    {{addURIAndSha "https://github.com/rajatjindal/kubectl-modify-secret/releases/download/{{ .TagName }}/darwin-amd64-{{ .TagName }}.tar.gz" .TagName }}
    files:
    - from: "*"
      to: "."
    bin: kubectl-modify-secret
  - selector:
      matchLabels:
        os: linux
        arch: amd64
    {{addURIAndSha "https://github.com/rajatjindal/kubectl-modify-secret/releases/download/{{ .TagName }}/linux-amd64-{{ .TagName }}.tar.gz" .TagName }}
    files:
    - from: "*"
      to: "."
    bin: kubectl-modify-secret
  shortDescription: modify secret with implicit base64 translations
  caveats: |
    Usage:
      kubectl modify-secret secret-name -n kube-system
    
    Read the documentation at:
      https://github.com/rajatjindal/kubectl-modify-secret
  description: |
    Usage:
      kubectl modify-secret secret-name -n kube-system

      This plugin fetches the given secret manifest from the cluster, decodes 
      the payload, opens an editor to make changes, and applies the modified 
      manifest when done.
  homepage: https://github.com/rajatjindal/kubectl-modify-secret
```

#### Add following file to your .github/workflows dir

Add the following file, and change the plugin name to your plugin name.

```yaml
name: main
on: 
  release: 
    types: ["published"]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v1 
      - name: Open PR for new release of Krew Plugin
        uses: rajatjindal/krew-plugin-release@v0.0.30
        with:
          plugin-name: 'your-plugin-name'
          upstream-krew-index-owner: 'kubernetes-sigs'
        env:
          KREW_PLUGIN_RELEASE_TOKEN: ${{ secrets.KREW_PLUGIN_RELEASE_TOKEN }}
```

# Bug Reports

if you find any issue with this action, please report an issue, and we will try to fix it asap.

# Security Issue
if you find any security issue with this action, please email to rajatjindal83 at gmail dot com

