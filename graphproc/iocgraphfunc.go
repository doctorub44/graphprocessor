package graphproc

import (
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type IOCs struct {
	Md5    []string `json:"md5"`
	Sha1   []string `json:"sha1"`
	Sha256 []string `json:"sha256"`
	Ipv4   []string `json:"ipv4"`
	Ipv6   []string `json:"ipv6"`
	Url    []string `json:"url"`
	Domain []string `json:"domain"`
}

var remd5, resha1, resha256, reipv4, reipv6, reurl, redomain *regexp.Regexp

//DeleteWhiteListIOC :
func DeleteWhiteListIOC(s *State, payload *Payload) error {
	var err error = nil

	return err
}

//WhiteListIOC :
func WhiteListIOC(s *State, payload *Payload) error {
	var err error = nil

	return err
}

//FilterWhiteListIOC :
func FilterWhiteListIOC(s *State, payload *Payload) error {
	var err error = nil

	return err
}

//CutFields :
func CutFields(s *State, payload *Payload) error {
	//Convert the payload raw buffer into a m x n matrix of strings
	matrix := BytesMatrix(payload.Raw)
	//Build bitmap row of fields to cut
	cutrow, err := bitmaprow(s, matrix[0])
	if err != nil {
		return errors.New("CutFields: unable to build rows to cut [" + err.Error() + "]")
	}
	//Using the cut row map, cut columns from the matrix
	matrix = MatrixCut(matrix, cutrow)
	//Empty the payload raw buffer and fill it the matrix converted to with string rows with \n and converted to bytes
	payload.Raw = payload.Raw[:0]
	payload.Raw, err = MatrixBytes(matrix, payload.Raw)
	return err
}

//SelectFields :
func SelectFields(s *State, payload *Payload) error {
	//Convert the payload raw buffer into a m x n matrix of strings
	matrix := BytesMatrix(payload.Raw)
	//Build bitmap row of fields to select
	selrow, err := bitmaprow(s, matrix[0])
	if err != nil {
		return errors.New("SelectFields: unable to build rows to select [" + err.Error() + "]")
	}
	//Using the select row map, select columns from the matrix
	matrix = MatrixSelect(matrix, selrow)
	//Empty the payload raw buffer and fill it the matrix converted to with string rows with \n and converted to bytes
	payload.Raw = payload.Raw[:0]
	payload.Raw, err = MatrixBytes(matrix, payload.Raw)
	return err
}

//bitmaprow : build bitmap row of fields to select
func bitmaprow(s *State, row []string) ([]bool, error) {
	var bitmaprow []bool

	for range row {
		bitmaprow = append(bitmaprow, false)
	}

	if flds, ok := s.Config("fields"); ok {
		flist := strings.Split(flds, " ")
		for _, f := range flist {
			col, err := strconv.Atoi(strings.TrimSpace(f))
			if err != nil {
				return nil, err
			}
			if col > len(bitmaprow) {
				return nil, errors.New("bitmaprow: column specified is greater than length of row")
			}
			bitmaprow[col] = true
		}
	}
	return bitmaprow, nil
}

//NormalIOC : convert some typical values to normal values
func NormalIOC(s *State, payload *Payload) error {
	var err error = nil
	var text string

	text = string(payload.Raw)
	text = strings.ReplaceAll(text, "hxxp", "http")
	text = strings.ReplaceAll(text, "http[:]", "http:")
	text = strings.ReplaceAll(text, "https[:]", "https:")
	payload.Raw = append(payload.Raw[:0], []byte(text)...)

	return err
}

//Md5IOC : search and return MD5
func Md5IOC(s *State, payload *Payload) error {
	var err error = nil
	if remd5 == nil {
		remd5 = regexp.MustCompile(`\b[a-fA-F0-9]{32}\b`)
	}
	payload.Raw, err = IOC(remd5, payload.Raw, "md5")
	return err
}

//Sha1IOC : search and return SHA1
func Sha1IOC(s *State, payload *Payload) error {
	var err error = nil
	if resha1 == nil {
		resha1 = regexp.MustCompile(`\b[a-fA-F0-9]{40}\b`)
	}
	payload.Raw, err = IOC(resha1, payload.Raw, "sha1")
	return err
}

//Sha256IOC : search and return SHA256
func Sha256IOC(s *State, payload *Payload) error {
	var err error = nil
	if resha256 == nil {
		resha256 = regexp.MustCompile(`\b[a-fA-F0-9]{64}\b`)
	}
	payload.Raw, err = IOC(resha256, payload.Raw, "sha256")
	return err
}

//Ipv4IOC : search and return IPV4
func Ipv4IOC(s *State, payload *Payload) error {
	var err error = nil
	if reipv4 == nil {
		reipv4 = regexp.MustCompile(`\b(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\b`)
	}
	payload.Raw, err = IOC(reipv4, payload.Raw, "ipv4")
	return err
}

//Ipv6IOC : search and return IPV6
func Ipv6IOC(s *State, payload *Payload) error {
	var err error = nil
	if reipv6 == nil {
		reipv6 = regexp.MustCompile(`\b(?:[a-fA-F0-9]{1,4}:){7}[a-fA-F0-9]{1,4}\b`)
	}
	payload.Raw, err = IOC(reipv6, payload.Raw, "ipv6")
	return err
}

//UrlIOC : search and return URL
func UrlIOC(s *State, payload *Payload) error {
	var err error = nil
	if reurl == nil {
		reurl = regexp.MustCompile(`(http|ftp|https)://([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@?^=%&:/~+#-]*[\w@?^=%&/~+#-])?`)
	}
	payload.Raw, err = IOC(reurl, payload.Raw, "url")
	return err
}

//DomainIOC : search and return DOMAIN
func DomainIOC(s *State, payload *Payload) error {
	var err error = nil
	if redomain == nil {
		redomain = regexp.MustCompile(``)
	}
	payload.Raw, err = IOC(reurl, payload.Raw, "domain")
	return err
}

//IOC : search and return the IOC
func IOC(re *regexp.Regexp, rawbytes []byte, hashtype string) ([]byte, error) {
	var err error = nil
	iocs := re.FindAllString(string(rawbytes), -1)
	rawbytes = rawbytes[:0]
	//size := len(iocs)

	for _, ioc := range iocs {
		newbytes := append([]byte(hashtype+"|"), ioc...)
		rawbytes = append(rawbytes, newbytes...)
		rawbytes = append(rawbytes, []byte(",")...)
	}

	return rawbytes, err
}

//IOCtoData : conver the Raw field IOC strings to a slice in the Data field
func IOCtoData(s *State, payload *Payload) error {
	var err error = nil

	payload.Data = strings.Split(string(payload.Raw), ",")

	return err
}

//IOCDataToJson : marshall the Data field to JSON in the Raw field
func IOCDataToJson(s *State, payload *Payload) error {
	var err error = nil
	var iocs IOCs
	var fields []string

	if payload.Data != nil {
		switch t := payload.Data.(type) {
		case []string:
			fields = payload.Data.([]string)
		default:
			return errors.New("IOCtoJSON: Data field invalid type passed: " + t.(string))
		}
	} else {
		return errors.New("IOCtoJSON: missing Data field")
	}

	for _, f := range fields {
		if strings.TrimSpace(f) == "" {
			continue
		}
		ioc := strings.Split(f, "|")
		switch ioc[0] {
		case "md5":
			iocs.Md5 = append(iocs.Md5, ioc[1])
		case "sha1":
			iocs.Sha1 = append(iocs.Sha1, ioc[1])
		case "sha256":
			iocs.Sha256 = append(iocs.Sha256, ioc[1])
		case "ipv4":
			iocs.Ipv4 = append(iocs.Ipv4, ioc[1])
		case "ipv6":
			iocs.Ipv6 = append(iocs.Ipv6, ioc[1])
		case "url":
			iocs.Url = append(iocs.Url, ioc[1])
		default:
			continue
		}
	}
	text, _ := json.Marshal(&iocs)
	nullioc := regexp.MustCompile(`"[a-z0-9]*":null,?`)
	trailcomm := regexp.MustCompile(`,}`)
	text = []byte(nullioc.ReplaceAllString(string(text), ""))
	text = []byte(trailcomm.ReplaceAllString(string(text), "}"))
	payload.Raw = append(payload.Raw[:0], text...)

	return err
}
