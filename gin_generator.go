package gin_generator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"strings"

	"github.com/ametsuramet/gin_generator/models"
	"github.com/ametsuramet/gin_generator/utils"
)

type Config struct {
	JsonFile string
	Path     string
	Data
}

type Data interface {
}

func Set(JsonFile string, path string, data interface{}) *Config {
	return &Config{
		JsonFile: JsonFile,
		Path:     path,
		Data:     data,
	}
}

func Test() string {
	return "gin_generator loaded"
}

func (c *Config) Generate() {
	fmt.Println(os.Getenv("GOPATH"))

	output, _ := c.Unmarshal()

	//copy config
	configPath := c.Path + "/config"

	utils.CopyDir(os.Getenv("GOPATH")+"/src/github.com/ametsuramet/gin_generator/config", configPath)
	utils.CopyFile(os.Getenv("GOPATH")+"/src/github.com/ametsuramet/gin_generator/default.yaml", c.Path+"/default.yaml")

	path := c.Path + "/models"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0777)
	}

	controllerPath := c.Path + "/controllers"
	if _, err := os.Stat(controllerPath); os.IsNotExist(err) {
		os.Mkdir(controllerPath, 0777)
	}
	for _, model := range output {
		err := c.createModel(model)
		if err != nil {
			fmt.Println(err)
		}

		err = c.createController(model)
		if err != nil {
			fmt.Println(err)
		}

	}

	c.createMain(output)
	c.createRouter(output)
	c.createBaseModel(output)

	// utils.CopyDir(os.Getenv("GOPATH")+"/src/github.com/ametsuramet/gin_generator/models/Base.go", modelsPath+"/Base.go")

}

// Unmarshal Json and bind to Struct
func (c *Config) Unmarshal() (result []models.Json, err error) {
	jsonFile, err := os.Open(c.JsonFile)
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully Opened ", c.JsonFile)
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal([]byte(byteValue), &result)
	return result, nil

}
func (c *Config) createBaseModel(models []models.Json) error {
	path := c.Path + "/models/Base.go"
	var _, err = os.Stat(path)
	segments := strings.Split(c.Path, "/")
	packageName := segments[len(segments)-1]
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()
	} else {
		os.Remove(path)
		os.Create(path)
	}
	var migrateString []string
	for _, model := range models {
		migrateString = append(migrateString, "&"+model.Name+"{}")
	}
	// open file using READ & WRITE permission
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	defer file.Close()

	file.WriteString("package models\n\n")
	file.WriteString("import (\n\tcfg \"" + packageName + "/config\"\n\t\"github.com/jinzhu/gorm\"\n\t_ \"github.com/jinzhu/gorm/dialects/mysql\"\n\t\"time\"\n")
	file.WriteString(")\n\nvar db *gorm.DB\nvar err error\nvar DB = &gorm.DB{}\n")
	file.WriteString("type BaseModel struct {\n")
	file.WriteString("\tID\tuint\t`gorm:\"primary_key\" json:\"id\"`\n\tCreatedAt\ttime.Time\t`json:\"created_at\"`\n\tUpdatedAt\ttime.Time\t`json:\"update_at\"`\n\tDeletedAt\t*time.Time\t`sql:\"index\" json:\"deleted_at\"`\n")
	file.WriteString("}\n")
	file.WriteString("func Set() {\n")
	file.WriteString("\tcfg.Init()\n\tdb, _ := gorm.Open(cfg.App.Database.Driver, cfg.App.Database.Uri)\n")
	file.WriteString("\tDB = db\n")
	file.WriteString("\tdb.AutoMigrate(" + strings.Join(migrateString, ", "))
	file.WriteString(")\n}\n")
	return nil
}
func (c *Config) createMain(models []models.Json) error {
	path := c.Path + "/main.go"
	var _, err = os.Stat(path)
	segments := strings.Split(c.Path, "/")
	packageName := segments[len(segments)-1]
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()
	} else {
		os.Remove(path)
		os.Create(path)
	}

	fmt.Println("==> done creating file", path)

	// open file using READ & WRITE permission
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	defer file.Close()
	// write some text line-by-line to file

	file.WriteString("package main\n\n")

	file.WriteString("import (\n\tmdl \"" + packageName + "/models\"\n\tctrl \"" + packageName + "/controllers\"\n\t\"github.com/gin-gonic/gin\"\n\t\"net/http\"\n\t\"os\"\n\t\"os/signal\"\n\t\"fmt\"\n")

	file.WriteString(")\n\n")

	file.WriteString("func init() {\n\tmdl.Set()\n}\n")
	file.WriteString("func main() {\n")
	file.WriteString("\tport := os.Getenv(\"PORT\")\n\tif port == \"\" {\n\t\tport = \"7000\"\n\t}\n\tr := setupRouter()\n")

	file.WriteString("\tsrv := &http.Server{\n\t\tAddr:\t\":\" + port,\n\t\tHandler:\tr,\n\t}\n")

	file.WriteString("\tgo func() {\n\t\tif err := srv.ListenAndServe(); err != nil {\n\t\t\tpanic(fmt.Errorf(\"Fatal error failed to start server: %s\", err))\n\t\t}\n\t}()\n")
	file.WriteString("\tquit := make(chan os.Signal)\n\tsignal.Notify(quit, os.Interrupt)\n\t<-quit\n}\n\n")
	return nil
}

