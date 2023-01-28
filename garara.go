package garara

import (
	"encoding/xml"
	"errors"
	"sync"
)

const UNKNOWN = "Unknown"

type MailRoot struct {
	XMLName xml.Name `xml:"mail"`
}

type CDATAString struct {
	CDATA string `xml:",cdata"`
}

func SimpleV1SendListBuilder(addressList []string) SendList {
	l := make([]Data, 0, len(addressList))
	for k, v := range addressList {
		var d Data
		d.ID = k
		d.Address.Value = v
		l = append(l, d)
	}
	var s SendList
	s.Data = l
	return s
}

func SimpleV1SendListDataBuilder(id int, deviceType DeviceType, address string, intText, extText, extImage []string, keyField string) (d Data, err error) {
	if deviceType.String() == "Unknown" {
		err = errors.New("DeviceType Unknown")
		return d, err
	}

	d.ID = id
	d.Address.Device = deviceType
	d.Address.Value = address
	d.KeyField = keyField

	i := make([]AttrIdCdata, 0, len(intText))
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for k, v := range intText {
			var it AttrIdCdata
			it.ID = k
			it.CDATA = v
			i = append(i, it)
		}
		d.IntText = i
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		d.ExtText = AttrIdStringsBuilder(extText)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		d.ExtImage = AttrIdStringsBuilder(extImage)
	}()
	wg.Wait()

	return d, nil
}

func AttrIdCDataBuilder(ss []string) []AttrIdCdata {
	size := len(ss)
	if size == 0 {
		return nil
	}

	a := make([]AttrIdCdata, 0, size)
	for k, v := range ss {
		var c AttrIdCdata
		c.ID = k
		c.CDATA = v
		a = append(a, c)
	}
	return a
}

func AttrIdStringsBuilder(ss []string) []AttrIdString {
	size := len(ss)
	if size == 0 {
		return nil
	}

	a := make([]AttrIdString, 0, size)
	for k, v := range ss {
		var s AttrIdString
		s.ID = k
		s.Value = v
		a = append(a, s)
	}
	return a
}

// SimpleV1MailContentsBuilder はメール送信時のcontentを作成します。
// imagesはContent.Imageにid,0から順番にデータを作成します。
// textsはContent.Textにid,0から順番にデータを作成します。
// filesがattachMaxFileSize（現仕様では3）よりデータ量が大きい場合エラーを返却します。
func SimpleV1MailContentsBuilder(subject, body string, part PartType, images, texts, files []string) (cont Contents, err error) {
	if len(files) > attachMaxFileSize {
		return cont, errors.New("too many files")
	}

	if part.String() == "Unknown" {
		err = errors.New("PartType Unknown")
	}

	cont.Subject.CDATA = subject
	cont.Body.Part = part
	cont.Body.CDATA = body

	cont.Image = AttrIdStringsBuilder(images)

	ts := make([]AttrIdCdata, 0, len(texts))
	for k, t := range texts {
		var idc AttrIdCdata
		idc.ID = k
		idc.CDATA = t
		ts = append(ts, idc)
	}
	cont.Text = ts
	cont.AttachFile = files

	return cont, nil
}

func SimpleV1SettingBuilder(sendData, fromName, from, envelopeFrom string, throttle int, smime, opened UseType, option Option) (s Setting, err error) {
	err = nil
	if smime.String() == UNKNOWN {
		err = errors.New("smime type is unknown")
		return s, err
	}

	if opened.String() == UNKNOWN {
		err = errors.New("opened type is unknown")
		return s, err
	}

	s.SendDate = sendData
	s.FromName.CDATA = fromName
	s.From = from
	s.EnvelopeFrom = envelopeFrom
	s.Option = option
	s.Throttle = throttle
	s.SMime.Use = smime
	s.OpenedFlag.Use = opened

	return s, err
}

// V1MailDeliveryBuilder はメール送信用のDeliveryを作成します。
// Deliveryはメール送信用以外の項目もあるため、本builderを利用することで設定項目を特定します。
// 引数に取る各structについても各SimpleBuilderを用意しています。
func V1MailDeliveryBuilder(id int, requestID string, setting Setting, contents Contents, list SendList) (d Delivery) {
	d.Action = RESERVE
	d.ID = id
	d.RequestID = requestID
	d.Setting = setting
	d.Contents = contents
	d.SendList = list
	return d
}
