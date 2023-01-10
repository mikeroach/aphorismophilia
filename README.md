[![Go Report Card](https://goreportcard.com/badge/github.com/mikeroach/aphorismophilia)](https://goreportcard.com/report/github.com/mikeroach/aphorismophilia)

# Aphorismophilia, n.:

1. Love of aphorisms.
1. Flimsy pretense to gain hands-on experience with modern technology trends.

This is a proof-of-concept project to create a containerized Go "hello world" web application deployed to GitOps-managed IaaS/PaaS resources via CI, CD, & Infrastructure-as-Code pipelines. **See it in action [here](http://www.borrowingcarbon.net)!**

### Table of Contents
1. [Motivation](#motivation)
1. [Technology Justification](#technology-justification)
1. [Guiding Principles](#guiding-principles)
1. [Architecture](#architecture)
1. [Lessons Learned and Stray Thoughts](#lessons-learned-and-stray-thoughts)
1. [Next Steps](#next-steps)

### Motivation

I made this to gain hands-on experience implementing GitOps, Infrastructure-as-Code, Continuous Integration, Continuous Delivery, container orchestration, and Site Reliability Engineering concepts with a Terraform, Kubernetes, Jenkins, and Go-based stack on Google Cloud Platform.

During my past several years in an Operations/SRE leadership role supporting teams (aka "just talking about work instead of actually _doing_ it"), I found myself focusing most of my personal growth on soft skills rather than technical ones. While this was fulfilling in its own right, I also didn't want to let my technical knowledge atrophy beyond irrelevance as a technological renaissance obsoleted many of the practices I'd adopted throughout my career. Fortunately I had the privilege of working with exceptionally smart and talented engineers fully capable of adopting new industry trends into their work, but I felt somewhat disingenuous preaching about how certain emerging tech concepts would help us become more efficient when I hadn't yet practiced them personally.

I decided to take the opportunity while between jobs to catch up towards the state of the industry with a learning project which would satisfy my own curiosity and better prepare me for a new job involving contemporary engineering disciplines.

### Technology Justification

Why did I want to learn more about _these_ specific technologies?

* **Terraform:** All of my prior infrastructure-as-code experience has been with Salt's boto modules and CloudFormation, so I appreciated Terraform's vendor-agnostic and infrastructure-oriented (vs. configuration management-oriented) design as well as its widespread third-party provider and module support.

* **Kubernetes:** I was excited to hear how Kubernetes exposes primitives for orchestrating availability, scalability, deployment, and general operability considerations that I've spent most of my career addressing through more complex processes and conventions.

* **Jenkins:** Most software organizations where I've worked have used Hudson or Jenkins. While this longtime defacto CI tool seems to add modern CD concepts as an afterthought, its versatile plugin ecosystem and wide install base suggests it's worth learning even though my brief research revealed a variety of other strong cloud-native CI/CD solutions.

* **Go:** I was intrigued by Go's design philosophy around maintainability and suitability for modern web services, and wanted to challenge myself with a more structured and demanding programming language than Python.

* **Google Cloud Platform:** I wanted to try GCP for the novelty since most of my compute provider experience involves AWS. Their $300 worth of promotional credit for new accounts plus compute instance hours included in the Always Free tier also nicely supported my Smart Spending goal, and I was delighted to find that GCP's free Kubernetes control-plane saves me approximately $150/month per cluster compared to EKS.

### Guiding Principles

At first I didn't even know what I didn't know about what I wanted to learn, let alone exactly how I'd implement the broad goal of "write a Go service and deploy it with IaC CI/CD". So I defined a few guidelines to help drive my decision making process:

* **Smart Spending ®™:** Cost consciousness is a pillar of operational excellence regardless of whose money I'm spending, and when it's my own cash on a hobby project the budget should be comparable to what I'd pay PG&E to power a junker laptop's power supply.

* **Don't Shave the Yak:** I won't let perfect be the enemy of done. I'll get a minimum viable result working that satisfies the basic criteria for what I want to learn, and create a backlog of the things I'd like to improve and learn later.

* **Do It Yourself to Learn It Yourself:** The point of this project is to learn fundamentals, even if that means doing things the hard way and inventing squarish wheels when there are already perfectly round ones that do the job more elegantly. I'm not learning anything if I'm just blindly running someone else's software or copy-pasting from a Medium article/StackExchange post. Besides, I'll appreciate the abstracted solutions that much more once I'm familiar with the details of how these technologies work.

* **Transparency Builds Trust - Don't Keep Too Many Secrets:** Sunlight is the best disinfectant, so I'll come up with better solutions knowing that other people will review my work (that's you, constant reader!). Once I'm happy with version 1.0, I'll publish the entirety of what's required to construct and maintain my environments in public GitHub repositories without any hand-waving or hidden magic - with the notable exception of encrypting sensitive data like credentials and non-functional environment details.

* **Declarative is a Great Narrative:** Since I want source code to represent the desired state of the environment, my deployment processes should be derived from source code management strategies. I'll use a feature branching workflow for testing changes, pull request to `main` when I'm ready to merge, and once it's in `main` then it's shipping. Copy/pasting and nifty scripty CLI work has its place in IDE/local development, but the brass ring is limiting write interactions with my online environments to Git operations.

* **Ship It, Theseus:** Once launched, I'll maintain this environment on an ongoing basis for other tech experiments and migrate services and persistent data as necessary to gracefully maintain functionality when replatforming.

### Architecture

I develop and deploy the aphorismophilia service along with all its supporting infrastructure from these repositories:

| Repository | Purpose |
| ---------- | ------- |
| **[iac-bootstrap](https://github.com/mikeroach/iac-bootstrap)** | Terraform to instantiate a management and Jenkins build environment from which to run IaC and application CI/CD pipelines. |
| **[iac-bootstrap-salt](https://github.com/mikeroach/iac-bootstrap-salt)** | SaltStack masterless states and pillar for management host configuration. |
| **[iac-template-pipeline](https://github.com/mikeroach/iac-template-pipeline)** | Terraform module pipeline for templated infrastructure and application stacks\*. |
| **[iac-pipeline-auto](https://github.com/mikeroach/iac-pipeline-auto)** | Pipeline for environments receiving automatic continuous delivery deployments. |
| **[iac-pipeline-gated](https://github.com/mikeroach/iac-pipeline-gated)** | Pipeline for environments requiring manual approval to perform deployments. |
| **[aphorismophilia](https://github.com/mikeroach/aphorismophilia)** | Go "hello world" containerized webapp source and Jenkins-based CI pipeline. |
| **[aphorismophilia-terraform](https://github.com/mikeroach/aphorismophilia-terraform)** | Terraform module to manage a Kubernetes deployment of the aphorismophilia webapp\*. |

\* I was inspired to experiment with this approach by Kief Morris' [Template Stack Pattern](https://infrastructure-as-code.com/patterns/stack-replication/template-stack.html) as described in [Infrastructure as Code](https://infrastructure-as-code.com) and Nicki Watt's [Terraservices presentation slides](https://www.slideshare.net/opencredo/hashidays-london-2017-evolving-your-infrastructure-with-terraform-by-nicki-watt) and [video](https://www.youtube.com/watch?v=wgzgVm7Sqlk).

### Lessons Learned and Stray Thoughts

* My emphasis on budget consciousness drove some awkward technical implementation decisions not suitable for production traffic (like single node K8s clusters with NodePort-based Nginx ingress instead of load balancers), but... by Grabthar's hammer... what a savings.

* There's sophisticated cloud-native delivery tooling like [Jenkins-X](https://jenkins-x.io), [Argo CD](https://argoproj.github.io/argo-cd/), [Weaveworks Flux](https://www.weave.works/oss/flux/), and [Spinnaker](https://www.spinnaker.io) that would make a much better foundation for a real environment than what I cobbled together.

* I like being able to deploy a complete stack of applications and infrastructure as one discrete unit. In a collaborative environment it may be better to instead manage various components separately; modeled around both the organizational structure of the people involved and their skillsets/appetite for complexity.

* At one point while learning how to build a Go web application, I empathized with my development counterparts' sentiments from jobs past: I just wanted to work on writing new features without having to spend all that time on the details of getting it into production infrastructure!

* Using Terraform's Kubernetes provider to manage my clusters is more awkward than I thought it would be, and it was especially tedious to translate YAML into HCL by hand before thinking to look for [this](https://github.com/sl1pm4t/k2tf). At least it was convenient to limit the number of new things I learned at once (e.g. Helm & Tiller or Not to Tiller).

* I find it ironic that so much of modern "DevOps" techniques come down to string manipulation. Not that tossing a bunch of Makefiles into the mix is any better.

### Next Steps

1. Add a diagram of the architecture and workflow; a picture is worth 1,000+ lines of code.
1. Transcribe my yak shaving notes into an actionable backlog now that I have a basic working environment.
1. Implement true IaC unit and integration testing (e.g. Terratest, kitchen-terraform, InSpec, et al.) instead of CI-"lite".
