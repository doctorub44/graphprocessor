package graphproc

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"os"
)

//FileState : store that state of file scanner
type FileState struct {
	Status   string
	Ffile    *os.File
	Fscanner *bufio.Scanner
}

//SubGraph : graph stage to execute a sub-graph
//	The subgraph must have been made public by PublishGraphs()
func SubGraph(s *State, payload *Payload) error {
	graphs := PublicGraphs()
	name, ok := s.Config("name")
	if !ok {
		return errors.New("SubGraph: unable to lookup configuration field: " + name)
	}
	return graphs.Execute(name, payload)
}

//FileWrite : graph stage to write the raw bytes to a file
func FileWrite(s *State, payload *Payload) error {
	var fstate FileState
	var err error
	var file *os.File

	if s == nil {
		return errors.New("FileWrite: no file state")
	}

	name, ok := s.Config("name")
	if !ok || name == "" {
		return errors.New("FileWrite: unable to lookup configuration field: " + name)
	}

	if file, err = os.Create(name); err != nil {
		s.appstate = fstate
		return err
	}
	_, err = io.WriteString(file, string(payload.Raw))
	file.Close()

	s.appstate = fstate

	return err
}

//FileRead : graph stage to read a file into the raw bytes
func FileRead(s *State, payload *Payload) error {
	var fstate FileState
	var err error
	var file *os.File

	if s == nil {
		return errors.New("FileRead: no file state")
	}

	name, ok := s.Config("name")
	if !ok || name == "" {
		return errors.New("FileRead: unable to lookup configuration field: " + name)
	}

	if file, err = os.Open(name); err != nil {
		s.appstate = fstate
		return err
	}

	p, err := ioutil.ReadAll(file)
	payload.Raw = append(payload.Raw, p...)
	file.Close()

	s.appstate = fstate
	return err
}

//FileTextWriter : graph stage to write a stream to a text file one line at a time
//	looks for closed status in payload.Data to close the file
func FileTextWriter(s *State, payload *Payload) error {
	var fstate FileState
	var err error

	if s != nil {
		if s.appstate == nil {
			s.appstate = fstate
		} else {
			fstate = s.appstate.(FileState)
		}
	} else {
		return errors.New("FileTextWriter: no state passed")
	}

	if fstate.Status == "closed" {
		return errors.New("FileTextWriter: file stream closed")
	}

	if fstate.Status == "" {
		name, ok := s.Config("name")
		if !ok || name == "" {
			return errors.New("FileTextWriter: unable to lookup configuration field: " + name)
		}
		if fstate.Ffile, err = os.Create(name); err != nil {
			s.appstate = fstate
			return err
		}
		fstate.Status = "open"
	}

	status, err := payload.GetData("stream")
	if status.(string) == "closed" {
		fstate.Ffile.Close()
		fstate.Status = "closed"
	} else {
		_, err = io.WriteString(fstate.Ffile, string(payload.Raw)+"\n")
	}

	s.appstate = fstate
	return err
}

//FileTextReader : graph stage to read a text file one line at a time and send as a stream
//	state of the stream, open or closed, is stored in the payload.Data field. Closed is sent after
//	file has been read
func FileTextReader(s *State, payload *Payload) error {
	var fstate FileState
	var err error

	if s != nil {
		if s.appstate == nil {
			s.appstate = fstate
		} else {
			fstate = s.appstate.(FileState)
		}
	} else {
		return errors.New("FileTextReader: no state passed")
	}

	if fstate.Status == "closed" {
		return errors.New("FileTextReader: file stream closed")
	}

	if fstate.Ffile == nil {
		name, ok := s.Config("name")
		if !ok || name == "" {
			return errors.New("FileTextReader: unable to lookup configuration field: " + name)
		}
		if fstate.Ffile, err = os.Open(name); err != nil {
			s.appstate = fstate
			return err
		}
		fstate.Fscanner = bufio.NewScanner(fstate.Ffile)
		fstate.Fscanner.Split(bufio.ScanLines)
		fstate.Status = "open"
		payload.SetData("stream", "open")
	}

	result := fstate.Fscanner.Scan()
	if result {
		text := []byte(fstate.Fscanner.Text())
		payload.Raw = append(payload.Raw[:0], text...)
	} else {
		fstate.Ffile.Close()
		payload.SetData("stream", "closed")
		fstate.Status = "closed"
	}

	s.appstate = fstate
	return nil
}
