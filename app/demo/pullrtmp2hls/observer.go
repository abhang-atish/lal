package main

import (
	"encoding/json"
	"fmt"
	"github.com/grafov/m3u8"
	"github.com/q191201771/lal/pkg/base"
	"path/filepath"
	"strings"
	"time"
)

func PrettyPrint(data interface{}) {
	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(out))
}

type Observer struct {
}

func MakeManifestM3U8(m3u8Path, tsName string, duration float64) {
	manifestPath := strings.Replace(m3u8Path, "playlist", "manifest", -1)
	manifestFile := ManifestFile{
		Path:     manifestPath,
		Key:      "manifest_key",
		Name:     "manifest.m3u8",
		IsMaster: false,
	}

	manifestFile.Read(false)

	recreateList := func(playlist *m3u8.MediaPlaylist, isNew bool, isDownloaded bool) {
		nonNilPlaylist := GetNonNilSegments(playlist.Segments)
		if !isNew && len(nonNilPlaylist) >= 1 {
			if nonNilPlaylist[0].Discontinuity {
				playlist.DiscontinuitySeq += 1
			}
			_ = playlist.Remove()
		}

		tsName = fmt.Sprintf("http://localhost:54321/%v/", tsName)
		var x = m3u8.MediaSegment{Discontinuity: false, URI: tsName, Duration: 2}
		pdt, err := ConvertTimeToRFC3339Nano(time.Now())
		if err == nil {
			x.ProgramDateTime = pdt
		}

		err = playlist.AppendSegment(&x)
		if err != nil {
			panic(fmt.Sprintf("Add segment #%s to a media playlist failed: %s", manifestFile.Name, err))
		}
		manifestFile.Bytes = playlist.Encode().Bytes()
		manifestFile.SaveToDisk()
	}

	var isDownloaded = false

	if len(manifestFile.Bytes) == 0 {
		manifestFile.Bytes, _ = manifestFile.GetOldCopy()
		isDownloaded = true
	}

	if len(manifestFile.Bytes) == 0 {
		p, _ := m3u8.NewMediaPlaylist(2, 2)
		recreateList(p, true, isDownloaded)
		return
	}

	play, _, _ := ParseM3U8(manifestFile.Bytes, manifestFile.Key)
	playlist, ok := play.(*m3u8.MediaPlaylist)
	if !ok {
		return
	}
	recreateList(playlist, false, isDownloaded)
}

func (o *Observer) OnHlsMakeTs(info base.HlsMakeTsInfo) {
	if info.Event == "open" {
		MakeManifestM3U8(info.LiveM3u8File, filepath.Base(info.TsFile), info.Duration)
	}
}

func (o *Observer) OnFragmentOpen() {
	//fmt.Println("OnFragmentOpen")
}
