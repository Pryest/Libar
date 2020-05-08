package main

import (
	"fmt"
	"time"
	"strings"
	// mysql connector
	_ "github.com/go-sql-driver/mysql"
	sqlx "github.com/jmoiron/sqlx"
)

const (
	User     = "root"
	Password = "123456"
	DBName   = "ass3"
)

type Library struct {
	db *sqlx.DB
}

func (lib *Library) ConnectDB() error{
	db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", User, Password, DBName))
	if err != nil {
		panic(err)
	}
	lib.db = db
	return nil;
}

// CreateTables created the tables in MySQL
func (lib *Library) CreateTables() error{
	_, _ = lib.db.Exec(fmt.Sprintf("create table S(S int not null, primary key(S));"))
	_, _ = lib.db.Exec(fmt.Sprintf("create table T(T int not null,primary key(T));"))
	_, _ = lib.db.Exec(fmt.Sprintf("create table BT(B int not null,TITLE varchar(40) not null,primary key(B));"))
	_, _ = lib.db.Exec(fmt.Sprintf("create table BA(B int not null,AUTHOR varchar(40) not null,primary key(B));"))
	_, _ = lib.db.Exec(fmt.Sprintf("create table BI(B int not null,ISBN char(13) not null,primary key(B));"))
	_, _ = lib.db.Exec(fmt.Sprintf("create table BB(B int not null,BR char(1) not null,primary key(B));"))
	_, _ = lib.db.Exec(fmt.Sprintf("create table BE(B int not null,EXP varchar(40) not null,primary key(B));"))
	_, _ = lib.db.Exec(fmt.Sprintf("create table SBC(S int not null, B int not null, C int not null, primary key(S,B), foreign key(S) references S(S), foreign key(B) references BT(B));"))
	_, _ = lib.db.Exec(fmt.Sprintf("create table SBD(S int not null, B int not null, DT datetime not null, primary key(S,B), foreign key(S) references S(S), foreign key(B) references BT(B));"))
	_, _ = lib.db.Exec(fmt.Sprintf("create table SBR(S int not null, B int not null, RD datetime not null, primary key(S,B,RD), foreign key(S) references S(S), foreign key(B) references BT(B));"))
	return nil;
}



func (lib *Library) ReturnB(S,B int) error{
	var cnt int
	row,err:=lib.db.Query(fmt.Sprintf("select count(*) from SBD where B=%d;",B));
	row.Next()
	err=row.Scan(&cnt);
	if err!=nil{
		panic(err)
	}
	if cnt>0 {
		_,err =lib.db.Exec(fmt.Sprintf("update BB set BR='N' where B=%d",B));
		if err != nil {
			panic(err)
		}
		_,err =lib.db.Exec(fmt.Sprintf("delete from SBC where B=%d",B));
		if err != nil {
			panic(err)
		}
		_,err =lib.db.Exec(fmt.Sprintf("delete from SBD where B=%d",B));
		if err != nil {
			panic(err)
		}
		_,err =lib.db.Exec(fmt.Sprintf("insert into SBR values(%d,%d,'%s')",S,B,time.Now().Format("2006-1-2 15:04:05")));
		if err != nil {
			panic(err)
		}
	}else{
		fmt.Printf("No record about the borrowing.\n");
	}
	return nil;
}

func (lib *Library) Check(S int) error{
	cnt:=0
	cur:=time.Now().Format("2006-1-2 15:04:05")
	row,err:=lib.db.Query(fmt.Sprintf("select count(*) from SBD where S=%d and DT<'%s';",S,cur));
	row.Next()
	err=row.Scan(&cnt);
	if err!=nil{
		panic(err)
	}
	if cnt>0 {
		fmt.Printf("Overdue book(s) should be returned.\n");
	}else{
		fmt.Printf("No overdue book need be returned.\n");
	}
	return nil;
}

func (lib *Library) Check_sus(S int) bool{
	cnt:=0
	cur:=time.Now().Format("2006-1-2 15:04:05")
	row,err:=lib.db.Query(fmt.Sprintf("select count(*) from SBD where S=%d and DT<'%s';",S,cur));
	row.Next()
	err=row.Scan(&cnt);
	if err!=nil{
		panic(err)
	}
	if cnt>3 {
		return true
	}else{
		return false
	}
}

