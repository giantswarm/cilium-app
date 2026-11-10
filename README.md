[![CircleCI](https://circleci.com/gh/giantswarm/cilium-app.svg?style=shield)](https://circleci.com/gh/giantswarm/cilium-app)

# cilium chart

Giant Swarm offers a cilium App for workload clusters.

## Installing

This app is installed by default in every workload cluster, as part of the Giant Swarm release.

## Upgrading cilium version

Update the version in `vendir.yml` and run `./sync/sync.sh`. It's good to check the diffs from upstream. They are stored in the `./diffs` directory.


## Version compatibility

| Cilium app | Cilium version | CAPA | CAPZ | CAPV | CAPVCD | EKS | AKS |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 1.7.x | 1.20.x | -- | -- | -- | -- | -- | -- |
| 1.6.x | 1.20.x | v36 | v36 | v36 | v36 | v36 | v36 |
| 1.5.x | 1.19.x | v35 | v35 | v35 | v35 | -- | -- |
| 1.4.x | 1.19.x | v34 | v34 | v34 | v34 | -- | -- |
| 1.3.x | 1.18.x | v32,v33,v34 | v32,v33,v34 | v32,v33,v34 | v32,v33,v34 | v32,v34 | -- |
| 1.2.x | 1.17.x | v31 | v31 | v31 | v31 | -- | -- |
| 1.1.x | 1.17.x | -- | -- | -- | -- | -- | -- |
| 1.0.x | 1.17.x | -- | -- | -- | -- | -- | -- |
| 0.32.x | 1.17.x | -- | -- | -- | -- | -- | -- |
| 0.31.x | 1.16.x | v30 | v30 | v30 | v30 | -- | -- |
| 0.30.x | 1.16.x | -- | -- | -- | -- | -- | -- |
| 0.29.x | 1.16.x | -- | -- | -- | -- | -- | -- |
| 0.28.x | 1.16.x | -- | -- | -- | -- | -- | -- |
| 0.27.x | 1.16.x | -- | -- | -- | -- | -- | -- |
| 0.26.x | 1.16.x | -- | -- | -- | -- | -- | -- |
| 0.25.x | 1.15.x | v25-v29 | v25-v29 | v27-v29 | v27-v29 | -- | -- |
