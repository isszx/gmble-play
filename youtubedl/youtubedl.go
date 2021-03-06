package youtubedl

import (
	"bytes"
	"fmt"
	"layeh.com/gumble/gumbleffmpeg"
	"log"
	"os/exec"
	"strings"
)

func IsWhiteListedURL(url string) bool {
	// ! Don't forget to end url with / !
	whiteListedURLS := []string{"https://www.youtube.com/", "https://music.youtube.com/", "https://youtu.be/", "https://soundcloud.com/"}
	for i := range whiteListedURLS {
		if strings.HasPrefix(url, whiteListedURLS[i]) {
			return true
		}
	}
	return false
}

func GetYtdlTitle(url string) string {
	ytdl := exec.Command("youtube-dl", "-e", url)
	var output bytes.Buffer
	ytdl.Stdout = &output
	err := ytdl.Run()
	if err != nil {
		log.Println("Youtube-DL failed to get title for", url)
		return "missed title"
	}
	return output.String()
}

func GetYtdlSource(url string) gumbleffmpeg.Source {
	// TODO: Enforce --no-playlist AND/OR Make special handling for playlists
	gumbleSource := gumbleffmpeg.SourceExec("youtube-dl", "-f", "bestaudio", "--rm-cache-dir", "-q", "-o", "-", url)
	fmt.Println(gumbleSource)
	return gumbleSource
}
