Terraform Websupportsk Provider
=============================

- Website: https://www.terraform.io

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="300px">

Maintainers
-----------

This provider plugin is maintained by:

* Filip Ilavsky ([@fossilmfc](https://github.com/fossilmfc))

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.15.x and higher recommended
-	[Go](https://golang.org/doc/install) 1.16 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/fossilmfc/terraform-provider-websupportsk`

```sh
$ mkdir -p $GOPATH/src/github.com/fossilmfc; cd $GOPATH/src/github.com/fossilmfc
$ git clone https://github.com/fossilmfc/terraform-provider-websupportsk
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/fossilmfc/terraform-provider-websupportsk
$ goreleaser
```

Using the provider
----------------------
Check documentation of the provider in terraform registry by searching for `fossilmfc/websupportsk` provider.

Developing the Provider
---------------------------

Provider is currently being build by goreleaser (you need to install it locally first). 
You can also try to build by running command `goreleaser`.