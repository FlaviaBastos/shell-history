

Shell History


# OVERVIEW

The shell history is an essential resource for a system administrator as it not only provides a view of past activities, but can also be used as a cookbook, providing insights for future activities and be used for postmortem analysis.

Unfortunately the standard way the shell history works attaches the commands list to a terminal, session and/or individual host. It would be much more valuable to have access to one’s shell history from any terminal you are currently working from.

# GOALS

1. Implement a shell history which is accessible from any host within a given set of restrictions

2. Make the history easily searchable from the terminal or via a web interface

# SPECIFICATIONS

The system should implement four main features: 

1. The local redirector, which will capture the shell command, parse it and upload to a remote server

2. The remote daemon, which will receive a rpc call containing the interesting payload

3. The command line interface, which will allow the user to:

    1. Enable/disable sending history to the remote server

    2. Query history

    3. Delete / Purge history

4. The web interface, which will allow the user to:

    4. Query the history

    5. Delete / Purge history

## Local redirector

The most transparent way to use the local redirector seems to be via [bash-preexec](https://github.com/rcaloras/bash-preexec).

A simple example of .bashrc using it:

```
source ~/.bash-preexec.sh
source ~/go/src/github.com/ebastos/shell-history/precmd.sh
```

## Design / Technology choices

The main four pieces will be implemented in two different languages:

Go (golang) for the command line redirector and Python + Django for the remote daemon and web interface.

Communication between the redirector and the remote daemon shall be implemented via gRPC using protobufs.

### Rationale

* Go (golang) is a compiled, multi-platform language with an outstanding standard library, which will produce self-contained binaries for easy deployment and no dependency hell.

* Python is a popular, decently performing scripting language, also with an outstanding standard library. A huge number of 3rd party libraries is also available to be added if necessary

* Django is a web app framework built on Python, applying the DRY principle and lots of "free" features, like admin interface. Its MVC also makes it easy to interact directly with the models using APIs, which we will do from the remote daemon

* gRPC is a high performance, open-source universal RPC framework backed by Google, which excels in small calls and is highly efficient over the wire.

* Protocol Buffers (protobufs) are Google's language-neutral, platform-neutral, extensible mechanism for serializing structured data. Based on a model it’s possible to generate APIs for multiple languages so there is an easy and efficient way to exchange data

### Technology Choices out of scope

Due to Python and Django’s multi-platform and multi-database support both the operating system on the server and the backend database chosen are irrelevant and should be subject to the choice of each deployment leader.

### Open for investigation

Web Interface: Material Design?

## Django Backend

The following models will be necessary:

1. User

    1. Id (primary key)

    2. Username (str)

    3. Password (str)

    4. ??? (history management. Can delete/Archive)

2. Host

    5. Id (primary key)

    6. Hostname

3. Command

    7. Id (primary key)

    8. User (foreign key)

    9. Host (foreign key)

    10. Cwd (str)

    11. Timestamp (time)

    12. Command (str)

    13. ExitCode (int)

    14. Hash (str)

    15. SourceAddress (str)

# SECURITY AND PRIVACY

Shell history can be considered sensitive data as it’s common for commands to include IP addresses, hostnames, usernames and even passwords. Transferring the commands between the local host and the remote server must happen in an encrypted and authenticated way.

gRPC native support to TLS/SSL should easily address this demand. On the webserver side HTTPS is also a requirement. 

It is suggested that the remote daemon should have access restricted to known and trusted subnets, although having it exposed to the open internet is not a major security risk.

Encryption of data on rest with GPG was also considered, but dismissed for this first draft due to the added complexity for the initial release.

Some challenges still stand, like the authentication of the remote user when query data via command line. Solutions still pending investigation.

# ALTERNATIVE APPROACHES

Regardless of this particular project approach, some options to make the shell history more flexible or accessible already exist:

* [Bashhub](https://www.bashhub.com/)

* [Saving bash history in syslog](https://coderwall.com/p/anphha/save-bash-history-in-syslog-on-centos)

* Local function (see references section)


# MILESTONES (Phase 1)

## Local redirector

* ~~Find the best way to capture the last command line~~

* ~~Find how to parse it and store inside a protobuf~~

* Find how to securely connect to gRPC and transmit data

## Remote daemon

* Find how to to listen and receive gRPC connection securely

* ~~Find how to receive and parse protobuf~~

## Django 

* Create models

**Protobuf**

* ~~Define it’s properties~~

# REFERENCES

[You Should Be Logging Shell History](https://www.jefftk.com/p/you-should-be-logging-shell-history)

[Secure gRPC with TLS/SSL](https://bbengfort.github.io/programmer/2017/03/03/secure-grpc.html)