func (c *Config) createRouter(models []models.Json) error {
	path := c.Path + "/main.go"
	var _, err = os.Stat(path)
	// segments := strings.Split(c.Path, "/")
	// packageName := segments[len(segments)-1]
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()
	}
	// open file using READ & WRITE permission
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)

	file.WriteString("func setupRouter() *gin.Engine {\n")
	file.WriteString("\tr := gin.New()\n\tr.Use(gin.Logger())\n\tr.Use(gin.Recovery())\n\tr.GET(\"/ping\", ping)\n")
	file.WriteString("\tv1 := r.Group(\"/api/v1\")\n")
	file.WriteString("\tv1.Use()\n\t{\n")
	for _, model := range models {
		featureName := model.Name
		file.WriteString("\t\t" + model.Name + "Route := v1.Group(\"/" + featureName + "\")\n")
		file.WriteString("\t\t" + model.Name + "Route.Use()\n\t\t{\n")
		/*
			stores.GET("/", ctrl.GetStores)
							stores.GET("/:id", ctrl.GetStore)
							stores.PUT("/:id", ctrl.PutStore)
							stores.DELETE("/:id", ctrl.DeleteStore)
							stores.POST("/", ctrl.PostStore)
		*/
		file.WriteString("\t\t\t" + model.Name + "Route.GET(\"/\", ctrl.Index" + model.Name + ")\n")
		file.WriteString("\t\t\t" + model.Name + "Route.GET(\"/:id\", ctrl.Show" + model.Name + ")\n")
		file.WriteString("\t\t\t" + model.Name + "Route.POST(\"/\", ctrl.Store" + model.Name + ")\n")
		file.WriteString("\t\t\t" + model.Name + "Route.PUT(\"/:id\", ctrl.Update" + model.Name + ")\n")
		file.WriteString("\t\t\t" + model.Name + "Route.DELETE(\"/:id\", ctrl.Delete" + model.Name + ")\n")
		file.WriteString("\t\t}\n")
	}
	file.WriteString("\t}\n")

	file.WriteString("\treturn r\n}\n\nfunc ping(c *gin.Context) {\n\tc.JSON(200, gin.H{\n\t\t\"message\": \"pong\",\n\t})\n}\n")

	return nil
}

