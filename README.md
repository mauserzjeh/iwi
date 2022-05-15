![GitHub release (latest by date)](https://img.shields.io/github/v/release/mauserzjeh/iwi?style=flat-square)

# iwi
IWI (Infirnity Ward Image) processing library

# Supported games
- Call of Duty 2
- Call of Duty 4 Modern Warfare
- Call of Duty 5
- Call of Duty Modern Warfare 2
- Call of Duty Modern Warfare 3
- Call of Duty Black Ops
- Call of Duty Black Ops 2

# Installation
```
go get -u github.com/mauserzjeh/iwi
```

# Usage
```go
// import the library
import "github.com/mauserzjeh/iwi"

iwiFile, err := iwi.ReadIWI("path/to/iwi/file.iwi")

// check for errors
if err != nil {
    log.Fatal(err)
}

// iwiFile holds the data of the .iwi file
header := iwiFile.Header // access information about the header
info := iwiFile.Info // access information about the file
data := iwiFile.Data // access the image data that was read in (only the highest mip)
```

The library exposes various constants to check for version and format
```go
const (
	// IWi Versions
	IWI_VERSION_COD2   = 0x05 // CoD2
	IWI_VERSION_COD4   = 0x06 // CoD4
	IWI_VERSION_COD5   = 0x06 // CoD5
	IWI_VERSION_CODMW2 = 0x08 // CoDMW2
	IWI_VERSION_CODMW3 = 0x08 // CoDMW3
	IWI_VERSION_CODBO1 = 0x0D // CoDBO1
	IWI_VERSION_CODBO2 = 0x1B // CoDBO2

	// IWi Format
	IWI_FORMAT_ARGB32 = 0x01 // ARGB32
	IWI_FORMAT_RGB24  = 0x02 // RGB24
	IWI_FORMAT_GA16   = 0x03 // GA16
	IWI_FORMAT_A8     = 0x04 // A8
	IWI_FORMAT_DXT1   = 0x0B // DXT1
	IWI_FORMAT_DXT3   = 0x0C // DXT3
	IWI_FORMAT_DXT5   = 0x0D // DXT5
)
```
