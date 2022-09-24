# Peernet Web Gateway

This Web Gateway is a regular Peernet client that provides a web gateway to Peernet. It allows to access files in Peernet through an open website. Accessing files via this web gateway is useful for regular web 2 users. For example, one can share a file on Peernet and then use this web gateway to make it accessible on Twitter.

## Compile

Download the [latest version of Go](https://golang.org/dl/). To build:

```
go build
```

## Deploy

Todo.

## Config

The config filename is hard-coded to `Config.yaml` and is created on the first run. Please see the [core library](https://github.com/PeernetOfficial/core#configuration) for individual settings to change.

Todo.

## Static Pages

Access to `/` will show a generic information page.

The main functionality is to provide URLs in the format `/[blockchain public key]/[file hash]` and other variations (using a directory name instead of file hash). This page shall provide access to one or multiple files shared via a user's blockchain. Access shall include download, preview, and providing a native Peernet link. The blockchain public key may be the node ID or the peer ID. Currently only hex encoding is supported.

Native Peernet linsk should start with the `peernet://` custom scheme. On Windows such links can be opened with a registered protocol handler.

To provide the best user experience, those static pages shall be generated immediately and use JavaScript to dynamically load the requested resource using the embedded API.

## API

Todo.

## Cache

Todo. A cache folder should temporarily cache requested files so that they are available faster.

