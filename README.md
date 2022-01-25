# Kaleido Project [![Build Status](https://app.travis-ci.com/mrshah-at-ibm/solidity-project.svg?token=69E7mUEaHCqqAB5inXkw&branch=main)](https://app.travis-ci.com/mrshah-at-ibm/solidity-project)

This is a complete project to setup and run an app that meets the following requirements:

1. Sign Transactions

   The app [generates](https://github.com/mrshah-at-ibm/solidity-project/blob/5ca6bc26ee39c4184d29acd2b330e3c74809ba02/pkg/executer/utils.go#L11-L47) a private key on start and [saves it in a secret](https://github.com/mrshah-at-ibm/solidity-project/blob/5ca6bc26ee39c4184d29acd2b330e3c74809ba02/pkg/config/config.go#L155).

2. Deploy an instance of ERC-721 contract, and retain knowledge of its contract

   The app [deploys](https://github.com/mrshah-at-ibm/solidity-project/blob/5ca6bc26ee39c4184d29acd2b330e3c74809ba02/pkg/executer/executer.go#L75-L128) instance of an ERC-721 contract and [saves](https://github.com/mrshah-at-ibm/solidity-project/blob/5ca6bc26ee39c4184d29acd2b330e3c74809ba02/pkg/config/config.go#L113-L216) the address in a configmap

3. Mint, burn and transfer ERC-721 tokens

   The app provides API to [mint](https://github.com/mrshah-at-ibm/solidity-project/blob/5ca6bc26ee39c4184d29acd2b330e3c74809ba02/pkg/config/config.go#L113-L216), [burn](https://github.com/mrshah-at-ibm/solidity-project/blob/5ca6bc26ee39c4184d29acd2b330e3c74809ba02/pkg/config/config.go#L113-L216) and [transfer](https://github.com/mrshah-at-ibm/solidity-project/blob/5ca6bc26ee39c4184d29acd2b330e3c74809ba02/pkg/routes/routes.go#L180-L242) ERC-721 tokens

4. The app should be scaleable

   The app can be scaled by updating the replicas of the deployment spec (either using [terraform](https://github.com/mrshah-at-ibm/solidity-project/blob/5ca6bc26ee39c4184d29acd2b330e3c74809ba02/terraform/app/app.tf#L74) or manually).

5. The app should not use local storage

   The app stores and gets addresses and private keys. These are directly from kubernetes configmaps and secrets.

## Additional features:

1. TLS enabled for the API

   [cert-manager](https://github.com/mrshah-at-ibm/solidity-project/blob/5ca6bc26ee39c4184d29acd2b330e3c74809ba02/terraform/cert-manager/cert-manager.tf) based certificates are issued. terraform is provided for [certificates](https://github.com/mrshah-at-ibm/solidity-project/blob/5ca6bc26ee39c4184d29acd2b330e3c74809ba02/terraform/app/app.tf#L74) to be generated using AKME (Let's Encrypt) servers.

2. IaC to deploy network, certificate manager & app

   terraform is provided to deploy [network](./terraform/network), [certificate manager](./terraform/cert-manager) and [app](./terraform/app)

3. Tests

   Repository is clean as far as go-sec goes. A few unit-tests have been written using ginkgo.

4. Travis integration

   [Travis integration](https://app.travis-ci.com/github/mrshah-at-ibm/solidity-project) has been setup on the repository to perform CI. CD is not setup yet.

5. Option to use multiple addresses for the app itself

   Started with [ClaimAddress](https://github.com/mrshah-at-ibm/solidity-project/blob/5ca6bc26ee39c4184d29acd2b330e3c74809ba02/pkg/config/config.go#L113) to claim and [maintain ownership](https://github.com/mrshah-at-ibm/solidity-project/blob/5ca6bc26ee39c4184d29acd2b330e3c74809ba02/pkg/config/config.go#L155) of an address exclusively. Will switch to UseAddress to use one of the addresses in the list non-exclusively.

6. Tested on multiple types of kubernetes clusters

   Tested on IBM Cloud Free Cluster - this only works with NodePort type of service as ingress is not setup. So no certificates/ingress/tls setup.

   Tested on Azure Kubernetes Service - this has full ingress option to test with ingress/DNS/etc.

7. [Not a feature] The contract

   The [contract](./contracts/default.sol) is a copy of the preset provided by OpenZeppelin. The reason I made a copy of it and tried to edit is to be able to make successful `eth_calls` which you will read about in the Future section of the README.

8. In cluster communication

   The app communicates with the network using the `service.namespace` hostnames. All the nodes in network communicate in the same way.

9. Quick Build and Test loop

   The go build will take a long time to compile inside the container build as it cannot use any cache. So I started with the base build and pushed the image as `mrshah/kp:app-orig`. Then I would build the app locally and replace only the binary in the image. This was pushed as the image `mrshah/kp:app` and used everywhere. I would also restart the pod and watch it run - all in [one script](./tools/restart_and_watch.sh).

10. Auth

   The app is protected by tokens which can only be generated after authenticating against Github's SSO.

## Usage

### Prerequisites

1. [terraform cli](https://learn.hashicorp.com/tutorials/terraform/install-cli)

2. curl

3. kubeconfig.yaml for your cluster

   If you have kubectl installed, you can run `kubectl config view --flatten > kubeconfig.yaml` to get the kubeconfig for your cluster.

### Optional Prerequisites

   If you want to mess with stuff and look deeper, you will need the following tools.

1. [kubectl](https://kubernetes.io/docs/tasks/tools/)

### Deployment

_Note: The namespaces are hardcoded in terraform. Please update them if needed._

1. Network

   Go to the [terraform](./terraform/network) folder for network. Put kubeconfig.yaml to that folder and run apply it.

   ```bash
   cd ./terraform/network
   
   # Get kubeconfig file
   kubectl config view --flatten > kubeconfig.yaml
   
   # Review plan to be applied
   terraform plan

   # Apply terraform
   terraform apply -auto-approve
   ```

2. [Optional] Install cert-manager

   Go to the [terraform](./terraform/cert-manager) folder for cert-manager. Put kubeconfig.yaml to that folder and run apply it.

   ```bash
   cd ./terraform/cert-manager
   
   # Get kubeconfig file
   kubectl config view --flatten > kubeconfig.yaml
   
   # Review plan to be applied
   terraform plan

   # Apply terraform
   terraform apply -auto-approve
   ```

3. Install app

   Go to the [terraform](./terraform/cert-manager) folder for app. Put kubeconfig.yaml to that folder and run apply it.

   ```bash
   cd ./terraform/app
   
   # Get kubeconfig file
   kubectl config view --flatten > kubeconfig.yaml
   
   # Review plan to be applied
   terraform plan

   # Apply terraform
   terraform apply -auto-approve
   ```

4. Test the app

   _Notes: The app is currently running on https://app.mrshah.space_

   Go to the tests folder and run some shell scripts to test the APIs. Make sure to point the scripts to the domain of your app.

   ```bash
   cd ./tests

   DOMAIN="https://app.mrshah.space"
   NUMBER_OF_CALLS="5"
   AUTH_TOKEN="<Optional - put your auth token here>"

   # Test the app is running - it should return "Server running ok
   curl "${DOMAIN}"

   # Mint tokens - to address is hardcoded as per sample network
   ./mint.sh ${DOMAIN} ${NUMBER_OF_CALLS} ${AUTH_TOKEN}

   # Transfer tokens - addresses are hardcoded as per sample network
   ./transfer.sh ${DOMAIN} ${NUMBER_OF_CALLS} ${AUTH_TOKEN}

   # Burn tokens - address is hardcoded as per sample network   
   ./burn.sh ${DOMAIN} ${NUMBER_OF_CALLS} ${AUTH_TOKEN}
   ```

5. Enable Github SSO auth

   To enable Github's SSO based auth, add a secret to the application's namespace with clientid, clientsecret and redirect url.
   ```bash
   kubectl create secret generic -n mrshah githublogin --from-literal=clientid=<myclientid> --from-literal=clientsecret=<mysecret> --from-literal=redirecturl=https://<mydomain>/login/github/callback
   ```
   
## Future improvements:

1. Terraform variables

   Currently terraform is hardcoded with all parameters. Update to have variables.tf file and make it configurable.

2. Context setup

   Currently the app uses `context.TODO()` wherever a context is needed. A proper context should be setup and passed on.

3. Hardcoded gas

   Currently the app uses hardcoded gas value for each transaction. We should flow that in from API body.

4. Infinite loop on `nonce too low`

   If there is an `nonce too low` error on submitting transaction, a recurring call will happen to the function itself with incremented nonce. There is no limit on how many times it will try. We should limit and also add a delay between calls.

5. Tests

   As most other developers, I am a little behind on writing tests for this project. I generally like to keep up with the tests for my projects.
   
   I have written ginkgo based tests for a few modules. Some system tests have been written using shell scripts, that was the quickest way I could send transactions in parallel.

   The project is missing system tests completely. Travis should be setup with secrets to access kubernetes cluster and deploy the code using terraform, run a few transactions and report failures.

6. Resources

   Kubernetes best practice is to run the containers with limits and requests set for the resources of the containers. As I was running this on a free cluster, I wanted to keep it open and let processes fill up resources as and when required. We should add limits and requests in deployment specs.

7. NodePorts

   If we are using ingresses, the services can be changed to ClusterIPs. As I was running on a free IBM Kubernetes Services cluster, I had to use NodePort type service to be able to access the app.

8. eth_calls

   I was not able to make `eth_calls` successfully to retrieve balance, etc. So I removed the commented out code from the repo. The error I got was `no contract at given address`. I made sure that the contract is deployed as well as the transaction is minted.

9. Extend to deploy any contract

   Currently a static contract is deployed and metadata is stored. The code has ability to deploy any contract in future. The contract can be deployed and the abi generated on the fly by exposing an API. 

10. Metrics

   The app is not instrumented to have metrics server or capture metrics.

# Challenges

While this space was new for me, I faced a few challenges with setting up and making the things work.

- Spent most of my time refactoring [kaleido-go](https://github.com/kaleido-io/kaleido-go) and putting an API layer on top of it. I liked the way the executer->worker was being used. I built out a very generic framework on top it.

  Finally abandoned using the native rpcCalls as it was draining a lot of my time trying to find out how to retrieve the values from `eth_calls`. This was a rabbit hole. I wasted more than 5 days worth of time on this one. The repository is at [mrshah-at-ibm/kaleido-go](https://github.com/mrshah-at-ibm/kaleido-go). Reach out to me if you need access to it.

- Second most of my time went into learning and debuggin about how the whole system works - network, contracts, solc, abigen, eth_* rpcs, and more.

- Some issues that I found while debugging:

   - Finding geth logs. Took me sometime to find out why the network logs on stdout were not moving. It is because the geth logs are stored in a file inside `/qdata/logs/geth.log` file. It will be good if they were redirected to stdout.

   - quorum-tools has [Makefile](https://github.com/kaleido-io/quorum-tools/blob/master/Makefile) with `docker-clean` as default target. Running `make` will delete all the images.

   - [Cleanup in setup.sh](https://github.com/kaleido-io/quorum-tools/blob/master/examples/setup.sh#L89) needs to be done as root. It will be good if it can run a container as root and delete the `tmp` folder.

   - Raft ports are [hardcoded](https://github.com/kaleido-io/quorum-tools/blob/master/boot/lib/boot.js#L104-L108) to 50400 in `boot.js` file in quorum tools. To expose an ephermal port, it needs to be in 30000-32767 range. Changing the node in args of geth command does not work.
