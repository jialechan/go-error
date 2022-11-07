## panic or error?
1. 在`Go`中`panic`会导致程序直接退出，是一个致命的错误，如果使用`panic` `recover`进行处理的话，会存在很多问题
    - 性能问题，频繁`panic` `recover`性能不好
    - 容易导致程序异常退出，只要有一个地方没有处理到就会导致程序进程整个退出
    - 不可控，一旦`panic`就将处理逻辑移交给了外部，我们并不能预设外部包一定会进行处理
    - 有时候你拿到一个`panic`，你不知道从什么地方来的，不知道是什么地方来你就不知道怎么处理
2. 什么时候使用 panic 呢？
    - 对于真正意外的情况，那些表示不可恢复的程序错误，例如索引越界、不可恢复的环境问题、栈溢出，我们才使用`panic`
2. 使用 error 处理有哪些好处？
    - 简单
    - 考虑失败，而不是成功(Plan for failure, not success)。
    - 没有隐藏的控制流。
    - 完全交给你来控制 error
    - Error are values

# panic
1. 在程序启动的时候，如果有强依赖的服务出现故障时`panic`退出（✅）
2. 在程序启动的时候，如果发现有配置明显不符合要求,可以`panic`退出（✅）
3. 其他情况下只要不是不可恢复的程序错误，都不应该直接`panic`应该返回`error`（✅）
4. 在程序入口处，例如`gin`中间件需要使用`recover`预防`panic`程序退出（🤨，虽然HTTP服务端对所有的处理函数都会使用`recover`处理`panic`，但David Symonds在GitHub的评论中 https://oreil.ly/BGOmg 提到这其实是一个决策错误）
5. 在程序中我们应该避免使用野生的`goroutine`（💡，含有io的场景是考虑协程的好时机）
    - 如果是在请求中需要执行异步任务，应该使用异步`worker`，消息通知的方式进行处理，避免请求量大时大量`goroutine`创建（🤨，Go程序本来可以同时生成数百、数千甚至数万个`goroutine`）
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
1. 我们在应用程序中使用`github.com/pkg/errors`处理应用错误，**注意在公共库当中，我们一般不使用这个**（❓，为什么？）
2. `error`应该是函数的最后一个返回值，当`error`不为`nil`时，函数的其他返回值是不可用的状态，不应该对其他返回值做任何期待（✅，`func f()(io.Reader, *S1 error)` 在这里，我们不知道`io.Reader`中是否有数据，可能有，也有可能有一部分）
3. 错误处理的时候应该先判断错误，`if err != nil`出现错误及时返回，使代码是一条流畅的直线，避免过多的嵌套.(✅)
4. 在应用程序中出现错误时，使用`errors.New`或者`errors.Errorf`返回错误(✅)
    ```go
    func (u *usecese) usecase1() error {
        money := u.repo.getMoney(uid)
        if money < 10 {
            errors.Errorf("用户余额不足, uid: %d, money: %d", uid, money)
        }
        // 其他逻辑
        return nil
    }    
    ```    
5. 如果是调用应用程序的其他函数出现错误，请直接返回，如果需要携带信息，请使用`errors.WithMessage`(✅)
    ```go
    func (u *usecese) usecase2() error {
        name, err := u.repo.getUserName(uid)
        if err != nil {
            return errors.WithMessage(err, "其他附加信息")
        }

        // 其他逻辑
        return nil
    }  
    ```    
6. 如果是调用其他库（标准库、企业公共库、开源第三方库等）获取到错误时，请使用`errors.Wrap`添加堆栈信息(✅)
    - 切记，不要每个地方都是用`errors.Wrap`只需要在错误第一次出现时进行`errors.Wrap`即可(✅)
    - 根据场景进行判断是否需要将其他库的原始错误吞掉，例如可以把`repository`层的数据库相关错误吞掉，返回业务错误码，避免后续我们分割微服务或者更换`ORM`库时需要去修改上层代码(✅)
    - 注意我们在基础库，被大量引入的第三方库编写时一般不使用`errors.Wrap`避免堆栈信息重复(❓，没太get到)
        ```go
        func (u *usecese) usecase2() error {
            name, err := u.repo.getUserName(uid)
            if err != nil {
                return errors.WithMessage(err, "其他附加信息")
            }

            // 其他逻辑
            return nil
        }   
        ```
7. **禁止**每个出错的地方都打日志，只需要在进程的最开始的地方使用`%+v`进行统一打印，例如`http/rpc`服务的中间件(🤔️，如何实践)
8. 错误判断使用`errors.Is`进行比较(✅)
    ```go
    func f() error {
        err := A()
        if errors.Is(err, io.EOF){
            return nil
        }

        // 其他逻辑
        return nil
    } 
    ```
9. 错误类型判断，使用`errors.As`进行赋值
    ```go
    func f() error {
        err := A()

        var errA errorA
        if errors.As(err, &errA){
            // ...
        }

        // 其他逻辑
        return nil
    }
    ```
10. 如何判定错误的信息是否足够，想一想当你的代码出现问题需要排查的时候你的错误信息是否可以帮助你快速的定位问题，例如我们在请求中一般会输出参数信息，用于辅助判断错误(✅)
11. 对于业务错误，推荐在一个统一的地方创建一个错误字典，错误字典里面应该包含错误的`code`，并且在日志中作为独立字段打印，方便做业务告警的判断，错误必须有清晰的错误文档
不需要返回，被忽略的错误必须输出日志信息(🤔️，如何实践)
12. 同一个地方不停的报错，最好不要不停输出错误日志，这样可能会导致被大量的错误日志信息淹没，无法排查问题，比较好的做法是打印一次错误详情，然后打印出错误出现的次数(🤔️，如何实践)
13. 对同一个类型的错误，采用相同的模式，例如参数错误，不要有的返回`404`有的返回`200`(✅)
14. 处理错误的时候，需要处理已分配的资源，使用`defer`进行清理，例如文件句柄(✅)



