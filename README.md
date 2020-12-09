# jfrog-support-plugin

## About this plugin
This plugin is a template and a functioning example for a basic JFrog CLI plugin. 
This README shows the expected structure of your plugin's README.

## Installation with JFrog CLI
Since this plugin is currently not included in [JFrog CLI Plugins Registry](https://github.com/jfrog/jfrog-cli-plugins-reg), it needs to be built and installed manually. Follow these steps to install and use this plugin with JFrog CLI.
1. Make sure JFrog CLI is installed on you machine by running ```jfrog```. If it is not installed, [install](https://jfrog.com/getcli/) it.
2. Create a directory named ```plugins``` under ```~/.jfrog/``` if it does not exist already.
3. Clone this repository.
4. CD into the root directory of the cloned project.
5. Run ```go build``` to create the binary in the current directory.
6. Copy the binary into the ```~/.jfrog/plugins``` directory.

## Usage
### Commands
#### Name:
  jfrog support generate - Generates support bundle to supportlogs.

Usage:
  jfrog support generate [command options]

Options:
  --server-id          [Optional] Artifactory server ID configured using the config command.
  --send-to-support    [Default: false] Rather to upload the support bundle to JFrog support or not
  --ticket             [Optional] Ticket identifier for JFrog support team - must be provided when send-to-support = true
  --name               [Optional] Support bundle name - when empty will be auto generated
  --description        [Optional] Support bundle description
  --config             [Default: true] Include service configuration
  --system             [Default: true] Include service system information
  --logs               [Default: false] Include logs
  --dumps              [Default: false] Include thread dumps
  --dumps-count        [Optional] number of times to collect thread dump. Default:1
  --dumps-interval     [Optional] Interval between times of collection in milliseconds. Default:0
  --start              [Optional] start date from which to fetch the logs. pattern: YYYY-MM-DD
  --end                [Optional] end date until which to fetch the logs. pattern: YYYY-MM-DD
  
Environment Variables:
  SUPPORT_LOGS_URL
    [Default: https://supportlogs.jfrog.com/logs]
    Support logs base url - mostly for debug


#### Name:
  jfrog support upload - Uploads support bundle to supportlogs.

### Usage:
```
jfrog support upload <filepath> <ticket>

Arguments:
  filepath
    Bundle path on the local file system

  ticket
    Ticket number


Environment Variables:
  SUPPORT_LOGS_URL
    [Default: https://supportlogs.jfrog.com/logs]
    Support logs base url - mostly for debug
```

#### Name:
  jfrog support encrypt - Encrypt secret using Artifactory master key. Output will be in the form '<kid>.aesgcm256.<encrypted message>' or '<kid>.aesgcm128.<encrypted message>' depends on the key length

Usage:
```
jfrog support encrypt <plaintext> <key>

Arguments:
  plaintext
    Plain text to encrypt

  key
    Artifactory master key
```

#### Name:
  jfrog support decrypt - Decrypt secret using Artifactory master key. Currently supports only encrypted messages of the form: '<kid>.aesgcm256.<encrypted message>' or '<kid>.aesgcm128.<encrypted message>'

Usage:
```
jfrog support decrypt <secret> <key>

Arguments:
  secret
    The secret to decrypt

  key
    Artifactory master key

```