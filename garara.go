package garara

import (
	"encoding/xml"
	"errors"
	"sync"
)

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

	i := make([]AttrIdCdata, 0, len(intText))
	var et, ei []AttrIdString
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
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		et = attrIdStringsBuilder(extText)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ei = attrIdStringsBuilder(extImage)
	}()
	wg.Wait()
	d.IntText = i
	d.ExtText = et
	d.ExtImage = ei

	d.KeyField = keyField

	return d, nil
}

func attrIdStringsBuilder(ss []string) []AttrIdString {
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
