# **NeuraKube | Fused cluster intelligence.**

NeuraKube is a machine learning framework that aims to automate the complete lifecycle of ML applications.<br>
Accelerate your machine learning development and deployment processes with an End2End platform system.<br>

NeuraKube simplifies various kubernetes and ML development processes like the setup of GPU/TPU cluster infrastructure,<br>
remote debugging (cloud native development), handling of large datasets (data-providers, filters, enrichment),<br>
training experiment management and production inference.

## Alpha Status: Pre-release

**Caution**: This software is in an alpha development status (WIP).<br>
That means that not every part of the underlying source code is tested or abstracted fully yet.<br>
Neither the less it already works great within many cases and situations - just test it for yourself.<br>
<br>
Feel free to file a GitHub issue if you find any bugs that are bothering you.<br>

## Why NeuraKube?

- Hustle free infrastructure setup (public clouds, on-prem. kubeadm) & management with kubernetes (kubectl abstraction)
- Develop your ML projects cloud native easily with automated remote debugging setup (VSCode, ..)
- Train ML models on public cloud GPU/TPUs with ease (automated setup)
- Integrate and organise large scale datasets/lakes for your projects (soon)
- Deploy scalable APIs for your ML applications

## Setup NeuraCLI

Your kubernetes infrastructure and NeuraKube are easily manageable via NeuraCLI.<br>
The assistant enables a intuitively setup workflow experience and<br>
will guide you through all necessary configuration steps:
<details><summary>Open details</summary>
<p>

### Install NeuraCLI locally
1.1 Download the latest release for your OS: https://github.com/NeuraFuse/NeuraCLI/releases/latest<br>
1.2 Start the NeuraCLI setup file via your terminal with:<br>
```bash
./neuracli-[OS]-[Architecture]
```
1.3 The assistant will guide you through the installation process<br>
</p></details>

## Architecture

The NeuraKube architecture consists of NeuraCLI and NeuraKube which are based primarily on Go and Python.<br>
The framework is deeply integrated with Kubernetes to allow intuitive workflows with microservice based environments.<br>
You can setup NeuraKube on your own infrastructure (on-premise), public clouds (GCloud, AWS, ..), via the NeuraKube Cloud or even locally.<br>

### Versioning
:information_source: Version: **v1alpha1**<br>
:green_circle: Status: **Open Alpha**

## Automated infrastructure setup

Get running quickly on a public cloud (Gcloud/AWS/..) of your choice or setup your own bare metal kubernetes cluster with kubeadm.

<details><summary>Open details</summary>
<p>

#### 2. Start NeuraCLI Assistant

Setup your kubernetes cluster infrastructure easily with the assistant<br>
**(or configure your existing clusters if you like)**:
```bash
neuracli infrastructure create
```
[Terminal Gif]

### Setup NeuraKube

Deploy NeuraKube with NeuraCLI into your cluster to get ready:
```bash
neuracli api create
```
[Terminal Gif]

### Setup details
<details><summary>Open details</summary>
<p>

### Versioning
:information_source: Version: **v1alpha1**<br>
:green_circle: Status: **Open Alpha**

With NeuraKube you can easily deploy and manage a kubernetes cluster that provides for the necessary infrastructure for your ML workloads.
The NeuraCLI Assistant is integrated with the kubernetes-installer which is able to guide you step by step through the bare matel cluster setup process.
You can also automatically setup a new cluster on a public cloud. If you already have a cluster up and running you can directly provide
your kubeconfig credentials to initialize your existing infrastructure.

### Available infrastructure provider plugins
Just provide your credentials for your public cloud account or your already up and running kubernetes cluster:

- [x] **Self hosted** kubernetes cluster (authentication via existing kubeconfig file)
- [x] **Bare metal** setup
- [x] **Google** Cloud
- [ ] **Amazon** AWS (Soon)

### Automated kubernetes bare metal install
Provide only a host IP and SSH configs of a future kubernetes master node.

#### Supported host OS
- [ ] Ubuntu (Soon)
- [ ] CentOS (Soon)
</p></details>
</p></details>

### NeuraKube Cloud
Utilize the full power of NeuraKube without running your own infrastructure.<br>

