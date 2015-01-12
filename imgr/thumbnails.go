// Copyright 2015 Geofrey Ernest a.k.a gernest, All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package imgr

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
)

const (
	THUMNAIL_SIZE = 200
)

type Thumbnails struct {
	ImageManager
}

func (t *Thumbnails) Process(width, height int) error {
	var w, h int
	w = width
	h = height
	if width == 0 || height == 0 {
		w = THUMNAIL_SIZE
		h = THUMNAIL_SIZE
	}
	if len(t.Images) == 0 {
		return errors.New("There is nothing to process")
	}
	logThis.Event("Start Processing Images from %s", t.Source)
	for _, img := range t.Images {
		err := createThumbnail(img, t.Destinalion, w, h)
		if err != nil {
			return err
		}
	}
	logThis.Success("***Finished Processing***")
	return nil
}

func (l *Thumbnails) CreateThumbnail(src, dest string, width int, height int) error {
	img, err := getImageDetails(src)
	if err != nil {
		return err
	}
	pic := newImage(img, src)
	return createThumbnail(pic, dest, width, height)
}
func createThumbnail(img *Image, dest string, width int, height int) error {
	info, _ := os.Stat(filepath.Dir(img.Path))
	_, err := os.Stat(dest)
	if os.IsNotExist(err) {
		_ = os.Mkdir(dest, info.Mode())
	}
	pic, err := imaging.Open(img.Path)
	if err != nil {
		return err
	}
	destImg := imaging.Thumbnail(pic, width, height, imaging.Lanczos)
	destName := img.Name + "_thumbnail" + img.Ext
	destPath := filepath.Join(dest, destName)
	err = imaging.Save(destImg, destPath)
	if err != nil {
		return err
	}

	return nil
}
