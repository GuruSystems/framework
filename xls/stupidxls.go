package stupid

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/tealeg/xlsx"
	"github.com/GuruSystems/framework/auth"
	"math/rand"
	"strings"
	"time"
)

type Position struct {
	Name  string
	Sheet string
	Col   string
	Row   int
}

type RFC struct {
	Requestor string
	Approver  string
	ID        string
	Type      string
	Text      string
	file      *xlsx.File
}

var (
	phone    = flag.String("phone", "123456", "Phonenumber for requestor and approver")
	customer = flag.String("customer", "Foo bar", "Customer name")
	layout   []Position
	src      = rand.NewSource(time.Now().UnixNano())
)

/**************************************************
* helpers
***************************************************/
//https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
func RandString(n int) string {
	const letterBytes = "1234567890:-_=abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits

	)

	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

/**************************************************
* sheet stuff
***************************************************/

func AddPositions() {
	layout = layout[0:]
	AddPosition("customer", "RFC Form", "A", 4)
	AddPosition("requestor.name", "RFC Form", "D", 4)
	AddPosition("requestor.email", "RFC Form", "D", 5)
	AddPosition("requestor.telephone", "RFC Form", "D", 6)
	AddPosition("approver.name", "RFC Form", "I", 4)
	AddPosition("approver.email", "RFC Form", "I", 5)
	AddPosition("approver.telephone", "RFC Form", "I", 6)
	AddPosition("changetype", "RFC Form", "A", 9)
	AddPosition("configitem", "RFC Form", "A", 12)
	AddPosition("discipline", "RFC Form", "C", 9)
	AddPosition("changewindow.start", "RFC Form", "I", 9)
	AddPosition("changewindow.end", "RFC Form", "I", 10)
	AddPosition("changewindow.tz", "RFC Form", "I", 11)
	AddPosition("changetitle", "RFC Form", "A", 15)
	AddPosition("customer.reference", "RFC Form", "C", 12)
	AddPosition("businesscase", "RFC Form", "A", 22)
	AddPosition("impact", "RFC Form", "F", 22)
	AddPosition("rolloutplan", "RFC Form", "A", 25)
	AddPosition("rollbackplan", "RFC Form", "A", 29)
	AddPosition("testplan", "RFC Form", "A", 33)

	AddPosition("acl.type", "Network RFC Details", "B", 7)
	AddPosition("acl.device", "Network RFC Details", "C", 7)
	AddPosition("acl.interface", "Network RFC Details", "D", 7)
	AddPosition("acl.comment", "Network RFC Details", "E", 7)
	AddPosition("acl.purpose", "Network RFC Details", "F", 7)
	AddPosition("acl.source", "Network RFC Details", "G", 7)
	AddPosition("acl.destination", "Network RFC Details", "H", 7)
	AddPosition("acl.port", "Network RFC Details", "I", 7)
	AddPosition("acl.proto", "Network RFC Details", "J", 7)
	AddPosition("acl.action", "Network RFC Details", "K", 7)
	AddPosition("acl.bidirectional", "Network RFC Details", "L", 7)
}

func AddPosition(name string, sheet string, col string, row int) {
	layout = append(layout, Position{
		Name:  name,
		Sheet: sheet,
		Col:   col,
		Row:   row,
	})
}
func ColIDX(name string) int {
	var res int
	l := strings.ToUpper(name)
	b := []byte(l)
	res = int(b[0])
	res = res - 65
	return res
}
func (rfc *RFC) GetBuffer() *bytes.Buffer {
	buf := bytes.NewBuffer(nil)
	rfc.file.Write(buf)
	return buf
}
func (rfc *RFC) Set(name string, value string) {
	for _, pos := range layout {
		if pos.Name == "" {
			continue
		}
		if pos.Name == name {
			colpos := ColIDX(pos.Col)
			//fmt.Printf("Found %s in %s %d,%d\n", name, pos.Col, colpos, pos.Row)
			sheet := rfc.file.Sheet[pos.Sheet]
			cell := sheet.Cell(pos.Row-1, colpos)
			cell.SetValue(value)
			//fmt.Println("Cell:", cell)
			return
		}
	}
	fmt.Println("Not found: %s\n", name)
}

func (rfc *RFC) SetRequestor(user auth.User) {
	rfc.Requestor = user.Email
	rfc.Set("requestor.email", user.Email)
	rfc.Set("requestor.name", fmt.Sprintf("%s %s", user.FirstName, user.LastName))
	rfc.Set("requestor.telephone", *phone)
}
func (rfc *RFC) SetAuthoriser(user auth.User) {
	rfc.Approver = user.Email
	rfc.Set("approver.email", user.Email)
	rfc.Set("approver.name", fmt.Sprintf("%s %s", user.FirstName, user.LastName))
	rfc.Set("approver.telephone", *phone)
}
func (rfc *RFC) SetStandardStuff(user auth.User) {
	ref := fmt.Sprintf("%s-%s", user.FirstName, RandString(20))
	rfc.ID = ref
	rfc.Set("customer", *customer)
	fm := "2006-01-02 15:04"
	rfc.Set("changewindow.start", time.Now().Format(fm))
	rfc.Set("changewindow.end", time.Now().AddDate(0, 0, 7).Format(fm))
	rfc.Set("changewindow.tz", "GMT")
	rfc.Set("customer.reference", ref)
	rfc.Set("changetype", "Standard")
	rfc.Set("businesscase", "This change is part of \"normal operations\" and pre-approved without businesscase")
	rfc.Set("impact", "Infrastructure to operate normally afterwards")
}
func (rfc *RFC) Save(filename string) error {
	return rfc.file.Save(filename)
}
func (rfc *RFC) Print() {
	for _, pos := range layout {
		if pos.Name == "" {
			continue
		}
		colpos := ColIDX(pos.Col)
		sheet := rfc.file.Sheet[pos.Sheet]
		cell := sheet.Cell(pos.Row-1, colpos)
		v := cell.Value
		fmt.Printf("%s = %s\n", pos.Name, v)
	}
}

