// The content of this file is heavily inspired by the file magic.go from
// https://github.com/perkeep/perkeep/blob/master/internal/magic/magic.go
// some of the fucntions and structs are modified to fit my needs
//
// Please refer to Perkeep's license at https://github.com/perkeep/perkeep/blob/master/COPYING

package mime

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/alinz/pkg.go/stream"
)

type Match struct {
	Offset     int
	Prefix     []byte
	Type       string
	Categories []string
}

// Returns the string format of Match object
// it mainly uses for detecting duplicate entry
func (m Match) String() string {
	return fmt.Sprintf("%d:%x:%s", m.Offset, m.Prefix, m.Type)
}

type MIME struct {
	macthes []Match
}

func (m *MIME) RegisterImageFormats() {
	m.Register(
		Match{Prefix: []byte("GIF87a"), Type: "image/gif", Categories: categories("media", "image", "gif")},
		Match{Prefix: []byte("GIF89a"), Type: "image/gif", Categories: categories("media", "image", "gif")},
		Match{Prefix: []byte("\xff\xd8\xff\xe2"), Type: "image/jpeg", Categories: categories("media", "image", "jpeg", "jpg")},
		Match{Prefix: []byte("\xff\xd8\xff\xe1"), Type: "image/jpeg", Categories: categories("media", "image", "jpeg", "jpg")},
		Match{Prefix: []byte("\xff\xd8\xff\xe0"), Type: "image/jpeg", Categories: categories("media", "image", "jpeg", "jpg")},
		Match{Prefix: []byte("\xff\xd8\xff\xdb"), Type: "image/jpeg", Categories: categories("media", "image", "jpeg", "jpg")},
		Match{Prefix: []byte("\x49\x49\x2a\x00\x10\x00\x00\x00\x43\x52\x02"), Type: "image/cr2", Categories: categories("media", "image", "cr2")},
		Match{Prefix: []byte{137, 'P', 'N', 'G', '\r', '\n', 26, 10}, Type: "image/png", Categories: categories("media", "image", "png")},
		Match{Prefix: []byte{0x49, 0x20, 0x49}, Type: "image/tiff", Categories: categories("media", "image", "tiff")},
		Match{Prefix: []byte{0x49, 0x49, 0x2A, 0}, Type: "image/tiff", Categories: categories("media", "image", "tiff")},
		Match{Prefix: []byte{0x4D, 0x4D, 0, 0x2A}, Type: "image/tiff", Categories: categories("media", "image", "tiff")},
		Match{Prefix: []byte{0x4D, 0x4D, 0, 0x2B}, Type: "image/tiff", Categories: categories("media", "image", "tiff")},
		Match{Prefix: []byte("8BPS"), Type: "image/vnd.adobe.photoshop", Categories: categories("media", "image", "photoshop", "ps")},
		Match{Prefix: []byte("gimp xcf "), Type: "image/x-xcf", Categories: categories("media", "image", "xcf")},
		Match{Prefix: []byte("II\x1a\000\000\000HEAPCCDR"), Type: "image/x-canon-crw", Categories: categories("media", "image", "canon", "crw")},   // Canon CIFF raw image data
		Match{Prefix: []byte("II\x2a\000\x10\000\000\000CR"), Type: "image/x-canon-cr2", Categories: categories("media", "image", "canon", "cr2")}, // Canon CR2 raw image data
		Match{Prefix: []byte("MMOR"), Type: "image/x-olympus-orf", Categories: categories("media", "image", "olympus", "orf")},                     // Olympus ORF raw image data, big-endian
		Match{Prefix: []byte("IIRO"), Type: "image/x-olympus-orf", Categories: categories("media", "image", "olympus", "orf")},                     // Olympus ORF raw image data, little-endian
		Match{Prefix: []byte("IIRS"), Type: "image/x-olympus-orf", Categories: categories("media", "image", "olympus", "orf")},                     // Olympus ORF raw image data, little-endian
		Match{Offset: 12, Prefix: []byte("DJVM"), Type: "image/vnd.djvu", Categories: categories("media", "image", "djvu")},                        // DjVu multiple page document
		Match{Offset: 12, Prefix: []byte("DJVU"), Type: "image/vnd.djvu", Categories: categories("media", "image", "djvu")},                        // DjVu image or single page document
		Match{Offset: 12, Prefix: []byte("DJVI"), Type: "image/vnd.djvu", Categories: categories("media", "image", "djvu")},                        // DjVu shared document
		Match{Offset: 12, Prefix: []byte("THUM"), Type: "image/vnd.djvu", Categories: categories("media", "image", "djvu")},                        // DjVu page thumbnails
	)
}

