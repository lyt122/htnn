This doc tracks how to maintain the source code of HTNN.

## Release a new version

To release a new version, please follow the steps below:

1. Create a new release branch `release/v${version}` from the main branch. Do the work below on the new branch.
2. Create tag `api/v${version}`.
3. Commit the changes below (the CI will fail at this point):
    * Update those `go.mod` which depend on `mosn.io/htnn/$mod`.
4. Create tag `types/v${version}` for `types` module. Then do the same with `controller` and `plugins`. Rerun the `test` workflow to verify the changes. Don't panic for "server response: not found" error. The sync of sum.golang.org might take half an hour. Try again later.
5. Create tag `image/v${version}` to trigger image building.
6. Submit a new commit with the changes below:
    * Once the image is ready, update the version in the `manifests/charts/*/Chart.yaml`.
    * Run `make fmt-go`.
7. The CI will create a new chart package.
8. Merge the release branch to the main branch.

## Upgrade components

### Upgrade Istio

To upgrade Istio, please follow the steps below:

* Discuss the impact of the upgrade. For example, is there any break change, do we need to upgrade K8S, etc.
* Update the base image used in the integration tests.
* Update the ISTIO_VERSION we define in the `common.mk`.
* Update the link `/envoy/v1.xx.y/configuration/` in the doc to the new Envoy version. And `istio/istio/xxx` to the new Istio version.
* Update the charts' dependency versions used in the `manifests/charts/*/Chart.yaml`.

If this is a minor version upgrade, please follow the additional steps below:

* Sync the `manifests/charts/htnn-controller/*` to the latest istio's istiod chart.
