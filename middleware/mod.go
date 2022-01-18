package middleware

import (
	"fmt"
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

// func createNewMiddleware() Middleware {
// 	// 创建一个新的中间件
// 	middleware := func(next http.HandlerFunc) http.HandlerFunc {
// 		// 创建一个新的handler包裹next
// 		handler := func(w http.ResponseWriter, r *http.Request) {

// 			// 中间件的处理逻辑
// 			//           ......
// 			// 调用下一个中间件或者最终的handler处理程序
// 			next(w, r)
// 		}

// 		// 返回新建的包装handler
// 		return handler
// 	}

// 	// 返回新建的中间件
// 	return middleware
// }
func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	//通过m(f) 不断的填充f，最后返回的 f 是一个函数的连续调用，最开始的时候，就初始化了这个f,先执行里面的，再往外执行
	//f 是实际上的功能函数，后面的序列才是中间件
	return f
}
func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
}