func (lib *Library) ConT(S,B int) error{
	var cnt int
	row,err :=lib.db.Query(fmt.Sprintf("select count(*) from SBD where S=%d and B=%d;",S,B));
	if err!=nil{
		panic(err)
	}
	row.Next()
	err=row.Scan(&cnt);
	if err!=nil{
		panic(err)
	}
	if cnt==0 {
		fmt.Printf("There is no record about it, please check twice.\n");
	}else {
		row,err=lib.db.Query(fmt.Sprintf("select C from SBC where S=%d and B=%d",S,B));
		row.Next()
		err=row.Scan(&cnt);
		if err!=nil{
			panic(err)
		}
		if cnt>=3{
			fmt.Printf("No more latency!\n");
		}else{
			_,err =lib.db.Exec(fmt.Sprintf("update SBC set C=%d where S=%d and B=%d",cnt+1,S,B));
			if err != nil {
				panic(err)
			}
			cur:=time.Now()
			last,_:=time.ParseDuration("360h")
			due:=cur.Add(last)
			_,err =lib.db.Exec(fmt.Sprintf("update SBD set DT='%s' where S=%d and B=%d",due.Format("2006-1-2 15:04:05"),S,B));
			if err != nil {
				panic(err)
			}
			fmt.Printf("Duetime pushed to %s\n",due.Format("2006-1-2 15:04:05"));
		}
	}
	return nil;
} 


func (lib * Library)DueT(B int) error{
	due:=time.Now().Format("2006-1-2 15:04:05")
	row,err :=lib.db.Query(fmt.Sprintf("(select DT from SBD where B=%d);",B));
	if err!=nil{
		panic(err)
	}
	row.Next()
	err=row.Scan(&due);
	if err!=nil{
		panic(err)
	}
	fmt.Printf("Duetime: %s.\n",due);
	return nil;
}


func (lib *Library) BRS(S int) error{
	var Bs int
	row,err :=lib.db.Query(fmt.Sprintf("(select B from SBD where S=%d);",S));
	if err!=nil{
		panic(err)
	}
	cnt:=0
	for row.Next(){
		err = row.Scan(&Bs)
		if err != nil {
			panic(err)
		}
		cnt=cnt+1;
		fmt.Printf("Book number: %d.\n",Bs);
	}
	fmt.Println("Results:",cnt);
	return nil;	
}

func (lib *Library) History(S int)error{
	var Bs int
	row,err :=lib.db.Query(fmt.Sprintf("(select B from SBD where S=%d) union (select B from SBR where S=%d order by RD);",S,S));
	if err!=nil{
		panic(err)
	}
	cnt:=0
	for row.Next(){
		err = row.Scan(&Bs)
		if err != nil {
			panic(err)
		}
		cnt=cnt+1;
		fmt.Printf("Book number: %d.\n",Bs);
	}
	fmt.Println("Results:",cnt);
	return nil;	
}

func (lib *Library) Borrow(S,B int)error{
	rows, err := lib.db.Query(fmt.Sprintf("select BR from BB WHERE B=%d;",B));
	if err != nil {
		panic(err)
	}
	rows.Next()
	var BR string
	err = rows.Scan(&BR)
	
	if BR=="N" {
		_,err =lib.db.Exec(fmt.Sprintf("update BB set BR='Y' where B=%d",B));
		if err != nil {
			panic(err)
		}
		cur:=time.Now()
		last,_:=time.ParseDuration("360h")
		due:=cur.Add(last)
		_,err =lib.db.Exec(fmt.Sprintf("insert into SBD values(%d,%d,'%s')",S,B,due.Format("2006-1-2 15:04:05")));

		if err != nil {
			panic(err)
		}

		_,err =lib.db.Exec(fmt.Sprintf("insert into SBC values(%d,%d,0)",S,B));
		
		if err != nil {
			panic(err)
		}
	}else{
		fmt.Println("This book is already borrowed.");	
	}
	return nil;
}