func (rfc *RFC) SetTitle(title string) {
	rfc.Set("changetitle", title)
}

type OpenPortDef struct {
	PortDef  string
	VMName   string
	IP       string
	FromAddr string
	Nic      string
}

type DeleteServerDef struct {
	VMName string
	IP     string
}
type FreeTextDef struct {
	VMName string
	IP     string
	Text   string
}

/**************************************************
* RFC Create stuff
***************************************************/

func OpenPortRFC(stupidTemplate string, user auth.User, def OpenPortDef) *RFC {
	var err error
	AddPositions()

	rfc := new(RFC)
	rfc.Type = "OpenPort"
	rfc.file, err = xlsx.OpenFile(stupidTemplate)
	if err != nil {
		fmt.Printf("Failed to open %s: %s", stupidTemplate, err)
		return nil
	}
	fmt.Printf("Creating rfc for %s %s (%s)\n", user.FirstName, user.LastName, user.Email)
	rfc.SetStandardStuff(user)
	rfc.SetRequestor(user)
	rfc.SetAuthoriser(user)
	rfc.Set("rolloutplan", "Softcat to execute configuration change and inform customer upon completion. Customer verifies correct behavior and accepts or declines change")
	rfc.Set("rollbackplan", "Softcat to execute configuration rollback and inform customer upon completion")
	rfc.Set("testplan", "send traffic to given port from specified source address and verify traffic is arriving at target VM")
	rfc.Set("discipline", "Networks & Security")
	title := fmt.Sprintf("Open ports %s from %s to %s on machine %s", def.PortDef, def.FromAddr, def.IP, def.VMName)
	rfc.Text = title
	rfc.Set("configitem", def.IP)
	rfc.SetTitle(title)

	rfc.Set("acl.type", "Amendment")
	rfc.Set("acl.device", fmt.Sprintf("%s/%s", def.VMName, def.IP))
	rfc.Set("acl.interface", def.Nic)
	rfc.Set("acl.source", def.FromAddr)
	rfc.Set("acl.destination", def.IP)
	rfc.Set("acl.port", def.PortDef)
	rfc.Set("acl.proto", "UDP4 & TCP4")
	rfc.Set("acl.action", "Allow")
	rfc.Set("acl.bidirectional", "Yes")
	err = rfc.file.Save("/tmp/new-rfc.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
	return rfc
}

func DeleteServerRFC(stupidTemplate string, user auth.User, def DeleteServerDef) *RFC {
	var err error
	AddPositions()

	rfc := new(RFC)
	rfc.Type = "DeleteServer"
	rfc.file, err = xlsx.OpenFile(stupidTemplate)
	if err != nil {
		fmt.Printf("Failed to open %s: %s", stupidTemplate, err)
		return nil
	}
	fmt.Printf("Creating rfc for %s %s (%s)\n", user.FirstName, user.LastName, user.Email)
	rfc.SetStandardStuff(user)
	rfc.SetRequestor(user)
	rfc.SetAuthoriser(user)
	rfc.Set("rolloutplan", "Softcat to execute configuration change and inform customer upon completion. Customer verifies correct behavior and accepts or declines change")
	rfc.Set("rollbackplan", "Softcat to execute configuration rollback and inform customer upon completion")
	rfc.Set("testplan", "verify that server has been deleted from vSphere")
	rfc.Set("discipline", "Storage & Virtualisation")
	title := fmt.Sprintf("Delete server %s on IP address %s from vSphere", def.VMName, def.IP)
	rfc.Text = title
	rfc.Set("configitem", def.IP)
	rfc.SetTitle(title)

	err = rfc.file.Save("/tmp/new-rfc.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
	return rfc
}

func FreeTextRFC(stupidTemplate string, user auth.User, def FreeTextDef) *RFC {
	var err error
	AddPositions()

	rfc := new(RFC)
	rfc.Type = "FreeTextServer"
	rfc.file, err = xlsx.OpenFile(stupidTemplate)
	if err != nil {
		fmt.Printf("Failed to open %s: %s", stupidTemplate, err)
		return nil
	}
	fmt.Printf("Creating rfc for %s %s (%s)\n", user.FirstName, user.LastName, user.Email)
	rfc.SetStandardStuff(user)
	rfc.SetRequestor(user)
	rfc.SetAuthoriser(user)
	rfc.Set("rolloutplan", "Softcat to execute configuration change and inform customer upon completion. Customer verifies correct behavior and accepts or declines change")
	rfc.Set("rollbackplan", "Softcat to execute configuration rollback and inform customer upon completion")
	rfc.Set("testplan", "customer will verify")
	rfc.Set("discipline", "Networks & Security")
	title := fmt.Sprintf("%s", def.Text)
	rfc.Text = title
	rfc.Set("configitem", def.IP)
	rfc.SetTitle(title)

	err = rfc.file.Save("/tmp/new-rfc.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
	return rfc
}
