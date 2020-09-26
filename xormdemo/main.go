package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql" //不能忘记导入
	"github.com/go-xorm/xorm"
	"xorm.io/core"
)

func main() {
	//1. 创建数据库引擎对象
	engine, err := xorm.NewEngine("mysql", "root:linyifan@/elmcms?charset=utf8")
	if err != nil {
		panic(err.Error())
	}
	defer engine.Close()

	//数据库引擎设置
	engine.ShowSQL(true)                     //设置显示SQL语句
	engine.Logger().SetLevel(core.LOG_DEBUG) //设置日志级别
	engine.SetMaxOpenConns(10)               //设置最大连接数
	engine.SetMaxIdleConns(2)                //最大空闲链接数

	engine.Sync(new(Person)) // 根据结构体的结构映射到数据库，生成一张表

	session := engine.Table("user")
	count, err := session.Count()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(count)

	result, err := engine.Query("select `id` from user")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(result)

	//设置名称映射规则
	//engine.SetMapper(core.SnakeMapper{})
	//engine.Sync2(new(UserTable))

	//engine.SetMapper(core.SameMapper{})
	//engine.Sync2(new(StudentTable))

	engine.SetMapper(core.GonicMapper{})
	engine.Sync2(new(PersonTable))

	personEmpty, err := engine.IsTableEmpty(new(PersonTable))
	if err != nil {
		panic(err.Error())
	}
	if personEmpty {
		fmt.Println(" 人员表是空的 ")
	} else {
		fmt.Println(" 人员表不为空 ")
	}

	//判断表结构是否存在
	studentExist, err := engine.IsTableExist(new(StudentTable))
	if err != nil {
		panic(err.Error())
	}
	if studentExist {
		fmt.Println("学生表存在")
	} else {
		fmt.Println("学生表不存在")
	}

	// xorm增删查改
	// Get只返回一条数据
	// Find返回多条数据
	//1.ID查询
	var person PersonTable
	// select * from person_table where id = 1
	engine.Id(1).Get(&person)
	fmt.Println(person.PersonName)
	fmt.Println(person)

	//2.where多条件查询
	var person2 PersonTable
	// select * from person_table where person_age = 26 and person_sex = 2
	engine.Where(" person_age = ? and person_sex = ?", 28, 0).Get(&person2)
	fmt.Println(person2)

	//2.where多条件查询:返回多条数据
	var person3 []PersonTable
	// select * from person_table where person_age = 26 and person_sex = 2
	engine.Where(" person_age = ? and person_sex = ?", 28, 0).Find(&person3)
	fmt.Println(person3)

	//3.And条件查询
	var persons4 []PersonTable
	//select * from person_table where person_age = 26 and person_sex = 2
	err = engine.Where(" person_age = ? ", 26).And("person_sex = ? ", 2).Find(&persons4)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(persons4)
	fmt.Println()

	//4、Or条件查询
	var personArr []PersonTable
	//select * from person_table where person_age = 26 or person_sex = 1
	err = engine.Where(" person_age = ? ", 26).Or("person_sex = ? ", 1).Find(&personArr)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(personArr)
	fmt.Println()

	//5、原生SQL语句查询支持 like语法
	var personsNative []PersonTable
	err = engine.SQL(" select * from person_table where person_name like 't%' ").Find(&personsNative)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(personsNative)
	fmt.Println()

	//6、排序条件查询
	var personsOrderBy []PersonTable
	//select * from person_table orderby person_age  升序排列
	//engine.OrderBy(" person_age ").Find(&personsOrderBy)
	engine.OrderBy(" person_age desc ").Find(&personsOrderBy)
	fmt.Println(personsOrderBy)
	fmt.Println()

	//7、查询特定字段
	var personsCols []PersonTable
	engine.Cols("person_name", "person_age").Find(&personsCols)
	for _, col := range personsCols {
		fmt.Println(col)
	}

	//三、增加记录操作
	personInsert := PersonTable{
		PersonName: "Hello",
		PersonAge:  18,
		PersonSex:  1,
	}
	rowNum, err := engine.Insert(&personInsert)
	fmt.Println(rowNum) //rowNum 受影响的记录条数 返回的是插入条数
	fmt.Println()

	//四、删除操作
	rowNum, err = engine.Delete(&personInsert)
	fmt.Println(rowNum) //rowNum 受影响的记录条数
	fmt.Println()

	//五、更新操作
	rowNum, err = engine.Id(7).Update(&personInsert)
	fmt.Println(rowNum) //rowNum 受影响的记录条数
	fmt.Println()

	//六、统计功能count
	count, err = engine.Count(new(PersonTable))
	fmt.Println("PersonTable表总记录条数：", count)

	//七、事务操作
	personsArray := []PersonTable{
		PersonTable{
			PersonName: "Jack",
			PersonAge:  28,
			PersonSex:  1,
		},
		PersonTable{
			PersonName: "Mali",
			PersonAge:  28,
			PersonSex:  1,
		},
		PersonTable{
			PersonName: "Ruby",
			PersonAge:  28,
			PersonSex:  1,
		},
	}
	session = engine.NewSession()
	session.Begin()
	for i := 0; i < len(personsArray); i++ {
		_, err = session.Insert(personsArray[i])
		if err != nil {
			session.Rollback()
			session.Close()
		}
	}
	err = session.Commit()
	session.Close()
	if err != nil {
		panic(err.Error())
	}

}

/**
 * 人员结构表
 */
type PersonTable struct {
	Id         int64     `xorm:"pk autoincr"`   //主键自增
	PersonName string    `xorm:"varchar(24)"`   //可变字符
	PersonAge  int       `xorm:"int default 0"` //默认值
	PersonSex  int       `xorm:"notnull"`       //不能为空
	City       CityTable `xorm:"-"`             //不映射该字段
}

type CityTable struct {
	CityName      string
	CityLongitude float32
	CityLatitude  float32
}

type Person struct {
	Age  int
	Name string
}

type UserTable struct {
	UserId   int64  `xorm:"pk autoincr"`
	UserName string `xorm:"varchar(32)"` //用户名
	UserAge  int64  `xorm:"default 1"`   //用户年龄
	UserSex  int64  `xorm:"default 0"`   //用户性别
}

/**
 * 学生表
 */
type StudentTable struct {
	Id          int64  `xorm:"pk autoincr"` //主键 自增
	StudentName string `xorm:"varchar(24)"` //
	StudentAge  int    `xorm:"int default 0"`
	StudentSex  int    `xorm:"index"` //sex为索引
}
