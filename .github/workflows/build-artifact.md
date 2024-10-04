# terraform-provider-passwordsafe artifact

The artifact made by the `release.yml` Github Actions workflow is a **Terraform Provider** available on the **Terraform Registry** and **Jfrog Artifactory**.

## Build and Publish the artifact

To build the Terraform Provider, the `release` workflow starts executing `go mod`, and then import a **GPG** key, that is needed in the further steps.

Next, the workflow uses **Go Releaser**, that is the tool recommended by Terraform to build and release the provider ready for its registry. **Go Releaser** automates the release for the provider, to either publish a draft release in the Github repository or to make a standalone zip file artifact (which is the one that is send to **JFrog Artifact**). The release is signed using the imported **GPG** key.

Then, the zipped artifact is moved inside a new folder called `terraform-provider-passwordsafe`, which is uploaded to the Artifactory repository `eng-generic-dev-local`.