// AddBook add a book into the library
func (lib *Library) AddBook(title, author, ISBN string)error{
	rows, err := lib.db.Query("select count(*) from BT;")
	if err != nil {
		panic(err)
	}
	rows.Next()
	var cnt int
	err = rows.Scan(&cnt)
	if err != nil {
		panic(err)
	}
	cnt = cnt + 1
	_, err = lib.db.Exec(fmt.Sprintf("insert into BT values('%d','%s');",cnt, title))
	if err != nil {
		panic(err)
	}
	_, err = lib.db.Exec(fmt.Sprintf("insert into BA values('%d','%s');",cnt, author))
	if err != nil {
		panic(err)
	}
	_, err = lib.db.Exec(fmt.Sprintf("insert into BI values('%d','%s');",cnt, ISBN))
	if err != nil {
		panic(err)
	}
	_, err = lib.db.Exec(fmt.Sprintf("insert into BB values('%d','N');",cnt))
	if err != nil {
		panic(err)
	}
	return nil;
}


func (lib *Library) RemoveBook(B int,Exp string)error{
	rows, err := lib.db.Query(fmt.Sprintf("select count(*) from BT where B=%d;",B))
	if err != nil {
		panic(err)
	}
	rows.Next()
	var cnt int
	err = rows.Scan(&cnt)
	if err != nil {
		panic(err)
	}
	if cnt==0 {
		fmt.Println("The book does not exist.");
		return nil;
	}

	rows, err = lib.db.Query(fmt.Sprintf("select count(*) from BE where B=%d;",B))
	if err != nil {
		panic(err)
	}
	rows.Next()
	err = rows.Scan(&cnt)
	if err != nil {
		panic(err)
	}
	if cnt>0 {
		fmt.Println("The book is already removed.");
		return nil;
	}

	_, err = lib.db.Exec(fmt.Sprintf("insert into BE values('%d','%s');",B, Exp))
	if err != nil {
		panic(err)
	}
	return nil;
}

func (lib *Library) FindT(title string)error{
	var Bs int
	row,err :=lib.db.Query(fmt.Sprintf("SELECT B FROM BT WHERE TITLE='%s';",title));
	if err!=nil{
		panic(err)
	}
	cnt:=0
	for row.Next(){
		err = row.Scan(&Bs)
		if err != nil {
			panic(err)
		}
		cnt=cnt+1;
		fmt.Printf("Book number: %d.\n",Bs);
	}
	fmt.Println("Results:",cnt);
	return nil;	
}

func (lib *Library) FindA(author string)error{
	var Bs int
	row,err:=lib.db.Query(fmt.Sprintf("SELECT B FROM BA WHERE AUTHOR='%s';",author));
	if err!=nil{
		panic(err)
	}
	cnt:=0
	for row.Next(){
		err = row.Scan(&Bs)
		if err != nil {
			panic(err)
		}
		cnt=cnt+1;
		fmt.Printf("Book number: %d.",Bs);
	}
	fmt.Println("Results:",cnt);
	return nil;
}

func (lib *Library) FindI(ISBN string) error{
	var Bs int
	row,err:=lib.db.Query(fmt.Sprintf("SELECT B FROM BI WHERE ISBN='%s';",ISBN));
	if err!=nil{
		panic(err)
	}
	cnt:=0
	for row.Next(){
		err = row.Scan(&Bs)
		if err != nil {
			panic(err)
		}
		cnt=cnt+1;
		fmt.Printf("Book number: %d.",Bs);
	}
	fmt.Println("Results:",cnt);
	return nil;
}	

func (lib *Library) AddS()error{
	rows, err := lib.db.Query("select count(*) from S;")
	if err != nil {
		panic(err)
	}
	rows.Next()
	var cnt int
	err = rows.Scan(&cnt)
	if err != nil {
		panic(err)
	}
	cnt = cnt + 1
	_, err = lib.db.Exec(fmt.Sprintf("insert into S values('%d');",cnt))
	if err != nil {
		panic(err)
	}
	fmt.Println("New student id:",cnt);
	return nil;
}



