package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/martinlindhe/imgcat/lib"
	"github.com/yaegaki/itunes-app-interface"
)

func main() {
	err := NowPlaying()
	if err != nil {
		log.Fatal(err)
	}
}

func NowPlaying() error {
	err := itunes.Init()
	if err != nil {
		return err
	}

	defer itunes.UnInit()
	it, err := itunes.CreateItunes()
	if err != nil {
		return err
	}
	defer it.Close()

	t, err := it.CurrentTrack()
	if err != nil {
		return errors.New("Does not play song.")
	}
	defer t.Close()

	fmt.Printf("NowPlaying:%v %v\n", t.Name, t.Artist)

	artworks, err := t.GetArtworks()
	if err != nil {
		return err
	}

	artwork := <-artworks
	if artwork != nil {
		defer artwork.Close()
		dir, err := ioutil.TempDir("", "nowplaying")
		if err != nil {
			return err
		}
		defer os.RemoveAll(dir)

		path, err := artwork.SaveToFile(dir, "nowplaying")
		if err != nil {
			return err
		}

		imgcat.CatFile(path, os.Stdout)
	}

	return nil
}