func (c *Config) createController(model models.Json) error {
	controllerPath := c.Path + "/controllers/" + model.Name + ".go"
	segments := strings.Split(c.Path, "/")
	packageName := segments[len(segments)-1]

	// var strConv utils.StringConv
	// typeName := strConv.ToCamel(model.Name)
	//create file
	var _, err = os.Stat(controllerPath)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(controllerPath)
		if err != nil {
			return err
		}
		defer file.Close()
	}

	fmt.Println("==> done creating file", controllerPath)

	file, err := os.OpenFile(controllerPath, os.O_RDWR, 0644)
	defer file.Close()

	// write some text line-by-line to file

	file.WriteString("package controllers\n\n\n")
	file.WriteString("import (\n")
	file.WriteString("\t\"github.com/gin-gonic/gin\"\n")
	file.WriteString("\tmdl \"" + packageName + "/models\"\n")
	file.WriteString("\t\"github.com/gin-gonic/gin/binding\"\n")
	file.WriteString(")\n\n")

	//write index
	file.WriteString("func Index" + model.Name + "(c *gin.Context) {\n")
	// file.WriteString("\tdefer mdl.DB.Close()\n")
	file.WriteString("\tdata := []mdl." + model.Name + "{}\n")
	file.WriteString("\tmdl.DB.Find(&data)\n")

	file.WriteString("\tc.JSON(200, gin.H{\n")
	file.WriteString("\t\t\"message\": \"Index " + model.Name + "\",\n")
	file.WriteString("\t\t\"data\": data,\n")
	file.WriteString("\t})\n}\n\n")

	//write show
	file.WriteString("func Show" + model.Name + "(c *gin.Context) {\n")

	file.WriteString("\tid := c.Params.ByName(\"id\")\n")
	file.WriteString("\tvar data mdl." + model.Name + "\n")
	file.WriteString("\tif err := mdl.DB.First(&data, \"id=?\", id); err != nil {\n")
	file.WriteString("\t\tc.JSON(200, gin.H{\n")
	file.WriteString("\t\t\t\"message\": \"" + model.Name + " Not Found\",\n")
	file.WriteString("\t\t})\n")
	file.WriteString("\t\treturn\n")
	file.WriteString("\t}\n")
	file.WriteString("\tc.JSON(200, gin.H{\n")
	file.WriteString("\t\t\"message\": \"Show " + model.Name + "\",\n")
	file.WriteString("\t\t\"data\": data,\n")
	file.WriteString("\t})\n}\n\n")

	//write store
	file.WriteString("func Store" + model.Name + "(c *gin.Context) {\n")
	file.WriteString("\tvar input mdl." + model.Name + "\n")

	file.WriteString("\tif c.ShouldBindWith(&input, binding.JSON) == nil {\n")
	file.WriteString("\t\tmdl.DB.Create(&input)\n")
	file.WriteString("\t}\n")
	file.WriteString("\tc.JSON(200, gin.H{\n")
	file.WriteString("\t\t\"message\": \"Store " + model.Name + "\",\n")
	file.WriteString("\t\t\"data\": input,\n")
	file.WriteString("\t})\n}\n\n")

	//write update
	file.WriteString("func Update" + model.Name + "(c *gin.Context) {\n")
	file.WriteString("\tid := c.Params.ByName(\"id\")\n")
	file.WriteString("\tvar data, input mdl." + model.Name + "\n")
	file.WriteString("\tmdl.DB.First(&data, \"id=?\", id)\n")

	file.WriteString("\tif c.ShouldBindWith(&input, binding.JSON) == nil {\n")
	file.WriteString("\t\tmdl.DB.Model(&data).Update(&input)\n")
	file.WriteString("\t}\n")
	file.WriteString("\tc.JSON(200, gin.H{\n")
	file.WriteString("\t\t\"message\": \"Update " + model.Name + "\",\n")
	file.WriteString("\t\t\"data\": input,\n")
	file.WriteString("\t})\n}\n\n")

	// write delete
	file.WriteString("func Delete" + model.Name + "(c *gin.Context) {\n")
	file.WriteString("\tid := c.Params.ByName(\"id\")\n")
	file.WriteString("\tvar data mdl." + model.Name + "\n")
	file.WriteString("\tmdl.DB.First(&data, \"id=?\", id)\n")
	file.WriteString("\tmdl.DB.Delete(&data)\n")

	file.WriteString("\tc.JSON(200, gin.H{\n")
	file.WriteString("\t\t\"message\": \"Delete " + model.Name + "\",\n")
	file.WriteString("\t})\n}\n\n")

	return nil

}

// generate model from struct
func (c *Config) createModel(model models.Json) error {
	modelPath := c.Path + "/models/" + model.Name + ".go"
	var strConv utils.StringConv
	typeName := strConv.ToCamel(model.Name)
	//create file
	var _, err = os.Stat(modelPath)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(modelPath)
		if err != nil {
			return err
		}
		defer file.Close()
	}

	fmt.Println("==> done creating file", modelPath)

	//write models
	// open file using READ & WRITE permission
	file, err := os.OpenFile(modelPath, os.O_RDWR, 0644)
	defer file.Close()

	// write some text line-by-line to file

	file.WriteString("package models\n\n\n")

	file.WriteString("import (\n")
	for _, schema := range model.Schema {
		schemaType := strings.Split(schema.Type, "::")
		if schemaType[0] == "dateTime" {
			file.WriteString("\t\"time\"\n")
			break
		}
	}

	// file.WriteString("\t\"github.com/jinzhu/gorm\"\n")

	// file.WriteString("\t_ \"github.com/jinzhu/gorm/dialects/sqlite\"\n")

	file.WriteString(")\n\n")
	file.WriteString("type " + typeName + " struct {\n")
	file.WriteString("\tBaseModel\n")
	// write schema
	for _, schema := range model.Schema {
		schemaType := strings.Split(schema.Type, "::")

		switch schemaType[0] {
		case "boolean":
			schema.Type = "bool"
		case "integer":
			schema.Type = "int"
		case "text":
			schema.Type = "string"
		case "float":
			schema.Type = "float32"
		case "dateTime":
			schema.Type = "time.Time"
		default:
			schema.Type = schemaType[0]
		}

		file.WriteString("\t" + strConv.ToCamel(schema.Field) + "\t\t" + schema.Type + "\t\t`json:\"" + schema.Field + "\"`" + "\n")

	}

	file.WriteString("}")

	// save changes
	err = file.Sync()

	fmt.Println("==> done writing model", model.Name)

	return nil
}
