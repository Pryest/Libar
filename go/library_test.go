package main

import (
	"testing"
)

func TestMain(t *testing.T) {
	var err error
	lib := Library{}
	lib.ConnectDB()
	err = lib.CreateTables()
	if err != nil {
		t.Errorf("can't create tables")
	}
	err = lib.AddS()
	if err != nil {
		t.Errorf("can't add student")
	}
	err = lib.AddBook("name","people","1234561234560")
	if err != nil {
		t.Errorf("can't add book")
	}
	err = lib.RemoveBook(1,"lost")
	if err != nil {
		t.Errorf("can't remove book")
	}
	err = lib.AddBook("name2","people2","1234561334560")
	if err != nil {
		t.Errorf("can't add book")
	}
	err = lib.FindT("name2")
	if err != nil {
		t.Errorf("can't find book")
	}
	err =lib.Borrow(1,1);
	if err != nil {
		t.Errorf("can't borrow book")
	}
	err =lib.Borrow(1,2);
	if err != nil {
		t.Errorf("can't borrow book")
	}
	err =lib.History(1);
	if err != nil {
		t.Errorf("can't fetch history")
	}
	err =lib.BRS(1);
	if err != nil {
		t.Errorf("can't see borrows")
	}
	err =lib.DueT(2);
	if err != nil {
		t.Errorf("can't see duetime")
	}
	err =lib.ConT(1,2);
	if err != nil {
		t.Errorf("can't continue")
	}
	err =lib.Check(1);
	if err != nil {
		t.Errorf("can't check overdue")
	}
	err =lib.ReturnB(1,2);
	if err != nil {
		t.Errorf("can't return book")
	}
}
