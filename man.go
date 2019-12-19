// Do not edit. This file is generated by running 'make generated'.
package main

import "strings"

const generated = `.\"t
.\" Automatically generated by Pandoc 2.7.3
.\"
.TH "CITOP" "1" "" "" "version <version>"
.hy
.SH NAME
.PP
\f[B]citop\f[R] \[en] Continuous Integration Table Of Pipelines
.SH SYNOPSIS
.PP
\f[C]citop [-r REPOSITORY | --repository REPOSITORY] [COMMIT]\f[R]
.PP
\f[C]citop -h | --help\f[R]
.PP
\f[C]citop --version\f[R]
.SH DESCRIPTION
.PP
citop monitors the CI pipelines associated to a specific commit of a git
repository.
.PP
citop currently integrates with the following online services.
Each of the service is one or both of the following:
.IP \[bu] 2
A \[lq]source provider\[rq] that is used to list the pipelines
associated to a given commit of an online repository
.IP \[bu] 2
A \[lq]CI provider\[rq] that is used to get detailed information about
CI builds
.PP
.TS
tab(@);
lw(14.6n) lw(8.8n) lw(7.8n) lw(27.2n).
T{
Service
T}@T{
Source
T}@T{
CI
T}@T{
URL
T}
_
T{
GitHub
T}@T{
yes
T}@T{
no
T}@T{
<https://github.com/>
T}
T{
GitLab
T}@T{
yes
T}@T{
yes
T}@T{
<https://gitlab.com/>
T}
T{
AppVeyor
T}@T{
no
T}@T{
yes
T}@T{
<https://www.appveyor.com/>
T}
T{
CircleCI
T}@T{
no
T}@T{
yes
T}@T{
<https://circleci.com/>
T}
T{
Travis CI
T}@T{
no
T}@T{
yes
T}@T{
<https://travis-ci.org/> <https://travis-ci.com/>
T}
T{
Azure Devops
T}@T{
no
T}@T{
yes
T}@T{
<https://dev.azure.com>
T}
.TE
.SH POSITIONAL ARGUMENTS
.SS \f[C]COMMIT\f[R]
.PP
Specify the commit to monitor.
COMMIT is expected to be the SHA identifier of a commit, or the name of
a tag or a branch.
If this option is missing citop will monitor the commit referenced by
HEAD.
.PP
Example:
.IP
.nf
\f[C]
# Show pipelines for commit 64be3c6
citop 64be3c6
# Show pipelines for the commit referenced by the tag \[aq]0.9.0\[aq]
citop 0.9.0
# Show pipelines for the commit at the tip of a branch
citop feature/doc
\f[R]
.fi
.SH OPTIONS
.SS \f[C]-r=REPOSITORY, --repository=REPOSITORY\f[R]
.PP
Specify the git repository to work with.
REPOSITORY can be either a path to a local git repository, or the URL of
an online repository hosted at GitHub or GitLab.
Both web URLs and git URLs are accepted.
.PP
In the absence of this option, citop will work with the git repository
located in the current directory.
If there is no such repository, citop will fail.
.PP
Examples:
.IP
.nf
\f[C]
# Work with the git repository in the current directory
citop
# Work with the repository specified by a web URL
citop -r https://gitlab.com/nbedos/citop
citop -r github.com/nbedos/citop
# Git URLs are accepted
citop -r git\[at]github.com:nbedos/citop.git
# Paths to a local repository are accepted too
citop -r /home/user/repos/myrepo
\f[R]
.fi
.SS \f[C]-h, --help\f[R]
.PP
Show usage of citop
.SS \f[C]--version\f[R]
.PP
Print the version of citop being run
.SH INTERACTIVE COMMANDS
.PP
Below are the default commands for interacting with citop.
.PP
.TS
tab(@);
lw(10.7n) lw(46.7n).
T{
Key
T}@T{
Action
T}
_
T{
Up, j
T}@T{
Move cursor up by one line
T}
T{
Down, k
T}@T{
Move cursor down by one line
T}
T{
Page Up
T}@T{
Move cursor up by one screen
T}
T{
Page Down
T}@T{
Move cursor down by one screen
T}
T{
o, +
T}@T{
Open the fold at the cursor
T}
T{
O
T}@T{
Open the fold at the cursor and all sub-folds
T}
T{
c, -
T}@T{
Close the fold at the cursor
T}
T{
C
T}@T{
Close the fold at the cursor and all sub-folds
T}
T{
/
T}@T{
Open search prompt
T}
T{
Escape
T}@T{
Close search prompt
T}
T{
Enter, n
T}@T{
Move to the next match
T}
T{
N
T}@T{
Move to the previous match
T}
T{
v
T}@T{
View the log of the job at the cursor[a]
T}
T{
b
T}@T{
Open with default web browser
T}
T{
q
T}@T{
Quit
T}
T{
?
T}@T{
View manual page
T}
.TE
.IP \[bu] 2
[a] Note that if the job is still running, the log may be incomplete.
.SH CONFIGURATION FILE
.SS Location
.PP
citop follows the XDG base directory specification [2] and expects to
find the configuration file at one of the following locations depending
on the value of the two environment variables \f[C]XDG_CONFIG_HOME\f[R]
and \f[C]XDG_CONFIG_DIRS\f[R]:
.IP "1." 3
\f[C]\[dq]$XDG_CONFIG_HOME/citop/citop.toml\[dq]\f[R]
.IP "2." 3
\f[C]\[dq]$DIR/citop/citop.toml\[dq]\f[R] for every directory
\f[C]DIR\f[R] in the comma-separated list
\f[C]\[dq]$XDG_CONFIG_DIRS\[dq]\f[R]
.PP
If \f[C]XDG_CONFIG_HOME\f[R] (resp.
\f[C]XDG_CONFIG_DIRS\f[R]) is not set, citop uses the default value
\f[C]\[dq]$HOME/.config\[dq]\f[R] (resp.
\f[C]\[dq]/etc/xdg\[dq]\f[R]) instead.
.SS Format
.PP
citop uses a configuration file in TOML version
v0.5.0 (https://github.com/toml-lang/toml/blob/master/versions/en/toml-v0.5.0.md)
format.
The configuration file is made of keys grouped together in tables.
The specification of each table is given in the example below.
.SS Example
.PP
This example describes and uses all existing configuration options.
.IP
.nf
\f[C]
#### CITOP CONFIGURATION FILE ####
# This file is a complete, valid configuration file for citop
# and should be located at $XDG_CONFIG_HOME/citop/citop.toml
# 

## PROVIDERS ##
[providers]
# The \[aq]providers\[aq] table is used to define credentials for 
# accessing online services. citop relies on two types of
# providers:
#
#    - \[aq]source providers\[aq] are used for listing the CI pipelines
#    associated to a given commit (GitHub and GitLab are source
#    providers)
#    - \[aq]CI providers\[aq] are used to get detailed information about
#    CI pipelines (GitLab, AppVeyor, CircleCI, Travis and Azure
#    Devops are CI providers)
#
# citop requires credentials for at least one source provider and
# one CI provider to run. Feel free to remove sections below 
# as long as this rule is met.
#
# Note that for all providers, not setting an API token or 
# setting \[ga]token = \[dq]\[dq]\[ga] will cause the provider to make
# unauthenticated API requests. 
#

### GITHUB ###
[[providers.github]]
# GitHub API token (optional, string)
#
# Note: Unauthenticated API requests are heavily rate-limited by 
# GitHub (60 requests per hour and per IP address) whereas 
# authenticated clients benefit from a rate of 5000 requests per
# hour. Providing an  API token is strongly encouraged: without
# one, citop will likely reach the rate limit in a matter of
# minutes.
#
# GitHub token management: https://github.com/settings/tokens
token = \[dq]\[dq]


### GITLAB ###
[[providers.gitlab]]
# Name shown by citop for this provider
# (optional, string, default: \[dq]gitlab\[dq])
name = \[dq]gitlab\[dq]

# GitLab instance URL (optional, string, default: \[dq]https://gitlab.com\[dq])
# (the GitLab instance must support GitLab REST API V4)
url = \[dq]https://gitlab.com\[dq]

# GitLab API token (optional, string)
#
# Note: GitLab prevents access to pipeline jobs for 
# unauthenticated users meaning if you wish to use citop
# to view GitLab pipelines you will have to provide
# appropriate credentials. This is true even for pipelines
# of public repositories.
#
# gitlab.com token management:
#     https://gitlab.com/profile/personal_access_tokens
token = \[dq]\[dq]


### TRAVIS CI ###
[[providers.travis]]
# Name shown by citop for this provider
# (optional, string, default: \[dq]travis\[dq])
name = \[dq]travis\[dq]

# URL of the Travis instance. \[dq]org\[dq] and \[dq]com\[dq] can be used as
# shorthands for the full URL of travis.org and travis.com
# (string, mandatory)
url = \[dq]org\[dq]

# API access token for the travis API (string, optional)
# Travis tokens are managed at:
#    - https://travis-ci.org/account/preferences
#    - https://travis-ci.com/account/preferences
token = \[dq]\[dq]


# Define another account for accessing travis.com
[[providers.travis]]
name = \[dq]travis\[dq]
url = \[dq]com\[dq]
token = \[dq]\[dq]


### APPVEYOR ###
[[providers.appveyor]]
# Name shown by citop for this provider
# (optional, string, default: \[dq]appveyor\[dq])
name = \[dq]appveyor\[dq]

# AppVeyor API token (optional, string)
# AppVeyor token managemement: https://ci.appveyor.com/api-keys
token = \[dq]\[dq]


### CIRCLECI ###
[[providers.circleci]]
# Name shown by citop for this provider
# (optional, string, default: \[dq]circleci\[dq])
name = \[dq]circleci\[dq]

# Circle CI API token (optional, string)
# See https://circleci.com/account/api
token = \[dq]\[dq]


### AZURE DEVOPS ###
[[providers.azure]]
# Name shown by citop for this provider
# (optional, string, default: \[dq]azure\[dq])
name = \[dq]azure\[dq]

# Azure API token (optional, string)
# Azure token management is done at https://dev.azure.com/ via
# the user settings menu
token = \[dq]\[dq]
\f[R]
.fi
.SH ENVIRONMENT
.SS ENVIRONMENT VARIABLES
.IP \[bu] 2
\f[C]BROWSER\f[R] is used to find the path of the default web browser
.IP \[bu] 2
\f[C]PAGER\f[R] is used to view log files.
If the variable is not set, citop will call \f[C]less\f[R]
.IP \[bu] 2
\f[C]HOME\f[R], \f[C]XDG_CONFIG_HOME\f[R] and \f[C]XDG_CONFIG_DIRS\f[R]
are used to locate the configuration file
.SS LOCAL PROGRAMS
.PP
citop relies on the following local executables:
.IP \[bu] 2
\f[C]git\f[R] to translate the abbreviated SHA identifier of a commit
into a non-abbreviated SHA
.IP \[bu] 2
\f[C]less\f[R] to view log files, unless \f[C]PAGER\f[R] is set
.IP \[bu] 2
\f[C]man\f[R] to show the manual page
.SH EXAMPLES
.PP
Monitor pipelines of the current git repository
.IP
.nf
\f[C]
# Move to a directory containing a git repository of your choosing
git clone git\[at]github.com:nbedos/citop.git && cd citop
# Run citop to list the pipelines associated to the last commit of the repository 
citop

# Show pipelines associated to a specific commit, tag or branch
citop a24840c
citop 0.1.0
citop master
\f[R]
.fi
.PP
Monitor pipelines of other repositories
.IP
.nf
\f[C]
# Show pipelines of a repository identified by a URL or path
citop -r https://gitlab.com/nbedos/citop        # Web URL
citop -r git\[at]github.com:nbedos/citop.git        # Git URL
citop -r github.com/nbedos/citop                # URL without scheme
citop -r /home/user/repos/repo                  # Path to a repository

# Specify both repository and git reference
citop -r github.com/nbedos/citop master
\f[R]
.fi
.SH BUGS
.PP
Questions, bug reports and feature requests are welcome and should be
submitted on GitHub (https://github.com/nbedos/citop/issues).
.SH NOTES
.IP "1." 3
\f[B]citop repository\f[R]
.RS 4
.IP \[bu] 2
<https://github.com/nbedos/citop>
.RE
.IP "2." 3
\f[B]XDG base directory specification\f[R]
.RS 4
.IP \[bu] 2
<https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html>
.RE
.SH AUTHORS
Nicolas Bedos.`

func manualPage() string {
	return strings.Replace(generated, "<version>", Version, 1)
}

