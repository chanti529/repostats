# repostats

## About this plugin
This plugin utilizes the statistical metadata of the Artifacts that is implicitly generated and maintained by Artifactory. These statisctics provide insights into Artifactory utilization and can help organizations to implement checks, chargeback Users/Teams, perform cleanups and analyse the consumption of resources.

This can help Users to find out the most poularly downlaoded Artifacts in a given repository, Artifacts that are consuming the most space in a given Repository with various levels of customization avaialble. Results obtained can also be sorted and filtered. 

## Installation with JFrog CLI
Since this plugin is currently not included in [JFrog CLI Plugins Registry](https://github.com/jfrog/jfrog-cli-plugins-reg), it needs to be built and installed manually. Follow these steps to install and use this plugin with JFrog CLI.
1. Make sure JFrog CLI is installed on you machine by running ```jfrog```. If it is not installed, [install](https://jfrog.com/getcli/) it.
2. Create a directory named ```plugins``` under ```~/.jfrog/``` if it does not exist already.
3. Clone this repository.
4. CD into the root directory of the cloned project.
5. Run ```go build``` to create the binary in the current directory.
6. Copy the binary into the ```~/.jfrog/plugins``` directory.


## download command - Download Count statistics
This command provides the download count statistics on a Repository/Folder path/Artifact level with an option to filter the results based on the properties and last downloaded timestamp or interval. It also provides the option to select the server id of the interested Artifactory instance.


## size command - Repository size statistics in Bytes
This command gives the info about the size of an Artifact/Folder/Repository in Bytes with an option to filter the results based on the properties and last modified timestamp or interval. It also provides the option to select the server id of the interested Artifactory instance.



## Usage
### Commands
* download - Get repo download count statistics.
    - Arguments:
        - type - Type of component to get statistics. Valid values: artifact, folder, repos, user
    - Options:
        - repos: [Mandatory] Comma separated list of repositories.
        - path: [Optional] Regular Expression to filter the full path of artifacts.
        - properties: [Optional] Comma separeted list of properties and values to filter in the format property_name=pattern
        - server-id: [Optional] Artifactory server ID configured using the config command.
        - lastdownloadedfrom: [Optional] Filter artifacts last downloaded after given timestamp in RFC3339 format.
        - lastdownloadedto: [Optional] Filter artifacts last downloaded before given timestamp in RFC3339 format.
        - limit: [Default: 10] Max number or results. Set value to 0 to disable limit
        - sort: [Default: desc] Results order. Valid values: desc, asc, alpha
        - page-size: [Default: 50000] Number of items to be processed at once per a single worker
        - max-workers: [Default: 5] Max number of concurrent workers processing items in parallel at a given time
        - max-depth: [Default: 4] Max depth to group folders when using folder command type
    - Examples:
    ```
    $ jfrog repostats download artifact --repos jcenter-remote --path .+.jar
    
    $ jfrog repostats downlaod folder --repos jcenter-remote --path .+.jar --lastdownloadedto 2020-05-12T15:55:00Z --limit 0 --max-depth 2
    
    $ jfrog repostats download user --repos jcenter-remote --path .+.jar --lastdownloadedto 2020-05-12T15:55:00Z --limit 0

    ```

* size - Get repo size statistics in Bytes.
    - Arguments:
        - type - Type of component to get statistics. Valid values: artifact, folder, repos, user
    - Options:
        - repos:        [Mandatory] Comma separated list of repositories.
        - path:         [Optional] Regular Expression to filter the full path of artifacts.
        - properties:   [Optional] Comma separeted list of properties and values to filter in the format property_name=pattern
        - server-id:    [Optional] Artifactory server ID configured using the config command.
        - modifiedfrom: [Optional] Filter artifacts modified after given timestamp in format RFC3339.
        - modifiedto:   [Optional] Filter artifacts modified before given timestamp in format RFC3339.
        - limit:        [Default: 10] Max number or results. Set value to 0 to disable limit
        - sort:         [Default: desc] Results order. Valid values: desc, asc, alpha
        - page-size:    [Default: 50000] Number of items to be processed at once per a single worker
        - max-workers:  [Default: 5] Max number of concurrent workers processing items in parallel at a given time
        - max-depth:    [Default: 4] Max depth to group folders when using folder command type
    - Examples:
    ```
    $ jfrog repostats size artifact --repos maven-local --path .+.jar

    ```

### Environment variables
None.

## Additional info
None.

## Release Notes
The release notes are available [here](RELEASE.md).
