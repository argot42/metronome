package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/ebitengine/oto/v3"
)

//go:embed tick.raw
var data []byte

const (
	SampleRate = 44100
	Channels   = 1
	Format     = oto.FormatFloat32LE
)

func main() {
	if len(os.Args) < 1 {
		fmt.Fprintf(os.Stderr, "%s <bpm>\n", os.Args[0])
		return
	}

	bpm := 100

	if len(os.Args) > 1 {
		var err error

		bpm, err = strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid bpm: %s", err)
			return
		}
	}

	if err := metronome(bpm); err != nil {
		fmt.Fprintf(os.Stderr, "metronome failed: %s", err)
		return
	}
}

func metronome(bpm int) error {
	o := &oto.NewContextOptions{
		SampleRate:   SampleRate,
		ChannelCount: Channels,
		Format:       Format,
		BufferSize:   time.Millisecond * 12,
	}

	otoCtx, ready, err := oto.NewContext(o)
	if err != nil {
		return err
	}

	<-ready

	r := bytes.NewReader(data)
	player := otoCtx.NewPlayer(r)
	defer player.Close()

	interval := time.Minute / time.Duration(bpm)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		player.Seek(0, io.SeekStart)
		player.Play()
	}

	return nil
}