<details><summary>Open details</summary>
<p>

&#8594; Instant setup<br>
&#8594; No public cloud account necessary<br>

#### Versioning
:information_source: Version: **v1alpha1**<br>
:red_circle: Status: **Closed Alpha**<br>

##### Login to the NeuraKube Cloud with:
```bash
neuracli cloud
```
</p></details>

## Intuitive kubernetes cluster management

NeuraKube has an kubernetes client builtin which provides you with many useful kubectl workflows.<br>
The NeuraCLI Assistant can also help you in various ways to interact faster with your kubernetes clusters.<br>

<details><summary>Open details</summary>
<p>

#### Access container logs fast with:

```bash
neuracli cluster logs [deployment_name]
```

#### Sync files from local working directories into running containers:

```bash
neuracli cluster sync [deployment_name]
```

#### Get a quick overlook over your cluster:

```bash
neuracli cluster inspect
```

[Terminal Gif]
</p></details>

## Cloud native development

Develop your ML models cloud natively on hardware accelerators and benefit from faster prototyping cycles.<br>
Save your developed models in experiment groups to manage a big variety of neural network architectures (soon).<br>

### Remote debugging with GPU/TPUs
&#8594; Quickly launch and connect to your code running on a hardware accelerator with your favorite IDE (VSCode, ..).
<details><summary>Open details</summary>
<p>

#### Connect your python code (soon)
Easily done with the NeuraKube Python Client:

```python
import neurakube.client as NeuraKube
client = NeuraKube.Client()
```

#### After importing the client just hit:
```bash
neuracli develop remote
```

#### Versioning
:information_source: Version: **v1alpha1**<br>
:green_circle: Status: **Open Alpha**

#### Supported base languages

- [x] **Python**
    - [x] PyTorch
    - [x] TensorFlow
    - [x] Other python ML frameworks
- [ ] Golang (Soon)
- [ ] C++ (Soon)

#### Supported IDEs

- [x] **VSCode**
- [ ] IntelliJ (Soon)
[Terminal Gif]
</p></details>

## Large scale data management

Integrate and organise big datasets for your projects easily with the NeuraKube DataLake.<br>
Download and prepare terabytes big datasets in your cloud and connect them with your experiment groups.<br>

<details><summary>Open details</summary>
<p>

### Versioning
:information_source: Version: **v1alpha1**<br>
:red_circle: Status: **Closed Alpha**

#### Providers

- [x] Web crawler
- [x] Common crawl archive

#### Filters/Enrichment

- [x] Text mining from html web pages
</p></details>

## Training with integrated experiment management

NeuraKube provides you with powerful training capabilities on scale.<br>
Train models that are terabytes in size on a cluster of accelerators based on your experiment groups and connected datasets.<br>
Monitor them via the WebUI to track and compare endless different experiments to find the best architectures and hyperparameters for a network (soon).<br>

<details><summary>Open details</summary>
<p>

Start your app with NeuraCLI:

```bash
neuracli app [appID] create
```

### Supported ML plugins

- [x] PyTorch
- [ ] TensorFlow (Soon)

### Versioning
:information_source: Version: **v1alpha1**<br>
:red_circle: Status: **Closed Alpha**

</p></details>

## Production inference

Deploy your trained ML models easily and use the REST client to interact with them via the NeuraKube API.

<details><summary>Open details</summary>
<p>

### Versioning
:information_source: Version: **v1alpha1**<br>
:red_circle: Status: **Closed Alpha**

Send a request to NeuraKube API to register and start an inference server for [appID]:
```bash
neuracli app [appID] inference create
```

Send a request with text data to the now reachable inference server of [appID]:
```bash
neuracli app [appID] inference request text [text_data]
```

### Supported ML plugins

- [x] PyTorch
- [ ] TensorFlow (Soon)
</p></details>

## Community

Join our Slack channel (soon) to interact with people that are also interested in or working on NeuraKube.

### Contributing

Feel free to contribute code improvements and additional plugins to NeuraKube.
Please consider the contribution guide while doing this.

## Copyright & License

NeuraKube is available open source & for free under the Apache License (Version 2.0).