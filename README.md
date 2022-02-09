<div align="center">
<br />
<img src="https://user-images.githubusercontent.com/3789273/128085813-92845abd-7c26-4fa2-9f98-928ce2246616.png" width="120px">

# DevLake
<p>

  </p>
  <p>

[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat&logo=github&color=2370ff&labelColor=454545)](http://makeapullrequest.com)
[![Discord](https://img.shields.io/discord/844603288082186240.svg?style=flat?label=&logo=discord&logoColor=ffffff&color=747df7&labelColor=454545)](https://discord.gg/83rDG6ydVZ)
![badge](https://github.com/merico-dev/lake/actions/workflows/test.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/merico-dev/lake)](https://goreportcard.com/report/github.com/merico-dev/lake)


| [English](README.md) | [中文](README-zh-CN.md) |
| --- | --- |
</div>
<br>
<div align="left">

### What is DevLake?
DevLake brings your DevOps data into one practical, customized, extensible view. Ingest, analyze, and visualize data from an ever-growing list of developer tools, with this open source product.

DevLake is designed for developer teams looking to make better sense of their development process and to bring a more data-driven approach to their own practices. You can ask DevLake many questions regarding your development process. Just connect and query.


#### Get started with just a few clicks




<br>


<div align="left">
<img src="https://user-images.githubusercontent.com/2908155/130271622-827c4ffa-d812-4843-b09d-ea1338b7e6e5.png" width="100%" alt="Dev Lake Grafana Dashboard" style="border-radius:15px;" />
<p align="center">Dashboard Screenshot</p><br>
<img src="https://user-images.githubusercontent.com/14050754/145056261-ceaf7044-f5c5-420f-80ca-54e56eb8e2a7.png" width="100%" alt="User Flow" style="border-radius:15px;"/>
<p align="center">User Flow</p>



### Why DevLake?
1. Comprehensive understanding of software development lifecycle, digging workflow bottlenecks
2. Timely review of team iteration performance, rapid feedback, agile adjustment
3. Quickly build scenario-based data dashboards and drill down to analyze the root cause of problems


### What can be accomplished with DevLake?
1. Collect DevOps performance data for the whole process
2. Share abstraction layer with similar tools to output standardized performance data
3. Built-in 20+ performance metrics and drill-down analysis capability
4. Support custom SQL analysis and drag and drop to build scenario-based data views
5. Flexible architecture and plug-in design to support fast access to new data sources

<br>

## Contents

<table>
    <tr>
        <td><b>Section</b></td>
        <td><b>Sub-section</b></td>
        <td><b>Description</b></td>
        <td><b>Documentation Link</b></td>
    </tr>
    <tr>
        <td>Data Sources</td>
        <td>Supported Data Sources</td>
        <td>Links to specific plugin usage & details</td>
        <td><a href="#data-source-plugins">View Section</a></td>
    </tr>
    <tr>
        <td rowspan="3">Setup Guide</td>
        <td>User Setup</td>
        <td>Set up Dev Lake locally as a user</td>
        <td><a href="#user-setup">View Section</a></td>
    </tr>
    <tr>
        <td>Developer Setup</td>
        <td>Set up development environment locally</td>
        <td><a href="#dev-setup">View Section</a></td>
    </tr>
    <tr>
        <td>Cloud Setup</td>
        <td>Set up DevLake in the cloud with Tin</td>
        <td><a href="#cloud-setup">View Section</a></td>
    </tr>
   <tr>
        <td>Tests</td>
        <td>Tests</td>
        <td>Commands for running tests</td>
        <td><a href="#tests">View Section</a></td>
    </tr>
    <tr>
        <td rowspan="4">Make Contribution</td>
        <td>Understand the architecture of DevLake</td>
        <td>See the architecture diagram</td>
        <td><a href="#architecture">View Section</a></td>
    </tr>
    <tr>
        <td>Build a Plugin</td>
        <td>Details on how to make your own plugin</td>
        <td><a href="#plugin">View Section</a></td>
    </tr>
   <tr>
        <td>Add Plugin Metrics</td>
        <td>Guide to add metrics</td>
        <td><a href="#metrics">View Section</a></td>
    </tr>
    <tr>
        <td>Contribution specs</td>
        <td>How to contribute to this repo</td>
        <td><a href="#contributing">View Section</a></td>
    </tr>
    <tr>
        <td rowspan="4">User Guide, Help, and more</td>
        <td>Grafana</td>
        <td>How to visualize the data</td>
        <td><a href="#grafana">View Section</a></td>
    </tr>
    <tr>
        <td>Need Help</td>
        <td>Message us on Discord</td>
        <td><a href="#help">View Section</a></td>
    </tr>
    <tr>
        <td>FAQ</td>
        <td>Frequently asked questions by users</td>
        <td><a href="#faq">View Section</a></td>
    </tr>
    <tr>
        <td>License</td>
        <td>The project license</td>
        <td><a href="#license">View Section</a></td>
    </tr>
</table>

<br>

## Data Sources We Currently Support<a id="data-source-plugins"></a>

Below is a list of _data source plugins_ used to collect & enrich data from specific sources. Each has a `README.md` file with basic setup, troubleshooting, and metrics info.

For more information on building a new _data source plugin_, see [Build a Plugin](plugins/README.md).

Section | Section Info | Docs
------------ | ------------- | -------------
Jira | Summary, Data & Metrics, Configuration, Plugin API | <a href="plugins/jira/README.md" target="_blank">Link</a>
GitLab | Summary, Data & Metrics, Configuration, Plugin API | <a href="plugins/gitlab/README.md" target="_blank">Link</a>
Jenkins | Summary, Data & Metrics, Configuration, Plugin API | <a href="plugins/jenkins/README.md" target="_blank">Link</a>
GitHub | Summary, Data & Metrics, Configuration, Plugin API | <a href="plugins/github/README.md" target="_blank">Link</a>
GitExtractor | Summary, Data & Metrics, Configuration, Plugin API | <a href="plugins/gitextractor/README.md" target="_blank">Link</a>
RefDiff | Summary, Data & Metrics, Configuration, Plugin API | <a href="plugins/refdiff/README.md" target="_blank">Link</a>
<br>

****

## Setup Guide
There're 3 ways to set up DevLake: user setup, developer setup and cloud setup.

### User setup<a id="user-setup"></a>

- If you only plan to run the product locally, this is the **ONLY** section you should need.
- Commands written `like this` are to be run in your terminal.

#### Required Packages to Install<a id="user-setup-requirements"></a>

- [Docker](https://docs.docker.com/get-docker)
- [docker-compose](https://docs.docker.com/compose/install/)

NOTE: After installing docker, you may need to run the docker application and restart your terminal

#### Commands to run in your terminal<a id="user-setup-commands"></a>

1. Clone repository:

   ```sh
   git clone https://github.com/merico-dev/lake.git devlake
   cd devlake
   cp .env.example .env
   ```
2. Start Docker on your machine, then run `docker compose up -d` to start the services.

3. Visit `localhost:4000` to setup configuration files.
   >- Navigate to desired plugins pages on the Integrations page
   >- You will need to enter the required information for the plugins you intend to use.
   >- Please reference the following for more details on how to configure each one:
   >-> <a href="plugins/jira/README.md" target="_blank">Jira</a>
   >-> <a href="plugins/gitlab/README.md" target="_blank">GitLab</a>
   >-> <a href="plugins/jenkins/README.md" target="_blank">Jenkins</a>
   >-> <a href="plugins/github/README.md" target="_blank">GitHub</a>

   >- Submit the form to update the values by clicking on the **Save Connection** button on each form page

   >- `devlake` takes a while to fully boot up. if `config-ui` complaining about api being unreachable, please wait a few seconds and try refreshing the page.
   >- To collect this repo for a quick preview, please provide a Github personal token on **Data Integrations / Github** page.

4. Visit `localhost:4000/create-pipeline` to RUN a Pipeline and trigger data collection.


Pipelines Runs can be initiated by the new "Create Run" Interface. Simply enable the **Data Source Providers**
you wish to run collection for, and specify the **Project ID** for Gitlab and **Repository Name** for GitHub.

Once a valid pipeline configuration has been created, press **Create Run** to start/run the pipeline.
After the pipeline starts, you will be automatically redirected to the **Pipeline Activity** screen
to monitor collection activity.

**Pipelines** is accessible from the **Main Menu** for easy access.

- **Manage All Pipelines** `http://localhost:4000/pipelines`
- **Create Pipeline RUN** `http://localhost:4000/create-pipeline`
- **Pipeline Activity** `http://localhost:4000/pipelines/activity/[RUN_ID]`

For advanced use cases and complex pipelines, please use the Raw JSON API to manually initiate
a run using **cURL** or graphical API tool such as **Postman**. `POST` the following request
to the DevLake API Endpoint.

```json
[
    [
        {
            "Plugin": "github",
            "Options": {
            "repo": "lake",
            "owner": "merico-dev"
            }
        }
    ]
]
```

Please refer to this wiki [How to trigger data collection](https://github.com/merico-dev/lake/wiki/How-to-use-the-triggers-page).

5. Click *View Dashboards* button when done (username: `admin`, password: `admin`). The button will be shown on the Trigger Collection page when data collection has finished.

#### Setup cron job

To synchronize data periodically, we provide [`lake-cli`](./cmd/lake-cli/README.md) for easily sending data collection requests along with [a cron job](./devops/sync/README.md) to periodically trigger the cli tool.

<br>

****

<br>

### Developer Setup<a id="dev-setup"></a>

#### Requirements

- <a href="https://github.com/libgit2/libgit2#quick-start" target="_blank">Libgit2</a>
- <a href="https://docs.docker.com/get-docker" target="_blank">Docker</a>
- <a href="https://golang.org/doc/install" target="_blank">Golang</a>
- Make
  - Mac (Already installed)
  - Windows: [Download](http://gnuwin32.sourceforge.net/packages/make.htm)
  - Ubuntu: `sudo apt-get install build-essential`

#### How to setup dev environment
1. Navigate to where you would like to install this project and clone the repository:

   ```sh
   git clone https://github.com/merico-dev/lake.git
   cd lake
   ```

2. Install Go packages:

    ```sh
    make install
    ```

3. Copy the sample config file to new local file:

    ```sh
    cp .env.example .env
    ```

4. Update the following variables in the file `.env`:
    * `DB_URL`: Replace `mysql:3306` with `127.0.0.1:3306`

5. Start the MySQL and Grafana containers:

    > Make sure the Docker daemon is running before this step.

    ```sh
    docker compose up
    ```

6. Run lake and config UI in dev mode in two seperate terminals:

    ```sh
    # run lake
    make dev
    # run config UI
    make configure-dev
    ```

7. Visit config UI at `localhost:4000` to configure data sources.
   >- Navigate to desired plugins pages on the Integrations page
   >- You will need to enter the required information for the plugins you intend to use.
   >- Please reference the following for more details on how to configure each one:
   >-> <a href="plugins/jira/README.md" target="_blank">Jira</a>
   >-> <a href="plugins/gitlab/README.md" target="_blank">GitLab</a>,
   >-> <a href="plugins/jenkins/README.md" target="_blank">Jenkins</a>
   >-> <a href="plugins/github/README.md" target="_blank">GitHub</a>

   >- Submit the form to update the values by clicking on the **Save Connection** button on each form page

8. Visit `localhost:4000/triggers` to trigger data collection.

   > - Please refer to this wiki [How to trigger data collection](https://github.com/merico-dev/lake/wiki/How-to-use-the-triggers-page). Data collection can take up to 20 minutes for large projects. (GitLab 10k+ commits or Jira 5k+ issues)

9. Click *View Dashboards* button when done (username: `admin`, password: `admin`). The button is shown in the top left.
<br>

****

<br>

### Cloud setup<a id="cloud-setup"></a>
If you want to run DevLake in a clound environment, you can set up DevLake with Tin. [See detailed setup guide](https://github.com/merico-dev/lake/wiki/How-to-Set-Up-Dev-Lake-with-Tin)

**Disclaimer:**
> To protect your information, it is critical for users of the Tin hosting to set passwords to protect DevLake applications. We built DevLake as a self-hosted product, in part to ensure users have total protection and ownership of their data, while the same remains true for the Tin hosting, this risk point can only be eliminated by the end-user.

<br>

## Tests<a id="tests"></a>

To run the tests:

```sh
make test
```
<br>

## Make Contribution
This section list all the documents to help you contribute to the repo.

### Understand the Architecture of DevLake<a id="architecture"></a>
![devlake-architecture](https://user-images.githubusercontent.com/14050754/143292041-a4839bf1-ca46-462d-96da-2381c8aa0fed.png)
<p align="center">Architecture Diagram</p>

### Add a Plugin<a id="plugin"></a>

[plugins/README.md](/plugins/README.md)

### Add Plugin Metrics<a id="metrics"></a>

[plugins/HOW-TO-ADD-METRICS.md](/plugins/HOW-TO-ADD-METRICS.md)

### Contributing Spec<a id="contributing"></a>

[CONTRIBUTING.md](CONTRIBUTING.md)

<br>

## User Guide, Help and more
### Grafana<a id="grafana"></a>

We use <a href="https://grafana.com/" target="_blank">Grafana</a> as a visualization tool to build charts for the data stored in our database. Using SQL queries, we can add panels to build, save, and edit customized dashboards.

All the details on provisioning and customizing a dashboard can be found in the [Grafana Doc](docs/GRAFANA.md).



### Need help?<a id="help"></a>

Message us on <a href="https://discord.com/invite/83rDG6ydVZ" target="_blank">Discord</a>


### FAQ<a id="faq"></a>

Q: When I run ``` docker-compose up -d ``` I get this error: "qemu: uncaught target signal 11 (Segmentation fault) - core dumped". How do I fix this?

A: M1 Mac users need to download a specific version of docker on their machine. You can find it <a href="https://docs.docker.com/desktop/mac/apple-silicon/" target="_blank">here</a>.


### License<a id="license"></a>

This project is licensed under Apache License 2.0 - see the [`LICENSE`](LICENSE) file for details.
