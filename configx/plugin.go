package configx

//register function
//anyFileMap key is the plugin name which provide by flag args "-c=any"
//anyFileMap value is the plugin load function,such as the loadAny function in goany.go
func loadPlugin() {
	anyFileMap["any"] = loadAny
}
