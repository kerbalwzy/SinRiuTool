package plugins

import (
	"testing"
)

func TestErpCli(t *testing.T) {
	cli := NewSinRiuErpCli()
	err := cli.Login("????", "????")
	if nil != err {
		t.Fatal(err)
	}
	err = cli.FlushProfile()
	if nil != err {
		t.Fatal(err)
	} else {
		t.Log(cli.ErpProfile)
	}
}
