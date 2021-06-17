package ocrfreeapi

import (
	"ConcurrentProg/util"
	ocr "github.com/ranghetto/go_ocr_space"
	"github.com/spf13/viper"
)

type parseWorker struct {
	exec ocr.Config
}

var worker *parseWorker

/**
Arabic = ara
Bulgarian = bul
Chinese(Simplified) = chs
Chinese(Traditional) = cht
Croatian = hrv
Czech = cze
Danish = dan
Dutch = dut
English = eng
Finnish = fin
French = fre
German = ger
Greek = gre
Hungarian = hun
Korean = kor
Italian = ita
Japanese = jpn
Polish = pol
Portuguese = por
Russian = rus
Slovenian = slv
Spanish = spa
Swedish = swe
Turkish = tur
*/
func init() {
	err := util.InitConfig("config", "json", "B:/study/ConcurrentProg/ocrfreeapi")
	if err != nil {
		panic(err)
	}
	apiKey := viper.Get("api_key").(string)
	language := viper.Get("language").(string)
	worker = &parseWorker{
		//setting up the configuration
		exec: ocr.InitConfig(apiKey, language),
	}
}

func NewParseWorker() *parseWorker {
	return worker
}

func (w *parseWorker) ParseFromLocal(imageFilePath string) (string, error) {
	result, err := w.exec.ParseFromLocal(imageFilePath)
	if err != nil {
		return "", err
	}
	//printing the just the parsed text
	return result.JustText(), nil
}

func (w *parseWorker) ParseFromUrl(imageUrl string) (string, error) {
	result, err := w.exec.ParseFromUrl(imageUrl)
	if err != nil {
		return "", err
	}
	//printing the just the parsed text
	return result.JustText(), nil
}

func (w *parseWorker) ParseFromBase64(base64String string) (string, error) {
	result, err := w.exec.ParseFromBase64(base64String)
	if err != nil {
		return "", err
	}
	//printing the just the parsed text
	return result.JustText(), nil
}