func main() {
	fmt.Println("Welcome to the Library Management System!")
	var lib Library
	lib.ConnectDB()
	lib.CreateTables()
	for ;; {
		fmt.Println("Login as a student,please input 1, login as a teacher,please input 2.");
		var mode int
		fmt.Scanln(&mode);
		if mode==1 {
			fmt.Println("Please input your student id.");
			var id int
			fmt.Scanln(&id);
			row,err:=lib.db.Query(fmt.Sprintf("SELECT COUNT(*) FROM S WHERE S=%d;",id));
			if err != nil {
				panic(err)
			}
			row.Next()
			var cnt int
			err = row.Scan(&cnt)
			if err != nil {
				panic(err)
			}
			if cnt > 0 {
				if(lib.Check_sus(id)){
					fmt.Println("Your id is suspended, please return books before more operations.");
					for ;; {
						fmt.Println("Input 1 to return book or anyother to exist.");
						var op int;
						fmt.Scanln(&op);
						if op == 1 {
							fmt.Println("Please input the returned book's id.");
							var bid int;
							fmt.Scanf("%d",bid);
							lib.ReturnB(id,bid);
							fmt.Println("If you are sure there are no more than 3 overdue books,try to login again.");
						}else{
							break;
						}

					}
				}else{
					for ;; {
						fmt.Println("Input 1 to return book.");
						fmt.Println("Input 2 to find books by title.");
						fmt.Println("Input 3 to find books by author.");
						fmt.Println("Input 4 to find books by ISBN.");
						fmt.Println("Input 5 to borrow one book by the book id.");
						fmt.Println("Input 6 to query the histroy.");
						fmt.Println("Input 7 to query your books that are not returned.");
						fmt.Println("Input 8 to query the duetime for one book by bookid.");
						fmt.Println("Input 9 to push the duetime for one book by bookid.");
						fmt.Println("Input 10 to check if there is overdue book(s) you need to return.");
						fmt.Println("Input anyother number to exist.");
						var op int;
						fmt.Scanln(&op);
						if op == 1 {
							fmt.Println("Please input the returned book's id.");
							var bid int;
							fmt.Scanln(&bid);
							lib.ReturnB(id,bid);
						}else if op == 2{
							fmt.Println("Please input the title.");
							var ss string;
							fmt.Scanln(&ss);
							lib.FindT(ss);
						}else if op == 3{
							fmt.Println("Please input the author.");
							var ss string;
							fmt.Scanln(&ss);
							lib.FindA(ss);
						}else if op == 4{
							fmt.Println("Please input the ISBN.");
							var ss string;
							fmt.Scanln(&ss);
							if(strings.Count(ss,"")!=14){
								fmt.Println("ISBN's lenth should be 13.");
							}else{
								lib.FindT(ss);
							}
						}else if op == 5{
							fmt.Println("Please input the bookid.");
							var bid int;
							fmt.Scanln(&bid);
							lib.Borrow(id,bid);
						} else if op == 6{
							lib.History(id);
						}else if op ==7{
							lib.BRS(id)
						}else if op == 8 {
							fmt.Println("Please input the bookid.");
							var bid int;
							fmt.Scanln(&bid);
							lib.DueT(bid);
						}else if op==9{
							fmt.Println("Please input the bookid.");
							var bid int;
							fmt.Scanln(&bid);
							lib.ConT(id,bid);
						}else if op==10{
							lib.Check(id);
						}else{
							break;
						}
					}
				}
				
			}else{
				fmt.Println("This id is not registered yet.");
			}
			
		}else if mode == 2 {
			fmt.Println("Please input your teacher id.");
			var id int
			fmt.Scanln(&id);
			row,err:=lib.db.Query(fmt.Sprintf("SELECT COUNT(*) FROM T WHERE T=%d;",id));
			if err != nil {
				panic(err)
			}
			row.Next()
			var cnt int
			err = row.Scan(&cnt)
			if err != nil {
				panic(err)
			}
			if cnt>0{
				for ;; {
					fmt.Println("Input 1 to add book.");
					fmt.Println("Input 2 to remove book.");
					fmt.Println("Input 3 to add student account.");
					fmt.Println("Input anyother number to exist.");
					var op int;
					fmt.Scanln(&op);
					if op==1 {
						fmt.Println("Please input the book's title,author and ISBN");
						var tt,aa,ii string;
						fmt.Scanln(&tt,&aa,&ii);
						if(strings.Count(ii,"")!=14){
							fmt.Println("ISBN's lenth should be 13.");
						}else{
							lib.AddBook(tt,aa,ii);
						}
					}else if op==2{ 
						fmt.Println("Please input the bookid and explanation");
						var bid int
						var exp string
						fmt.Scanln(&bid,&exp)
						lib.RemoveBook(bid,exp);
					}else if op==3{
						lib.AddS();
					}else{
						break;
					}
				}
			}else{
				fmt.Println("This id is not registered yet.");
			}
			
		}
		fmt.Println();
	}

}