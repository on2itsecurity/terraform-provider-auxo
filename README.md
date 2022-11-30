# Terraform Provider Auxo

*Work In Progress*

However, the protectsurface, location and state resources are there and should work!

Supported resources:
* location
* protectsurface
* protectsurface_state

The provider can be found here: https://registry.terraform.io/providers/on2itsecurity/auxo/latest

## Usage

### First resource creation

To prevent your token from being exposed, you can use the `AUXOTOKEN` environment variable to pass the token to Terraform.
This way it doesn't need to be set in the terraform file(s).

```shell
export AUXOTOKEN="YOURSECRETTOKEN"
```

Now you can create a `main.tf` file and add the provider configuration.

```hcl
terraform {
  required_providers {
    auxo = {
      version = "0.0.1"
      source  = "on2itsecurity/auxo"
    }
  }
}
```

Add a resource to your `main.tf` file, i.e. a location.

```hcl
resource "auxo_location" "loc_zaltbommel" {
  name      = "Datacenter Zaltbommel"
  latitude  = 51.7983645
  longitude = 5.2548381
}
```

Initialize the workspace.

```shell
terraform init
```

Now you can create the plan and apply it.

```shell
terraform plan --out myfirstdeployment.plan
terraform apply myfirstdeployment.plan
```

If you want to destroy the resource, you can do so with the following command.

```shell
terraform destroy
```

## Development

### Install the provider for development

You can install the provider (locally) by cloning this repo and running `make install`.

```shell
git clone
git mod vendor
make install
```

This will make the provider avaialble to your local Terraform installation.

### Using the locally build provider on a different machine
Below example assumes a Linux system.

* Build or get the binary for the destination machine OS and architecture.

```shell
GOOS=linux GOARCH=amd64 go build -o terraform-provider-auxo
```

* Copy the binary to the destination server.

```shell
scp ./terraform-provider-auxo <USER>@<DESTINATION.SERVER>:~/
```

* Create a folder on the destination server in your home-directory and move the binary

```shell
mkdir -p ~/.terraform.d/plugins/on2itsecurity/auxo/0.1/linux_amd64
mv ~/terraform-provider-auxo ~/.terraform.d/plugins/on2itsecurity/auxo/0.1/linux_amd64/
```
