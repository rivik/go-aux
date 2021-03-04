package rvdb

import (
	"encoding/json"
	"testing"
)

func TestDSNConfig(t *testing.T) {
	conf := []byte(`{
	"dsn": "mongodb://localhost:27017,localhost:27018,localhost:27019/?replicaSet=rs1",
	"user": "testuser",
	"password": "",
	"password_file": ".test.dbpass",
	"db_params": {
		"timeout": 30,
		"alt_hosts": "host1,host2",
		"mode": "test"
	}
}`)

	dsnConf := DSNConfig{}
	err := json.Unmarshal(conf, &dsnConf)
	if err != nil {
		t.Error("Can't parse config: ", err)
		return
	}

	err = dsnConf.Prepare()
	if err != nil {
		t.Error("Can't prepare dsn config: ", err)
	}

	// db-driver specific
	dsnConf.DBParams.Add("username", dsnConf.User)
	dsnConf.DBParams.Add("password", dsnConf.Password)

	want := `mongodb://localhost:27017,localhost:27018,localhost:27019/?alt_hosts=host1%2Chost2&mode=test&password=testsecretpass&replicaSet=rs1&timeout=30&username=testuser`
	got := dsnConf.DSN().String()

	if got != want {
		t.Errorf("Bad result DSN: got %v; want %v", got, want)
	}
}
