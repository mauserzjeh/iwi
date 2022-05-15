package iwi

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"sort"
)

type (
	IWI struct {
		Header header
		Info   info
		Data   []byte
	}

	header struct {
		Magic   [3]byte
		Version byte
	}

	info struct {
		Format byte
		Usage  byte
		Width  uint16
		Height uint16
		Depth  uint16
	}

	offsets []int32
	mipmap  struct {
		offset int32
		size   int32
	}
	mipmaps []mipmap
)

const (
	// IWI Versions
	IWI_VERSION_COD2   = 0x05 // CoD2
	IWI_VERSION_COD4   = 0x06 // CoD4
	IWI_VERSION_COD5   = 0x06 // CoD5
	IWI_VERSION_CODMW2 = 0x08 // CoDMW2
	IWI_VERSION_CODMW3 = 0x08 // CoDMW3
	IWI_VERSION_CODBO1 = 0x0D // CoDBO1
	IWI_VERSION_CODBO2 = 0x1B // CoDBO2

	// IWI Format
	IWI_FORMAT_ARGB32 = 0x01 // ARGB32
	IWI_FORMAT_RGB24  = 0x02 // RGB24
	IWI_FORMAT_GA16   = 0x03 // GA16
	IWI_FORMAT_A8     = 0x04 // A8
	IWI_FORMAT_DXT1   = 0x0B // DXT1
	IWI_FORMAT_DXT3   = 0x0C // DXT3
	IWI_FORMAT_DXT5   = 0x0D // DXT5
)

// supportedVersions
func supportedIWiVersions() []uint8 {
	return []uint8{
		IWI_VERSION_COD2,
		IWI_VERSION_COD4,
		IWI_VERSION_COD5,
		IWI_VERSION_CODMW2,
		IWI_VERSION_CODMW3,
		IWI_VERSION_CODBO1,
		IWI_VERSION_CODBO2,
	}
}

// isSupported
func (h header) isSupported() bool {
	if !(h.Magic[0] == 'I' && h.Magic[1] == 'W' && h.Magic[2] == 'i') {
		return false
	}

	for _, sv := range supportedIWiVersions() {
		if sv == h.Version {
			return true
		}
	}

	return false
}

// mipMaps calculate mipmap offsets and sizes from the given offsets
func (o offsets) mipMaps(first int32, size int32) mipmaps {
	m := make(mipmaps, len(o))

	for i := 0; i < len(o); i++ {
		switch i {
		case 0:
			m = append(m, mipmap{
				offset: o[i],
				size:   size - o[i],
			})

		case len(o) - 1:
			m = append(m, mipmap{
				offset: first,
				size:   o[i] - first,
			})

		default:
			m = append(m, mipmap{
				offset: o[i],
				size:   o[i-1] - o[i],
			})
		}
	}

	return m
}

func (m mipmaps) Len() int           { return len(m) }                // Sort interface - Len
func (m mipmaps) Less(i, j int) bool { return m[i].size < m[j].size } // Sort interface-  Less
func (m mipmaps) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }      // Sort interface - Swap

// ReadIWi reads an IWI file and returns a pointer to
// an IWI structure holding its data
func ReadIWI(file string) (*IWI, error) {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewReader(f)

	// read header
	var h header
	err = binary.Read(buf, binary.LittleEndian, &h)
	if err != nil {
		return nil, err
	}

	if !h.isSupported() {
		return nil, fmt.Errorf("unsupported version: %s%v", h.Magic, h.Version)
	}

	if h.Version == IWI_VERSION_CODMW2 || h.Version == IWI_VERSION_CODMW3 {
		_, err := buf.Seek(0x8, io.SeekStart)
		if err != nil {
			return nil, err
		}
	}

	var i info
	err = binary.Read(buf, binary.LittleEndian, &i)
	if err != nil {
		return nil, err
	}

	ofs := make(offsets, 4)
	if h.Version == IWI_VERSION_CODBO1 || h.Version == IWI_VERSION_CODBO2 {
		ofs = make(offsets, 8)
		_, err := buf.Seek(0x10, io.SeekStart)
		if err != nil {
			return nil, err
		}

		if h.Version == IWI_VERSION_CODBO2 {
			_, err := buf.Seek(0x20, io.SeekStart)
			if err != nil {
				return nil, err
			}
		}
	}

	err = binary.Read(buf, binary.LittleEndian, &ofs)
	if err != nil {
		return nil, err
	}

	curr, err := buf.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}
	mms := ofs.mipMaps(int32(curr), int32(buf.Size()))
	sort.Sort(sort.Reverse(mms))
	mm := mms[0]

	_, err = buf.Seek(int64(mm.offset), io.SeekStart)
	if err != nil {
		return nil, err
	}

	data := make([]byte, mm.size)
	err = binary.Read(buf, binary.LittleEndian, &data)
	if err != nil {
		return nil, err
	}

	return &IWI{
		Header: h,
		Info:   i,
		Data:   data,
	}, nil
}
