package run

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

type GithubBody struct {
	FrpcVersion string `json:"tag_name"`
}

func DownloadFrpc() {

	githubProxy := "download.fastgit.org"
	frpGithubApiUrl := "http://github.tanwen.net/api/repos/fatedier/frp/releases/latest"
	frpGithubReleaseUrl := "https://github.com/fatedier/frp/releases"
	frpFileSuffix := ".tar.gz"

	if SysMachine == "windows" {
		frpFileSuffix = ".zip"
		FrpcBin = path.Join(FrpcBasePath, "frpc.exe")
	}

	_, err := os.Stat(FrpcBin)
	if err != nil {
		resp, err := http.Get(frpGithubApiUrl)
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		var data GithubBody
		json.Unmarshal(body, &data)

		frpDownloadUrl := fmt.Sprintf("%s/download/%s/frp_%s_%s_%s%s",
			strings.Replace(frpGithubReleaseUrl, "github.com", githubProxy, 1),
			data.FrpcVersion, strings.TrimLeft(data.FrpcVersion, "v"),
			SysMachine, SysArch, frpFileSuffix)
		frpFile := path.Base(frpDownloadUrl)

		log.Printf("download %s", frpFile)
		req, err := http.NewRequest("GET", frpDownloadUrl, nil)
		if err != nil {
			log.Panicln(err)
		}

		r, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatalln(err, 123)
		}
		defer r.Body.Close()
		if r.StatusCode != http.StatusOK {
			log.Fatalf("当前系统架构 %s_%s 没有找到匹配的包，请手动下载 %s 解压frpc文件%s到目录下\n", SysMachine, SysArch, frpGithubReleaseUrl, FrpcBasePath)
		}
		if SysMachine == "windows" {
			unZip(r.Body, FrpcBin)
		} else {
			unTargz(r.Body, FrpcBin)
		}

	}

}

func unTargz(r io.Reader, target string) {
	gReader, err := gzip.NewReader(r)
	if err != nil {
		log.Fatal(err)
	}
	tarReader := tar.NewReader(gReader)

	for {
		cur, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalln(err)
		}
		switch cur.Typeflag {
		case tar.TypeReg:
			if path.Base(cur.Name) == "frpc" {
				f, err := os.Create(target)
				if err != nil {
					log.Fatalln(err)
				}
				defer f.Close()
				f.ReadFrom(tarReader)
				f.Chmod(cur.FileInfo().Mode())
			}
		}
	}
}

func unZip(r io.ReadCloser, target string) {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatalln(err)
	}

	zReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		log.Fatalln(err)
	}

	for _, file := range zReader.File {
		if path.Base(file.Name) == "frpc.exe" {
			rc, err := file.Open()
			if err != nil {
				log.Fatal(err)
			}
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR|os.O_TRUNC, file.Mode())
			if err != nil {
				log.Fatalln(err)
			}
			defer f.Close()
			f.ReadFrom(rc)
		}
	}

}
