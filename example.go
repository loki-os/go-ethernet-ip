package go_ethernet_ip

import (
	"log"
	"sync"
)

func Init(ip string) {
	conn, err := NewTCP(ip, nil)
	if err != nil {
		// cannot resolve host
		log.Fatalln(err)
	}

	err = conn.Connect()
	if err != nil {
		// cannot connect to host
		log.Fatalln(err)
	}

	Tags, err := conn.AllTags()
	if err != nil {
		// cannot get tags
		log.Fatalln(err)
	}

	targetTag := Tags["tagName"]

	err = targetTag.Read()
	if err != nil {
		// cannot read tag
		log.Fatalln(err)
	}

	// IF U Need any other type of tag, you should implement it yourself.only support int32 and string.
	i32Result := targetTag.Int32()
	log.Println(i32Result)

	stringResult := targetTag.String()
	log.Println(stringResult)

	targetTag.SetInt32(123)
	err = targetTag.Write()
	if err != nil {
		// cannot write tag
		log.Fatalln(err)
	}

	// tag group, multiple tags sync read/write.
	lock := new(sync.Mutex)
	tagGroup := NewTagGroup(lock)
	targetTag1 := Tags["tagName1"]
	targetTag2 := Tags["tagName2"]
	tagGroup.Add(targetTag1)
	tagGroup.Add(targetTag2)
	targetTag1.SetInt32(123)
	targetTag2.SetString("hello")
	err = tagGroup.Write()
	if err != nil {
		// cannot write tag
		log.Fatalln(err)
	}

	// lower level api, you can use it directly.
	identities, err := conn.ListIdentity()
	log.Println(identities)

	interfaces, err := conn.ListInterface()
	log.Println(interfaces)

	services, err := conn.ListServices()
	log.Println(services)

	// open connection forward, but it disconnected after few times.
	// if you need this future you need to fix it.
	// actually, I don't have abplc device to test after business done.
	conn.ForwardOpen()
	var tag = new(Tag)
	conn.InitializeTag("OP.UDT_Alarm.DINT_065_096", tag)
	log.Println("Name: ", tag.Name())
	log.Println(tag.GetValue())

	tag = new(Tag)
	conn.InitializeTag("wowtag", tag)
	log.Println("Name: ", tag.Name())
	log.Println(tag.GetValue())

	tag = new(Tag)
	conn.InitializeTag("wotag[0]", tag)
	log.Println("Name: ", tag.Name())
	log.Println(tag.GetValue())

	tag = new(Tag)
	conn.InitializeTag("wotag[1]", tag)
	log.Println("Name: ", tag.Name())
	log.Println(tag.GetValue())

	tag = new(Tag)
	conn.InitializeTag("OP_Format[0].REAL_Performance[0]", tag)
	log.Println("Name: ", tag.Name())
	//log.Println("Type: ", tag.Type)
	log.Println(tag.GetValue())

	tag = new(Tag)
	conn.InitializeTag("wwwtag[1,0,1]", tag)
	log.Println("Name: ", tag.Name())
	log.Println(tag.GetValue())

	tag = new(Tag)
	conn.InitializeTag("stringtag", tag)
	log.Println("Name: ", tag.Name())
	log.Println(tag.GetValue())
}
