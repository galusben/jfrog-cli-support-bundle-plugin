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

## NAME:
   support - Perform support operations like creating and uploading support bundles, encrypt decrypt passwords with master key etc.

## USAGE:
   support [global options] command [command options] [arguments...]
   
## VERSION:
   v0.0.1
   
## COMMANDS:
   upload, up    Uploads support bundle to supportlogs.
   generate, up  Generates support bundle to supportlogs.
   decrypt, up   Decrypt secret using Artifactory master key. Currently supports only encrypted messages of the form: '<kid>.aesgcm256.<encrypted message>' or '<kid>.aesgcm128.<encrypted message>'
   encrypt, up   Encrypt secret using Artifactory master key. Output will be in the form '<kid>.aesgcm256.<encrypted message>' or '<kid>.aesgcm128.<encrypted message>' depends on the key length
   help, h       Shows a list of commands or help for one command
   
## GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
   