func (m *MIME) RegisterAudioFormats() {
	m.Register(
		Match{Prefix: []byte("fLaC\x00\x00\x00"), Type: "audio/x-flac", Categories: categories("media", "audio")},
		Match{Prefix: []byte{'I', 'D', '3'}, Type: "audio/mpeg", Categories: categories("media", "audio", "mpeg")},
		Match{Prefix: []byte("MThd"), Type: "audio/midi", Categories: categories("media", "audio", "midi")},              // Standard MIDI data
		Match{Prefix: []byte("MAC\040"), Type: "audio/ape", Categories: categories("media", "audio", "ape")},             // Monkey's Audio compressed format
		Match{Prefix: []byte("MP+"), Type: "audio/musepack", Categories: categories("media", "audio", "musepack")},       // Musepack audio
		Match{Offset: 8, Prefix: []byte("WAVE"), Type: "audio/x-wav", Categories: categories("media", "audio", "wav")},   // WAVE audio
		Match{Offset: 8, Prefix: []byte("AIFF"), Type: "audio/x-aiff", Categories: categories("media", "audio", "aiff")}, // AIFF audio
		Match{Offset: 8, Prefix: []byte("AIFC"), Type: "audio/x-aiff", Categories: categories("media", "audio", "aiff")}, // AIFF-C compressed audio
		Match{Offset: 8, Prefix: []byte("8SVX"), Type: "audio/x-aiff", Categories: categories("media", "audio", "aiff")}, // 8SVX 8-bit sampled sound voice
	)
}

func (m *MIME) RegisterVideoFormats() {
	m.Register(
		// Definition data extracted automatically from the file utility source code.
		// See: http://darwinsys.com/file/ (version used: 5.19)
		Match{Prefix: []byte{0, 0, 1, 0xB7}, Type: "video/mpeg", Categories: categories("media", "video", "mpeg")},
		Match{Prefix: []byte{0, 0, 0, 0x14, 0x66, 0x74, 0x79, 0x70, 0x71, 0x74, 0x20, 0x20}, Type: "video/quicktime", Categories: categories("media", "video", "quicktime")},
		Match{Prefix: []byte{0x1A, 0x45, 0xDF, 0xA3}, Type: "video/webm", Categories: categories("media", "video", "webm")},
		Match{Prefix: []byte("FLV\x01"), Type: "application/vnd.adobe.flash.video", Categories: categories("media", "video", "flash")},
		Match{Offset: 4, Prefix: []byte("moov"), Type: "video/quicktime", Categories: categories("media", "video", "quicktime")},         // Apple QuickTime
		Match{Offset: 4, Prefix: []byte("mdat"), Type: "video/quicktime", Categories: categories("media", "video", "quicktime")},         // Apple QuickTime movie (unoptimized)
		Match{Offset: 8, Prefix: []byte("isom"), Type: "video/mp4", Categories: categories("media", "video", "mp4")},                     // MPEG v4 system, version 1
		Match{Offset: 8, Prefix: []byte("mp41"), Type: "video/mp4", Categories: categories("media", "video", "mp4")},                     // MPEG v4 system, version 1
		Match{Offset: 8, Prefix: []byte("mp42"), Type: "video/mp4", Categories: categories("media", "video", "mp4")},                     // MPEG v4 system, version 2
		Match{Offset: 8, Prefix: []byte("mmp4"), Type: "video/mp4", Categories: categories("media", "video", "mp4")},                     // MPEG v4 system, 3GPP Mobile
		Match{Offset: 8, Prefix: []byte("3ge"), Type: "video/3gpp", Categories: categories("media", "video", "3gpp")},                    // MPEG v4 system, 3GPP
		Match{Offset: 8, Prefix: []byte("3gg"), Type: "video/3gpp", Categories: categories("media", "video", "3gpp")},                    // MPEG v4 system, 3GPP
		Match{Offset: 8, Prefix: []byte("3gp"), Type: "video/3gpp", Categories: categories("media", "video", "3gpp")},                    // MPEG v4 system, 3GPP
		Match{Offset: 8, Prefix: []byte("3gs"), Type: "video/3gpp", Categories: categories("media", "video", "3gpp")},                    // MPEG v4 system, 3GPP
		Match{Offset: 8, Prefix: []byte("avc1"), Type: "video/3gpp", Categories: categories("media", "video", "3gpp")},                   // MPEG v4 system, 3GPP JVT AVC
		Match{Offset: 8, Prefix: []byte("3g2"), Type: "video/3gpp2", Categories: categories("media", "video", "3gpp2")},                  // MPEG v4 system, 3GPP2
		Match{Offset: 8, Prefix: []byte("AVI\040"), Type: "video/x-msvideo", Categories: categories("media", "video", "msvideo", "avi")}, // AVI
	)
}

