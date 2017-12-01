### Release Process Workflow

#### Pre-Requisites

This workflow should be zero-touch and lead to regular, stable, releases. Releases should satisfy the following conditions:

- Passes build and compile time error checking for all platforms
- Passes tests
  - 100% coverage of test functions
  - Tests test for most common theoretical conditions
- Passes linting using all supported linters (using `gometalinter`)
- Contains substantial or necessary changes from previous releases
- Contains a `CHANGELOG.MD` entry with changes listed
- All necessary documentation created or updated
- Signed off on and acceptable for general use (except in "pre-release" conditions)
- Contains all elements
  - `ksync` CLI binary
  - `radar` server side binary
  - Docker container image for deployment

#### Process
Releases will be handled via CI and be released using the GitHub release mechanism. Initiating a release (after the [pre-requisites](#pre-requisites) have been met) should involve a single manual step, pushing a tag to the repository.

###### Steps
1. A new tag is created for a specific commit and pushed to GitHub
2. Push triggers CI
  1. `master` branch is pulled
  2. Standard CI flow is run
  3. If passing, tag creation triggers secondary flow
3. github.com/tcnksm/ghr is used to create a prerelease setup from tag and build artifacts
4. Pre-release object is pushed to GitHub
5. (Optional) Pre-release is manually verified and promoted to full release.
