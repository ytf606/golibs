# golibs

## Description
The Go Library contains a larger set of useful Google Go packages for different purposes.

## Packages

### configx package
```golang
var _gc interface{}
if err := configx.InitKratos(conf, &_gc); err != nil {
    panic(err)
}
defer configx.Close()
    
db := configx.GetValue("data.database.dsn.rw")    
```
OR
```golang
configx.InitConfig(conf_path string) //default conf path: conf/conf.ini
    
/*
[Redis]
redis = 127.0.0.1:6379 127.0.0.1:7379
*/
//return "127.0.0.1:6379 127.0.0.1:7379"
addr := configx.GetConf("Redis", "redis")

//return "127.0.0.1:80"
addr := configx.GetConfDefault("Redis", "redis2", "127.0.0.1:80") 

//return []string{127.0.0.1:6379,127.0.0.1:7379}
addr := configx.GetConfs("Redis", "redis") 

// return  map[string][]string{"redis":[127.0.0.1:6379,127.0.0.1:7379]}
addr := configx.GetConfArrayMap("Redis", "redis")

// return map[string]string{"redis":"127.0.0.1:6379 127.0.0.1:7379"}
addr := configx.GetConfStringMap("Redis")  
```

### logx package

#### Init package
```golang
name := "project name"
func InitLogger(name string) (func(), error) {
    config := logx.NewLogConfig()
    config.SetConfigMap(map[string]string{
        "LogPath":    "/tmp/default.log",
        "RotateSize": "1G",
        "Rotate":     "true", //"false" 
        "Retention":  "3",
        "Console":    "true",
        "Level":      "DEBUG", //"TRACE","INFO","ERROR","WARNING"
    })
    logx.InitLogWithConfig(config)
    builder := new(builders.TraceBuilder)
    builder.SetTraceDepartment(name)
    builder.SetTraceVersion("1.0")
    logx.SetBuilder(builder)
    logx.NewLogger(builder)
    return func() {
        logx.Close()
    }, nil
}
```

#### How to use
```golang
tag := "tag label"
logx.I(tag, "this is info msg")
msg := "xxxxxx info yyyyyy"
logx.I(tag, "this is info:%+v", msg)
ctx := context.Background()
logx.Ix(ctx, tag, "this is info msg")
logx.Wx(ctx, tag, "this is warn msg")
logx.Ex(ctx, tag, "this is error msg")
```

### ginx package
```golang
//init gin bind validate for chinese
if err := validate.InitTrans("zh"); err != nil {
    panic(err)
}

//init gin engine
func InitGinEngine(r router.Router) *ginx.Engine {
    app := ginx.New(configx.GetValue("server.gin.mode"),
        mw.Logger(),
        mw.Recovery(),
        mw.LoggerMiddleware(),
        mw.NoCacheMiddleware(),
    )

    if configx.GetValue("server.cors.is_open") == "true" {
        config := mw.Config{
            AllowOrigins:     []string{configx.GetValue("server.cors.allow_origins")},
            AllowMethods:     []string{configx.GetValue("server.cors.allow_methods")},
            AllowHeaders:     []string{configx.GetValue("server.cors.allow_headers")},
            AllowCredentials: configx.GetValue("server.cors.allow_credentials") == "true",
            MaxAge:           timex.Duration(cast.ToInt(configx.GetValue("server.cors.max_age"))),
        }
        app.Use(mw.CorsMiddleware(config), mw.OptionsMiddleware())
    }
  
    r.Register(app)
    return app
}

// gin response
func Login(c *ginx.Context) {
    ctx := ginx.StdCtx(c)
    var param v1.LoginReq
    if err := ginx.ParseCheckJson(c, &param, errorx.AuthParamErr); err != nil {
        return
    }
    data, err := Login(ctx, &param)
    if err != nil {
        ginx.ErrResponse(c, err)
        return
    }
    ginx.SuccResponse(c, data)
}
```

### db package
```golang
func InitDb(name string) (*db.ClusterConn, error) {

    err := db.NewCluster().SetName(name).SetConf(&db.ClusterConfig{
        W: &db.Config{
            DSN:   configx.GetValue("data.database.dsn.rw"),
            Debug: configx.GetValue("data.database.mode") == "debug",
        },
        R: &db.Config{
            DSN:   configx.GetValue("data.database.dsn.ro"),
            Debug: configx.GetValue("data.database.mode") == "debug",
        },
    }).Build()
    if err != nil {
        return nil, err
    }
    return db.GetClusterInstance(name), nil
}
```

### jwtx package
```golang
// gen jwt token
j := jwtx.NewJwt(secret)
ret, err := j.Create(ctx, schema.Xxxx{
    Id:       id,
    Username: username,
    Email:    email,
    StandardClaims: jwtx.StandardClaims{
        Issuer:    "issuer"
        ExpiresAt: expired,
}})

// parse jwt token
j := jwtx.NewJwt(secret)
token, err := j.Parse(ctx, t, &schema.UserAuthInfo{})
if err != nil {
    ginx.ErrResponse(c, err)
    c.Abort()
    return
}
if claims, ok := token.Claims.(*schema.xxx); ok && token.Valid {
    c.Set("_xxx", *claims)
    return
}
```


