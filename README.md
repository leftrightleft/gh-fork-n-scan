# gh-fork-n-scan

A gh cli extension that forks all the repos from a source org into a destination org so you can run code scanning and secret scanning on them.

A couple notes:
* fork-n-scan will only copy over repos that are sourced in the source org.  Meaning, it will not fork forks of repos
* It will only fork public repos
* It will only fork non-archived repos
* Once the repos are forked into your destination org, fork-n-scan will add a topic to those forked repos.  That topic is the source organization.  This helps you identify the repos you care about in your destination org.

## Installation

This is an extension to the [gh cli](https://cli.github.com/).  Make sure you have that installed first.

Once the gh cli is installed, you can install this extension.  

From your local terminal:
```
gh extension install leftrightleft/gh-fork-n-scan
```

## Using fork-n-scan

To fork a batch of repos into your destination, you will execute the following commands:

```
gh fork-n-scan -s octodemo -d ghas-results
```

**Options**
```
gh fork-n-scan --help 
  -d string
        destination organization name
  -h    show help
  -r string
        source repository name
  -s string
        source organization name
```
