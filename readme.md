# Cloud Workstation Tool

Helper tool to start, stop and start a tunnel to connect to [Google Cloud Workstation](https://docs.cloud.google.com/workstations/docs/overview).

## Installation

Download from the [Releases](https://github.com/xaviatwork/cloud-workstation-tool/releases) page.

Cloud Workstation Tool is a single binary; no need to install ¯\\\_(ツ)_/¯

### Create the configuration folder

* MacOS: `mkdir -p ~/.config/cloud-workstation-config`
* Windows: `mkdir %USERPROFILE%\cloud-workstation-config`

### Create the configuration file

Create a JSON configuration file in the Cloud Workstation Tool configuration folder with the following structure:

```json
{
    "name": "WORKSTATION_NAME",
    "project_id": "GOOGLE_CLOUD_PROJECT_ID",
    "region": "REGION",
    "cluster": "CLOUD_WORKSTATION_CLUSTER_NAME",
    "cluster_config": "CLUSTER_CONFIG",
    "local_port": "LOCAL_PORT"
}
```

* WORKSTATION_NAME: Name of the Cloud Workstation. Example: `my-workstation`
* GOOGLE_CLOUD_PROJECT_ID: Google Cloud Project ID where the Workstation is located.
* REGION: Google Cloud region where the Workstation is located.
* CLOUD_WORKSTATION_CLUSTER_NAME: Name of the cluster managing the Cloud Workstation
* CLUSTER_CONFIG: Name of the Cloud Workstation configuration used by the Cloud Workstation
* LOCAL_PORT: Local port to connect to the SSH tunnel

## Run

To start the tunnel to the Google Cloud Workstation, run:

> It is recommended to rename `cw-tunnel-windows-amd64` to `cw-tunnel` in Windows.

```shell
$ cw-tunnel

Configuration:
- Workstation: WORKSTATION_NAME
- Project: GOOGLE_CLOUD_PROJECT_ID
- LocalPort: LOCAL_PORT
- Cluster: CLOUD_WORKSTATION_CLUSTER_NAME
- Config: CLUSTER_CONFIG
- Region: REGION
Starting tunnel to WORKSTATION_NAME:22 on localhost:LOCAL_PORT...
[gcloud] Listening on port [LOCAL_PORT].
✅ Tunnel is ready! You can now connect to localhost:LOCAL_PORT

```

## TODO

* [] Add command to create configuratio folder and empty config file
* [] Add `start` command
* [] Add `stop` command
* [] Unify `start`, `stop` and `tunnel` under a single `cw` tool (`cw start`, `cw stop`, `cw tunnel`)
* [] Custom configuration name and location ( `cw --config /path/to/cw.cfg`)
