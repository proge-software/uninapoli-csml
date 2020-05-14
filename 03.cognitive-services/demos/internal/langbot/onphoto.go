package langbot

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/proge-software/uninapoli-csml-csbot/pkg/wssface"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (b *Bot) onPhoto(m *tb.Message) {
	rc, err := b.tbot.GetFile(&m.Photo.File)
	if err != nil {
		log.Printf("error reading provided photo: %v", err)
		return
	}
	defer rc.Close()

	image, err := ioutil.ReadAll(rc)
	if err != nil {
		log.Printf("error reading provided photo: %v", err)
		return
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	frChan := b.invokeFaceAPI(ctx, m, image)
	faceResult := <-frChan
	if faceResult == nil {
		b.tbot.Send(m.Chat, "Hello, unfortunately I'm not smart enough to analyze photos!")
		return
	}

	faces := faceResult.Faces
	switch len(faces) {
	case 0:
		b.processNoFacePhoto(m.Chat, image)
	case 1:
		b.processSingleFacePhoto(m.Chat, faces[0])
	default:
		b.processGroupPhoto(m.Chat, faces)
	}
}

func (b *Bot) invokeFaceAPI(ctx context.Context, m *tb.Message, image []byte) chan *wssface.FaceResult {
	frChan := make(chan *wssface.FaceResult)
	if b.faceCli == nil {
		defer close(frChan)
		return frChan
	}

	go func() {
		defer close(frChan)

		rct := ioutil.NopCloser(bytes.NewReader(image))
		defer rct.Close()
		faceResult, err := b.faceCli.InvokeFace(ctx, rct)

		if err != nil {
			log.Printf("error invoking face API: %v", err)
		}
		frChan <- faceResult
	}()

	return frChan
}

func (b *Bot) processGroupPhoto(chat *tb.Chat, faces []wssface.FaceDetails) {
	message := fmt.Sprintf("What a nice group picture of %v of you", len(faces))
	b.tbot.Send(chat, message)
}

func (b *Bot) processSingleFacePhoto(chat *tb.Chat, face wssface.FaceDetails) {
	message := fmt.Sprintf("Hello %s, I guess you are %v years old. Why you so %s?",
		b.genderGreet(face.Gender), face.Age, face.Sentiment.Adjective())

	b.tbot.Send(chat, message)
}

func (b *Bot) processNoFacePhoto(chat *tb.Chat, image []byte) {
	ctx := context.Background()
	visionChan, formRecognizerChan := make(chan bool), make(chan bool)
	defer func() {
		close(visionChan)
		close(formRecognizerChan)
	}()

	go func() {
		rct := ioutil.NopCloser(bytes.NewReader(image))
		chanres := true
		defer func() {
			rct.Close()
			visionChan <- chanres
		}()

		res, err := b.visionCli.InvokeVision(ctx, rct)
		if err != nil {
			log.Printf(`error invoking computer vision service: %v`, err)
			b.tbot.Send(chat, `This picture makes me feel a sick! I'm sorry, I can't handle this request!`)
			chanres = false
			return
		}

		if res.Description == nil {
			b.tbot.Send(chat, `I'm sorry but I can't figure out what this picture represents`)
			return
		}

		message := fmt.Sprintf("It seems %s", *res.Description)
		b.tbot.Send(chat, message)
	}()

	go func() {
		rct := ioutil.NopCloser(bytes.NewReader(image))
		chanres := true
		defer func() {
			rct.Close()
			formRecognizerChan <- chanres
		}()

		res, err := b.formRecognizerCli.InvokeFormRecognizer(ctx, rct)
		if err != nil {
			log.Printf(`error invoking form recognizer service: %v`, err)
			chanres = false
			return
		}

		if !res.IsSucceeded() {
			b.tbot.Send(chat, `No receipt found in photo`)
			return
		}

		var message string
		total := res.Total()
		if total != nil {
			message = fmt.Sprintf("It seems you spent %s", *total)
		} else {
			message = "Can't figure out how much you spent"
		}

		if merchant := res.MerchantName(); merchant != nil {
			message += " at " + *merchant
		}
		if address := res.MerchantAddress(); address != nil {
			message += " in " + *address
		}
		if date := res.TransactionDate(); date != nil {
			message += " in date " + *date
		}

		b.tbot.Send(chat, message)
	}()

	vres := <-visionChan
	frres := <-formRecognizerChan
	if !vres && !frres {
		b.tbot.Send(chat, `This picture makes me feel a sick! I'm sorry, I can't handle this request!`)
	}
}

func (b *Bot) genderGreet(gender string) string {
	if gender == "male" {
		return "man"
	}
	return "darling"
}
