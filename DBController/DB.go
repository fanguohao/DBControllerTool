package DBController

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

var (
	updateMigrationSql = "INSERT INTO migrations (migration, batch) VALUES DummyString;"
)

var migrationPath string

var db *gorm.DB

//type rowScanner interface {
//	Scan(dst ...interface{}) error
//}

type Migration struct {
	ID        int64
	Migration string
	Batch     int64
}

type Tversion struct {
	Version float32
}

func init() {
	log.SetFlags(log.Ldate | log.Lshortfile)
}

//验证DB信息，并实例化db
func GetConnectionIfo(username string, password string, dbName string) bool {
	str := []string{username, ":", password, "@/", dbName}
	connInfo := strings.Join(str, "")
	var err error
	db, err = gorm.Open("mysql", connInfo)
	if err != nil {
		return false
	}
	return true
}

//获取DB当前版本；
func GetDBVersion() float32 {
	v := &Tversion{}
	err := db.Find(&Tversion{}).Scan(v).Error
	if err != nil {
		log.Println("数据库中版本表不存在！开始创建...")
		v := Tversion{ //默认从初始版本执行；
			Version: 0.0,
		}
		err := db.CreateTable(v).Error
		if err != nil {
			log.Println("版本库创建失败！")
		}
		err = db.Save(v).Error
		if err != nil {
			log.Println("版本初始化失败！")
		}

		return 0.0
	}

	return v.Version
}

func GetVersionDisplay() []string {
	return versionToDisplay
}

// 数据库升级接口；
func Update(address string, version float32) string {
	//fmt.Println(address)
	arr := strings.Split(address, "\\")
	migrationPath = strings.Join(arr[0:len(arr)-1], "\\") + "\\"
	err, Info := Migrate(version)
	if err != nil {
		return Info
	}
	return Info
}

const (
	NOTHINGTOUPDATA = "当前版本无需升级!"
	FILEDTOUPDATA   = "升级失败"
	SUCCESS         = "升级成功!"
)

var (
	fSlices          []string
	versionToDisplay []string
)

//func main() {
//	GetConnectionIfo("root", "Fan003174", "elmcms")
//	print(Update("C:\\Users\\Nova003174\\Documents\\migrations\\v2_addTable.sql", 3))
//}

func GetToversion(path string) float32 {
	arr := strings.Split(path, "\\")
	arr2 := strings.Split(arr[len(arr)-1], "_")
	f64, _ := strconv.ParseFloat(arr2[0][1:len(arr[0])], 32)

	return float32(f64)
}

// Migration
func Migrate(toVersion float32) (error, string) {
	currentDBVersion := GetDBVersion()
	if toVersion <= currentDBVersion {
		return nil, NOTHINGTOUPDATA
	} else {
		log.Println("开始版本升级！")
		files, err := ioutil.ReadDir(migrationPath)
		if err != nil {
			return err, "无法读取升级文件！"
		}

		var toMigrate []string
		for _, f := range files {
			arr := strings.Split(f.Name(), ".")
			fSlices = append(fSlices, arr[0])
			arr2 := strings.Split(arr[0], "_")
			versionToDisplay = append(versionToDisplay, arr2[0])

			vfloat64, err := strconv.ParseFloat(arr2[0][1:len(arr2[0])], 32)
			if err != nil {
				return err, FILEDTOUPDATA
			}
			if toVersion >= float32(vfloat64) && float32(vfloat64) > currentDBVersion {
				toMigrate = append(toMigrate, f.Name())
			}
		}

		// Migrate
		for _, v := range toMigrate {
			upSql, _ := ioutil.ReadFile(migrationPath + v)
			requests := strings.Split(string(upSql), ";")
			for _, request := range requests[0 : len(requests)-1] {
				err := db.Exec(request).Error
				if err != nil {
					log.Print(err.Error())
					//return err, FILEDTOUPDATA
				}
			}

			err = db.Model(&Tversion{}).Update("version", GetDBVersion()+1).Error

		}

		//err = db.Model(&Tversion{}).Update("version", toVersion).Error
		//if err != nil {
		//	log.Println(err)
		//	return err, FILEDTOUPDATA
		//} else {
		//	return nil, FILEDTOUPDATA
		//}
		return nil, SUCCESS
	}

	////fmt.Println(string(upSql))
	//cmd := exec.Command("C:\\Program Files\\MySQL\\MySQL Server 8.0\\bin\\mysql",  "-u","{root}", "-p{Fan003174}", "{elmcms}","<","C:\\Users\\Nova003174\\go\\src\\DBcontrollerdemo\\database\\migrations\\v1_init.sql")
	////cmd := exec.Command("mysql", "-U", "root", "-h", "Fan003174", "-d", "localhost", "-a", "-f", string(upSql))
	//err:=cmd.Run()
	//if err != nil {
	//	log.Print(err.Error())
	//	return err
	//
	//}

}