func (m *MIME) RegisterApplicationFormats() {
	m.Register(
		Match{Prefix: []byte{0, 0x6E, 0x1E, 0xF0}, Type: "application/vnd.ms-powerpoint", Categories: categories("ppt")},
		Match{Prefix: []byte{0x1F, 0x8B, 0x08}, Type: "application/x-gzip", Categories: categories("archive", "gzip")},
		Match{Prefix: []byte{0x37, 0x7A, 0xBC, 0xAF, 0x27, 0x1C}, Type: "application/x-7z-compressed", Categories: categories("archive", "7z")},
		Match{Prefix: []byte("BZh"), Type: "application/x-bzip2", Categories: categories("archive", "bzip2")},
		Match{Prefix: []byte{0xFD, 0x37, 0x7A, 0x58, 0x5A, 0}, Type: "application/x-xz", Categories: categories("archive", "xz")},
		Match{Prefix: []byte{'P', 'K', 3, 4, 0x0A, 0, 2, 0}, Type: "application/epub+zip", Categories: categories("ebook", "zip", "document")},
		Match{Prefix: []byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1}, Type: "application/vnd.ms-word", Categories: categories("document")},
		Match{Prefix: []byte{'P', 'K', 3, 4, 0x0A, 0x14, 0, 6, 0}, Type: "application/vnd.openxmlformats-officedocument.custom-properties+xml", Categories: categories("document", "office")},
		Match{Prefix: []byte{'P', 'K', 3, 4}, Type: "application/zip", Categories: categories("archive", "zip")},
		Match{Prefix: []byte("%PDF"), Type: "application/pdf", Categories: categories("pdf", "document")},
		Match{Prefix: []byte(".RMF\000\000\000"), Type: "application/vnd.rn-realmedia", Categories: categories("media")},     // RealMedia file
		Match{Prefix: []byte("OggS"), Type: "application/ogg", Categories: categories("media", "audio", "ogg")},              // Ogg data
		Match{Prefix: []byte("\000\001\000\000\000"), Type: "application/x-font-ttf", Categories: categories("ttf", "font")}, // TrueType font data
		Match{Prefix: []byte("d8:announce"), Type: "application/x-bittorrent", Categories: categories("torrent")},            // BitTorrent file
	)
}

func (m *MIME) RegisterMiscFormats() {
	m.Register(
		Match{Prefix: []byte("-----BEGIN PGP PUBLIC KEY BLOCK---"), Type: "text/x-openpgp-public-key", Categories: categories("pgp", "text")},
		Match{Prefix: []byte("{rtf"), Type: "text/rtf1", Categories: categories("rtf", "text")},
		Match{Prefix: []byte("BEGIN:VCARD\x0D\x0A"), Type: "text/vcard", Categories: categories("vcard", "text")},
		Match{Prefix: []byte("Return-Path: "), Type: "message/rfc822", Categories: categories("email", "message")},
	)
}

func (m *MIME) Register(macthes ...Match) {
	// rebuild the set for each already registred match
	set := make(map[string]struct{})
	for _, m := range m.macthes {
		set[m.String()] = struct{}{}
	}

	for _, match := range macthes {
		if _, ok := set[match.String()]; ok {
			continue
		}
		m.macthes = append(m.macthes, match)
	}
}

// Check returns type and categories for the given header bytes
// the header bytes is around 1kb of beginning of the content
func (m *MIME) Check(hdr []byte) (string, []string) {
	hlen := len(hdr)
	for _, match := range m.macthes {
		plen := match.Offset + len(match.Prefix)
		if hlen > plen && bytes.Equal(hdr[match.Offset:plen], match.Prefix) {
			return match.Type, match.Categories
		}
	}

	t := http.DetectContentType(hdr)
	t = strings.Replace(t, "; charset=utf-8", "", 1)
	if t != "application/octet-stream" && t != "text/plain" {
		return t, []string{}
	}

	return "", []string{}
}

// CheckReader returns type and categories for the given reader, it also returns the original r
// without advancing it, This is useful for checking the content type of a reader before processing it
func (m *MIME) CheckReader(r io.Reader) (string, []string, io.Reader) {
	hdr, r := stream.Peek(r, 1024)
	mime, categories := m.Check(hdr)
	return mime, categories, r
}

// CheckReadCloser is the same as CheckReader but also returns the same ReadCloser
func (m *MIME) CheckReadCloser(rc io.ReadCloser) (string, []string, io.ReadCloser) {
	mime, categories, r := m.CheckReader(rc)
	return mime, categories, stream.NewCloser(r, rc)
}

func New() *MIME {
	return &MIME{
		macthes: make([]Match, 0),
	}
}

func categories(c ...string) []string {
	return c
}
