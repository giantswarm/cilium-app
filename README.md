[![CircleCI](https://circleci.com/gh/giantswarm/cilium-app.svg?style=shield)](https://circleci.com/gh/giantswarm/cilium-app)

# cilium chart

Giant Swarm offers a cilium App which can be installed in workload clusters.

## Installing

There are several ways to install this app onto a workload cluster.

- [Using our web interface](https://docs.giantswarm.io/ui-api/web/app-platform/#installing-an-app).
- By creating an [App resource](https://docs.giantswarm.io/ui-api/management-api/crd/apps.application.giantswarm.io/) in the management cluster as explained in [Getting started with App Platform](https://docs.giantswarm.io/app-platform/getting-started/).

## Upgrading cilium version

The contents of the `helm` folder are being generated by the `make` target called `make update-chart`.

This target uses [`vendir`](https://carvel.dev/vendir/) to fetch the helm chart contained in [the fork of the cilium repository that we maintain](https://github.com/giantswarm/cilium-upstream).
Currently, the `main` branch on the fork contains the upstream tag `v1.14.3`, with our custom changes on top.

If you want to upgrade this `cilium-app` to use a newer version of cilium, you need to prepare our fork first. All changes are applied there, and then we use `vendir` to generate the contents of the `helm` folder.

We need to create a new branch on our fork based off the tag of the version we want to upgrade to. For example, if we want to upgrade to cilium `v1.14`, we need to create a new branch based off the upstream tag for the latest `v1.14` version, which is currently `v1.14.3`.
```
git fetch upstream
git checkout upstream/v1.14.3
git checkout -b update-to-1-14-3
```

Then we need to apply our custom changes on top of that new branch. Since those changes are already on the repository (although a different branch), we can use `cherry-pick` for apply those same commits to our newly created branch, for example
```
git cherry-pick a4b22dee87ba3663f967f6dd6d8e666c849c742d^..25c449534cc325a5798fc7c839b8ac33591b3516
```

It's probable that conflicts will happen, so we need to fix those when applying the commits.
One last thing we need to do in our fork is to update the `values.schema.json` file, because upstream does not provide one. You can do it with
```
helm schema-gen install/kubernetes/cilium/values.yaml > install/kubernetes/cilium/values.schema.json
```

Don't forget to commit the changes, if any.
Once we are done, we can push our new branch to our fork and it will be ready to be used by `vendir`.

Then, in this repository, we need to update the `vendir` configuration in `vendir.yml` to use the new branch we just pushed, and run the make target `APPLICATION=cilium make update-chart`.
With the generated changes, let's create a new pull request so that everyone can review the changes that will be applied to the cilium chart.
If we need further customizations, we can keep adding commits on the new branch on the fork, and re-run `make update-chart` to update the generated files.

Once we are happy with the changes, we can merge the changes in the fork to the fork's `main` branch. After merging, the `main` branch in our fork should contain cilium latest release (`v1.14.3` on the example) with our customizations on top.
Then update our pull request to use the fork's `main` branch again, and merge it.
Finally, we just need to create a new release of this cilium-app.
