package metadata_api

/*
Integrate with SF Metadata API using native SF binaries.
When contracting instance, pass path to SF binary.

# Example

## Using version from PATH
```go
// SF binary will be taken from PATH
CmdMetadataClient{binaryPath: "sf"}
```

## Using from exact location
```
// SF from exact location will be used
CmdMetadataClient{binaryPath: "/home/myUser/bin/sf"}
```
*/
type CmdMetadataClient struct {
	BinaryPath string
}
