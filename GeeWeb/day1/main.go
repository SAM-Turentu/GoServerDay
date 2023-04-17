package main

import (
	"fmt"
	"gee"
	"log"
	"net/http"
	"time"

	//"github.com/gin-gonic/gin"
)

/*

7 天一个项目

*/

func main() {
	fmt.Println("********************开始********************")
	//1
	http.HandleFunc("/", HomeHandler)
	//http.HandleFunc("/hello", HelloHandler)
	//log.Fatal(http.ListenAndServe("0.0.0.0:8082", nil))

	//2
	//engine := new(Engine)
	//log.Fatal(http.ListenAndServe("0.0.0.0:8082", engine))

	//3
	//geeweb()

	//4
	//geeweb_v2()

	//5
	//geeweb_v3()

	//6
	//geeweb_v4()

	//7
	//geeweb_v5()

	//7
	geeweb_v6()

	fmt.Println("********************结束********************")
}

//region 1.普通net/http 实现web handler方法
func HomeHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "Url.paht = %q\n", req.URL.Path)
}

func HelloHandler(w http.ResponseWriter, req *http.Request) {
	for k, v := range req.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
}

//endregion

//region 2.net/http 使用接口实现 handler 方法
type Engine struct{}

// 实现Handler接口方法
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
		fmt.Fprintln(w, "Url.paht = %q\n", req.URL.Path)
	case "/hello":
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	default:
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}

//endregion

//region 3.使用自定义的web框架实现net/http的功能
//func geeweb() {
//	r := gee.New()
//	r.GET("/", func(w http.ResponseWriter, req *http.Request) {
//		fmt.Fprintf(w, "url.path = %q\n", req.URL.Path)
//	})
//	r.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
//		for k, v := range req.Header {
//			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
//		}
//	})
//	r.Run("0.0.0.0:8082")
//}

//endregion

//region 4. geeweb 增加上下文Context和路由route
func geeweb_v2() {
	r := gee.New()
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello GEE</h1>")
	})

	r.GET("/hello", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.POST("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"user": c.PostForm("user"),
			"pwd":  c.PostForm("pwd"),
		})
	})

	r.Run("0.0.0.0:8082")
}

//endregion

//region 5.前缀路由router
func geeweb_v3() {
	r := gee.New()
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello GEE</h1>")
	})

	r.GET("/hello", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.GET("/hello/:name", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	r.POST("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"user": c.PostForm("user"),
			"pwd":  c.PostForm("pwd"),
		})
	})

	r.GET("/assets/*filepath", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{"filepath": c.Param("filepath")})
	})

	r.Run("0.0.0.0:8082")
}

//endregion

//region 6.分组路由
func geeweb_v4() {
	r := gee.New()
	r.GET("/index", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Index Page</h1>")
	})
	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *gee.Context) {
			c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
		})
		v1.GET("/hello", func(c *gee.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}
	v2 := r.Group("/v2")
	{
		v2.GET("/hello/:name", func(c *gee.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.GET("/login", func(c *gee.Context) {
			c.JSON(http.StatusOK, gee.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})
	}
	r.Run("0.0.0.0:8082")
}

//endregion

//region 7.中间件(分组路由)
func onlyForV2() gee.HandlerFunc {
	return func(c *gee.Context) {
		t := time.Now()
		c.Fail(500, "Internal Server Error")
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func geeweb_v5() {
	r := gee.New()
	r.Use(gee.Logger()) // 1中间件
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee, middleware logger()</h1>")
	})

	v2 := r.Group("/v2") // 2中间件，使用了 r 的中间件
	v2.Use(onlyForV2())  // 先调用了2的中间件，然后调用1的中间件
	{
		v2.GET("/hello/:name", func(c *gee.Context) {
			c.String(http.StatusOK, "hello %s", c.Param("name"))
		})
	}

	r.Run("0.0.0.0:8082")
}

//endregion

//region 8.错误恢复
func geeweb_v6() {
	//r := gee.New()
	//r.GET("/panic", func(c *gee.Context) {
	//	names := []string{"geektutu"}
	//	c.String(http.StatusOK, names[100])
	//})

	r := gee.Default()
	r.GET("/panic", func(c *gee.Context) {
		names := []string{"geektutu"}
		c.String(http.StatusOK, names[100])
	})

	r.Run("0.0.0.0:8082")
}

//endregion
