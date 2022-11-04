# panic
1. 在程序启动的时候，如果有强依赖的服务出现故障时`panic`退出（✅）
2. 在程序启动的时候，如果发现有配置明显不符合要求,可以`panic`退出（✅）
3. 其他情况下只要不是不可恢复的程序错误，都不应该直接`panic`应该返回`error`（✅）
4. 在程序入口处，例如`gin`中间件需要使用`recover`预防`panic`程序退出（💭，虽然HTTP服务端对所有的处理函数都会使用`recover`处理`panic`，但David Symonds在GitHub的评论中 https://oreil.ly/BGOmg 提到这其实是一个决策错误。）
5. 在程序中我们应该避免使用野生的goroutine（💡，含有io的场景是考虑协程的好时机）
    - 如果是在请求中需要执行异步任务，应该使用异步`worker`，消息通知的方式进行处理，避免请求量大时大量`goroutine`创建（💭，Go程序本来可以同时生成数百、数千甚至数万个`goroutine`）
    - 如果需要使用`goroutine`时，应该使用同一的`Go`函数进行创建，这个函数中会进行`recover`，避免因为野生`goroutine` `panic`导致主进程退出（✅，不过这个写法只支持`func()`这个类型的参数）
```go
func Go(f func()){
    go func(){
        defer func(){
            if err := recover(); err != nil {
                log.Printf("panic: %+v", err)
            }
        }()

        f()
    }()
}
```
# error
